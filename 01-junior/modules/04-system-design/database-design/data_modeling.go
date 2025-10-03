package database_design

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Data Modeling Patterns for System Design Interviews
// Interview Focus: Schema design, performance optimization, and scalability

// DataModelingStrategy defines interface for different modeling approaches
type DataModelingStrategy interface {
	DesignSchema(requirements SchemaRequirements) (SchemaDesign, error)
	OptimizeForQueries(queries []QueryPattern) ([]IndexRecommendation, error)
	ValidateNormalization(schema SchemaDesign) NormalizationAnalysis
	EstimateStorage(schema SchemaDesign, recordCount int64) StorageEstimate
}

// SchemaRequirements captures system requirements for data modeling
type SchemaRequirements struct {
	SystemType      string
	ExpectedScale   ScaleExpectation
	QueryPatterns   []QueryPattern
	ConsistencyNeed ConsistencyRequirement
	PerformanceGoal PerformanceTarget
}

type ScaleExpectation struct {
	DailyActiveUsers int64
	RecordsPerDay    int64
	PeakQPS          int64
	DataRetention    time.Duration
}

type QueryPattern struct {
	Type        QueryType
	Frequency   QueryFrequency
	Columns     []string
	FilterTypes []FilterType
	JoinTables  []string
	Complexity  QueryComplexity
}

type QueryType int

const (
	SelectQuery QueryType = iota
	InsertQuery
	UpdateQuery
	DeleteQuery
	AggregateQuery
	AnalyticsQuery
)

type QueryFrequency int

const (
	VeryHigh QueryFrequency = iota // > 1000 QPS
	High                           // 100-1000 QPS
	Medium                         // 10-100 QPS
	Low                            // < 10 QPS
)

type FilterType int

const (
	EqualsFilter FilterType = iota
	RangeFilter
	LikeFilter
	InFilter
	FullTextFilter
)

type QueryComplexity int

const (
	SimpleQuery QueryComplexity = iota
	MediumQuery
	ComplexQuery
)

type ConsistencyRequirement int

const (
	StrongConsistencyRequired ConsistencyRequirement = iota
	EventualConsistencyOK
	SessionConsistencyOK
)

type PerformanceTarget struct {
	MaxLatencyP99  time.Duration
	MinThroughput  int64
	MaxStorageSize int64
	CacheHitRatio  float64
}

// NormalizedDataModeling implements traditional 3NF approach
// FAANG Interview Point: ACID compliance, referential integrity
type NormalizedDataModeling struct {
	normalizationLevel int
	enforceConstraints bool
}

func NewNormalizedDataModeling() *NormalizedDataModeling {
	return &NormalizedDataModeling{
		normalizationLevel: 3, // 3NF by default
		enforceConstraints: true,
	}
}

// DesignSchema creates normalized schema design
// FAANG Interview Point: E-commerce, financial systems requiring consistency
func (ndm *NormalizedDataModeling) DesignSchema(requirements SchemaRequirements) (SchemaDesign, error) {
	var tables []TableDesign

	switch requirements.SystemType {
	case "e_commerce":
		tables = ndm.designECommerceSchema()
	case "social_media":
		tables = ndm.designSocialMediaSchema()
	case "banking":
		tables = ndm.designBankingSchema()
	default:
		return SchemaDesign{}, fmt.Errorf("unsupported system type: %s", requirements.SystemType)
	}

	schema := SchemaDesign{
		DatabaseType:  "PostgreSQL",
		Tables:        tables,
		Indexes:       []IndexDesign{},
		Constraints:   []ConstraintDesign{},
		ModelingType:  "Normalized",
		Justification: "Strong consistency requirements, complex relationships, ACID compliance needed",
	}

	// Add foreign key constraints for referential integrity
	schema.Constraints = ndm.generateConstraints(tables)

	return schema, nil
}

func (ndm *NormalizedDataModeling) designECommerceSchema() []TableDesign {
	return []TableDesign{
		{
			Name: "users",
			Columns: []ColumnDesign{
				{Name: "user_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "email", Type: "VARCHAR(255)", Unique: true, NotNull: true},
				{Name: "password_hash", Type: "VARCHAR(255)", NotNull: true},
				{Name: "first_name", Type: "VARCHAR(100)", NotNull: true},
				{Name: "last_name", Type: "VARCHAR(100)", NotNull: true},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
				{Name: "updated_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 10000000,
		},
		{
			Name: "products",
			Columns: []ColumnDesign{
				{Name: "product_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "name", Type: "VARCHAR(255)", NotNull: true},
				{Name: "description", Type: "TEXT"},
				{Name: "price", Type: "DECIMAL(10,2)", NotNull: true},
				{Name: "category_id", Type: "BIGINT", ForeignKey: "categories.category_id"},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
				{Name: "updated_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 1000000,
		},
		{
			Name: "categories",
			Columns: []ColumnDesign{
				{Name: "category_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "name", Type: "VARCHAR(100)", NotNull: true, Unique: true},
				{Name: "parent_category_id", Type: "BIGINT", ForeignKey: "categories.category_id"},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 10000,
		},
		{
			Name: "orders",
			Columns: []ColumnDesign{
				{Name: "order_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "user_id", Type: "BIGINT", NotNull: true, ForeignKey: "users.user_id"},
				{Name: "status", Type: "VARCHAR(50)", NotNull: true},
				{Name: "total_amount", Type: "DECIMAL(10,2)", NotNull: true},
				{Name: "shipping_address", Type: "TEXT", NotNull: true},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
				{Name: "updated_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 50000000,
		},
		{
			Name: "order_items",
			Columns: []ColumnDesign{
				{Name: "order_item_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "order_id", Type: "BIGINT", NotNull: true, ForeignKey: "orders.order_id"},
				{Name: "product_id", Type: "BIGINT", NotNull: true, ForeignKey: "products.product_id"},
				{Name: "quantity", Type: "INTEGER", NotNull: true},
				{Name: "unit_price", Type: "DECIMAL(10,2)", NotNull: true},
			},
			EstimatedRowCount: 150000000,
		},
	}
}

func (ndm *NormalizedDataModeling) designSocialMediaSchema() []TableDesign {
	return []TableDesign{
		{
			Name: "users",
			Columns: []ColumnDesign{
				{Name: "user_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "username", Type: "VARCHAR(50)", Unique: true, NotNull: true},
				{Name: "email", Type: "VARCHAR(255)", Unique: true, NotNull: true},
				{Name: "profile_picture_url", Type: "VARCHAR(500)"},
				{Name: "bio", Type: "TEXT"},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 100000000,
		},
		{
			Name: "posts",
			Columns: []ColumnDesign{
				{Name: "post_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "user_id", Type: "BIGINT", NotNull: true, ForeignKey: "users.user_id"},
				{Name: "content", Type: "TEXT", NotNull: true},
				{Name: "media_urls", Type: "TEXT[]"},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 1000000000,
		},
		{
			Name: "follows",
			Columns: []ColumnDesign{
				{Name: "follower_id", Type: "BIGINT", NotNull: true, ForeignKey: "users.user_id"},
				{Name: "following_id", Type: "BIGINT", NotNull: true, ForeignKey: "users.user_id"},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 500000000,
		},
	}
}

func (ndm *NormalizedDataModeling) designBankingSchema() []TableDesign {
	return []TableDesign{
		{
			Name: "customers",
			Columns: []ColumnDesign{
				{Name: "customer_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "ssn", Type: "VARCHAR(11)", Unique: true, NotNull: true},
				{Name: "first_name", Type: "VARCHAR(100)", NotNull: true},
				{Name: "last_name", Type: "VARCHAR(100)", NotNull: true},
				{Name: "date_of_birth", Type: "DATE", NotNull: true},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 10000000,
		},
		{
			Name: "accounts",
			Columns: []ColumnDesign{
				{Name: "account_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "customer_id", Type: "BIGINT", NotNull: true, ForeignKey: "customers.customer_id"},
				{Name: "account_number", Type: "VARCHAR(20)", Unique: true, NotNull: true},
				{Name: "account_type", Type: "VARCHAR(50)", NotNull: true},
				{Name: "balance", Type: "DECIMAL(15,2)", NotNull: true},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 20000000,
		},
		{
			Name: "transactions",
			Columns: []ColumnDesign{
				{Name: "transaction_id", Type: "BIGSERIAL", PrimaryKey: true},
				{Name: "from_account_id", Type: "BIGINT", ForeignKey: "accounts.account_id"},
				{Name: "to_account_id", Type: "BIGINT", ForeignKey: "accounts.account_id"},
				{Name: "amount", Type: "DECIMAL(15,2)", NotNull: true},
				{Name: "transaction_type", Type: "VARCHAR(50)", NotNull: true},
				{Name: "reference_number", Type: "VARCHAR(50)", Unique: true, NotNull: true},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true, DefaultValue: "NOW()"},
			},
			EstimatedRowCount: 1000000000,
		},
	}
}

func (ndm *NormalizedDataModeling) generateConstraints(tables []TableDesign) []ConstraintDesign {
	var constraints []ConstraintDesign

	for _, table := range tables {
		for _, column := range table.Columns {
			if column.ForeignKey != "" {
				parts := strings.Split(column.ForeignKey, ".")
				if len(parts) == 2 {
					constraints = append(constraints, ConstraintDesign{
						Type:             "FOREIGN KEY",
						SourceTable:      table.Name,
						SourceColumn:     column.Name,
						ReferencedTable:  parts[0],
						ReferencedColumn: parts[1],
					})
				}
			}
		}
	}

	return constraints
}

// DenormalizedDataModeling implements performance-optimized approach
// FAANG Interview Point: Read-heavy systems, analytics, caching layers
type DenormalizedDataModeling struct {
	redundancyLevel   float64 // 0.0 to 1.0
	optimizeForReads  bool
	materializedViews bool
}

func NewDenormalizedDataModeling() *DenormalizedDataModeling {
	return &DenormalizedDataModeling{
		redundancyLevel:   0.3, // 30% redundancy acceptable
		optimizeForReads:  true,
		materializedViews: true,
	}
}

// DesignSchema creates denormalized schema optimized for reads
// FAANG Interview Point: Social media feeds, recommendation engines
func (ddm *DenormalizedDataModeling) DesignSchema(requirements SchemaRequirements) (SchemaDesign, error) {
	var tables []TableDesign

	switch requirements.SystemType {
	case "social_media_feed":
		tables = ddm.designSocialMediaFeedSchema()
	case "analytics_dashboard":
		tables = ddm.designAnalyticsDashboardSchema()
	case "recommendation_engine":
		tables = ddm.designRecommendationSchema()
	default:
		return SchemaDesign{}, fmt.Errorf("unsupported system type for denormalization: %s", requirements.SystemType)
	}

	schema := SchemaDesign{
		DatabaseType:  "MongoDB", // Document-based for flexibility
		Tables:        tables,
		Indexes:       []IndexDesign{},
		Constraints:   []ConstraintDesign{}, // Fewer constraints in denormalized approach
		ModelingType:  "Denormalized",
		Justification: "Read-heavy workload, eventual consistency acceptable, performance over storage cost",
	}

	return schema, nil
}

func (ddm *DenormalizedDataModeling) designSocialMediaFeedSchema() []TableDesign {
	return []TableDesign{
		{
			Name: "user_timeline_posts",
			Columns: []ColumnDesign{
				{Name: "timeline_id", Type: "ObjectId", PrimaryKey: true},
				{Name: "user_id", Type: "BIGINT", NotNull: true},
				{Name: "post_id", Type: "BIGINT", NotNull: true},
				{Name: "post_content", Type: "TEXT", NotNull: true}, // Denormalized
				{Name: "author_id", Type: "BIGINT", NotNull: true},
				{Name: "author_username", Type: "VARCHAR(50)", NotNull: true}, // Denormalized
				{Name: "author_profile_pic", Type: "VARCHAR(500)"},            // Denormalized
				{Name: "like_count", Type: "INTEGER", DefaultValue: "0"},      // Denormalized
				{Name: "comment_count", Type: "INTEGER", DefaultValue: "0"},   // Denormalized
				{Name: "media_urls", Type: "JSON"},                            // Array of media
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true},
				{Name: "updated_at", Type: "TIMESTAMP", NotNull: true},
			},
			EstimatedRowCount: 10000000000, // Much higher due to duplication
		},
		{
			Name: "user_feed_cache",
			Columns: []ColumnDesign{
				{Name: "user_id", Type: "BIGINT", PrimaryKey: true},
				{Name: "feed_posts", Type: "JSON", NotNull: true}, // Pre-computed feed
				{Name: "last_updated", Type: "TIMESTAMP", NotNull: true},
			},
			EstimatedRowCount: 100000000,
		},
		{
			Name: "trending_posts",
			Columns: []ColumnDesign{
				{Name: "trend_id", Type: "ObjectId", PrimaryKey: true},
				{Name: "post_id", Type: "BIGINT", NotNull: true},
				{Name: "post_content", Type: "TEXT", NotNull: true},           // Denormalized
				{Name: "author_username", Type: "VARCHAR(50)", NotNull: true}, // Denormalized
				{Name: "engagement_score", Type: "FLOAT", NotNull: true},
				{Name: "trend_category", Type: "VARCHAR(50)", NotNull: true},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true},
			},
			EstimatedRowCount: 1000000,
		},
	}
}

func (ddm *DenormalizedDataModeling) designAnalyticsDashboardSchema() []TableDesign {
	return []TableDesign{
		{
			Name: "daily_user_metrics",
			Columns: []ColumnDesign{
				{Name: "metric_id", Type: "ObjectId", PrimaryKey: true},
				{Name: "date", Type: "DATE", NotNull: true},
				{Name: "user_id", Type: "BIGINT", NotNull: true},
				{Name: "username", Type: "VARCHAR(50)", NotNull: true},     // Denormalized
				{Name: "user_segment", Type: "VARCHAR(50)", NotNull: true}, // Denormalized
				{Name: "page_views", Type: "INTEGER", DefaultValue: "0"},
				{Name: "session_duration", Type: "INTEGER", DefaultValue: "0"},
				{Name: "clicks", Type: "INTEGER", DefaultValue: "0"},
				{Name: "purchases", Type: "INTEGER", DefaultValue: "0"},
				{Name: "revenue", Type: "DECIMAL(10,2)", DefaultValue: "0"},
			},
			EstimatedRowCount: 3650000000, // 10M users * 365 days
		},
		{
			Name: "hourly_system_metrics",
			Columns: []ColumnDesign{
				{Name: "metric_id", Type: "ObjectId", PrimaryKey: true},
				{Name: "timestamp", Type: "TIMESTAMP", NotNull: true},
				{Name: "active_users", Type: "INTEGER", NotNull: true},
				{Name: "requests_per_second", Type: "INTEGER", NotNull: true},
				{Name: "error_rate", Type: "FLOAT", NotNull: true},
				{Name: "avg_response_time", Type: "FLOAT", NotNull: true},
				{Name: "cpu_usage", Type: "FLOAT", NotNull: true},
				{Name: "memory_usage", Type: "FLOAT", NotNull: true},
			},
			EstimatedRowCount: 8760000, // 24 * 365 * years
		},
	}
}

func (ddm *DenormalizedDataModeling) designRecommendationSchema() []TableDesign {
	return []TableDesign{
		{
			Name: "user_recommendations",
			Columns: []ColumnDesign{
				{Name: "recommendation_id", Type: "ObjectId", PrimaryKey: true},
				{Name: "user_id", Type: "BIGINT", NotNull: true},
				{Name: "item_id", Type: "BIGINT", NotNull: true},
				{Name: "item_title", Type: "VARCHAR(255)", NotNull: true},    // Denormalized
				{Name: "item_category", Type: "VARCHAR(100)", NotNull: true}, // Denormalized
				{Name: "item_price", Type: "DECIMAL(10,2)"},                  // Denormalized
				{Name: "recommendation_score", Type: "FLOAT", NotNull: true},
				{Name: "recommendation_reason", Type: "VARCHAR(255)"},
				{Name: "created_at", Type: "TIMESTAMP", NotNull: true},
			},
			EstimatedRowCount: 1000000000,
		},
	}
}

// Indexing Strategies
// Interview Focus: Query optimization and performance tuning

type IndexStrategyAnalyzer struct {
	queryPatterns []QueryPattern
	cardinalities map[string]int64
	selectivities map[string]float64
}

func NewIndexStrategyAnalyzer() *IndexStrategyAnalyzer {
	return &IndexStrategyAnalyzer{
		cardinalities: make(map[string]int64),
		selectivities: make(map[string]float64),
	}
}

// OptimizeForQueries generates index recommendations based on query patterns
// FAANG Interview Point: Performance optimization, cost-benefit analysis
func (isa *IndexStrategyAnalyzer) OptimizeForQueries(queries []QueryPattern) ([]IndexRecommendation, error) {
	var recommendations []IndexRecommendation

	// Analyze each query pattern
	for _, query := range queries {
		switch query.Type {
		case SelectQuery:
			recs := isa.optimizeSelectQueries(query)
			recommendations = append(recommendations, recs...)
		case UpdateQuery, DeleteQuery:
			recs := isa.optimizeWriteQueries(query)
			recommendations = append(recommendations, recs...)
		case AggregateQuery:
			recs := isa.optimizeAggregateQueries(query)
			recommendations = append(recommendations, recs...)
		}
	}

	// Remove duplicates and rank by impact
	recommendations = isa.deduplicateAndRank(recommendations)

	return recommendations, nil
}

func (isa *IndexStrategyAnalyzer) optimizeSelectQueries(query QueryPattern) []IndexRecommendation {
	var recommendations []IndexRecommendation

	// Single column indexes for high-frequency equality filters
	for _, column := range query.Columns {
		for _, filterType := range query.FilterTypes {
			if filterType == EqualsFilter && query.Frequency >= High {
				recommendations = append(recommendations, IndexRecommendation{
					Type:            "B-TREE",
					Columns:         []string{column},
					Justification:   fmt.Sprintf("High-frequency equality filter on %s", column),
					ExpectedGain:    "50-90% query time reduction",
					StorageOverhead: "5-15% of table size",
					Priority:        "HIGH",
				})
			}
		}
	}

	// Composite indexes for multi-column filters
	if len(query.Columns) > 1 && query.Frequency >= Medium {
		recommendations = append(recommendations, IndexRecommendation{
			Type:            "COMPOSITE",
			Columns:         query.Columns,
			Justification:   "Multi-column filtering pattern detected",
			ExpectedGain:    "30-70% query time reduction",
			StorageOverhead: "10-25% of table size",
			Priority:        "MEDIUM",
		})
	}

	// Full-text search indexes
	for _, filterType := range query.FilterTypes {
		if filterType == FullTextFilter {
			recommendations = append(recommendations, IndexRecommendation{
				Type:            "GIN",
				Columns:         query.Columns,
				Justification:   "Full-text search requirements",
				ExpectedGain:    "100x improvement for text search",
				StorageOverhead: "20-50% of text column size",
				Priority:        "HIGH",
			})
		}
	}

	return recommendations
}

func (isa *IndexStrategyAnalyzer) optimizeWriteQueries(query QueryPattern) []IndexRecommendation {
	var recommendations []IndexRecommendation

	// Be conservative with indexes on write-heavy tables
	if query.Frequency >= High {
		recommendations = append(recommendations, IndexRecommendation{
			Type:            "PARTIAL",
			Columns:         query.Columns,
			Justification:   "Write-heavy table - partial index to minimize overhead",
			ExpectedGain:    "Reduced index maintenance cost",
			StorageOverhead: "2-10% of table size",
			Priority:        "LOW",
		})
	}

	return recommendations
}

func (isa *IndexStrategyAnalyzer) optimizeAggregateQueries(query QueryPattern) []IndexRecommendation {
	var recommendations []IndexRecommendation

	// Covering indexes for aggregation queries
	recommendations = append(recommendations, IndexRecommendation{
		Type:            "COVERING",
		Columns:         query.Columns,
		Justification:   "Covering index for aggregate queries - avoids table lookup",
		ExpectedGain:    "80-95% I/O reduction",
		StorageOverhead: "15-30% of table size",
		Priority:        "HIGH",
	})

	return recommendations
}

func (isa *IndexStrategyAnalyzer) deduplicateAndRank(recommendations []IndexRecommendation) []IndexRecommendation {
	// Simple deduplication by columns (in real implementation, would be more sophisticated)
	seen := make(map[string]bool)
	var unique []IndexRecommendation

	for _, rec := range recommendations {
		key := strings.Join(rec.Columns, ",")
		if !seen[key] {
			seen[key] = true
			unique = append(unique, rec)
		}
	}

	// Sort by priority (HIGH > MEDIUM > LOW)
	sort.Slice(unique, func(i, j int) bool {
		priorityOrder := map[string]int{"HIGH": 3, "MEDIUM": 2, "LOW": 1}
		return priorityOrder[unique[i].Priority] > priorityOrder[unique[j].Priority]
	})

	return unique
}

// Partitioning Schemes
// Interview Focus: Horizontal scaling and data lifecycle management

type PartitioningStrategy interface {
	DesignPartitions(table TableDesign, scale ScaleExpectation) (PartitioningScheme, error)
	EstimatePartitionSize(scheme PartitioningScheme, timeFrame time.Duration) PartitionSizeEstimate
	GeneratePartitionMaintenance(scheme PartitioningScheme) []MaintenanceTask
}

// TimeBasedPartitioning implements time-series partitioning
// FAANG Interview Point: Logs, metrics, financial transactions
type TimeBasedPartitioning struct {
	partitionInterval time.Duration
	retentionPeriod   time.Duration
}

func NewTimeBasedPartitioning(interval, retention time.Duration) *TimeBasedPartitioning {
	return &TimeBasedPartitioning{
		partitionInterval: interval,
		retentionPeriod:   retention,
	}
}

func (tbp *TimeBasedPartitioning) DesignPartitions(table TableDesign, scale ScaleExpectation) (PartitioningScheme, error) {
	// Determine optimal partition interval based on data volume
	var interval time.Duration
	if scale.RecordsPerDay > 100000000 { // 100M+ records per day
		interval = 24 * time.Hour // Daily partitions
	} else if scale.RecordsPerDay > 10000000 { // 10M+ records per day
		interval = 7 * 24 * time.Hour // Weekly partitions
	} else {
		interval = 30 * 24 * time.Hour // Monthly partitions
	}

	scheme := PartitioningScheme{
		Type:                "TIME_BASED",
		PartitionKey:        "created_at", // Assuming timestamp column
		PartitionInterval:   interval,
		RetentionPeriod:     scale.DataRetention,
		EstimatedPartitions: int64(scale.DataRetention / interval),
	}

	// Generate partition definitions
	partitionCount := int(scale.DataRetention / interval)
	for i := 0; i < partitionCount; i++ {
		startTime := time.Now().Add(-time.Duration(i) * interval)
		endTime := startTime.Add(interval)

		partition := PartitionDefinition{
			Name:          fmt.Sprintf("%s_%s", table.Name, startTime.Format("2006_01_02")),
			StartRange:    startTime.Format("2006-01-02"),
			EndRange:      endTime.Format("2006-01-02"),
			EstimatedSize: scale.RecordsPerDay * int64(interval.Hours()) / 24,
		}
		scheme.Partitions = append(scheme.Partitions, partition)
	}

	return scheme, nil
}

func (tbp *TimeBasedPartitioning) EstimatePartitionSize(scheme PartitioningScheme, timeFrame time.Duration) PartitionSizeEstimate {
	totalPartitions := int64(timeFrame / scheme.PartitionInterval)

	return PartitionSizeEstimate{
		TotalPartitions:      totalPartitions,
		AveragePartitionSize: 1000000,              // 1M records average
		TotalStorageGB:       totalPartitions * 10, // 10GB per partition
		QueryPerformance:     "90% improvement for time-range queries",
	}
}

func (tbp *TimeBasedPartitioning) GeneratePartitionMaintenance(scheme PartitioningScheme) []MaintenanceTask {
	return []MaintenanceTask{
		{
			Type:        "CREATE_PARTITION",
			Description: "Create new partition for upcoming time period",
			Schedule:    "Daily at 00:00 UTC",
			SQL:         fmt.Sprintf("CREATE TABLE %s_new PARTITION OF %s FOR VALUES FROM ('%s') TO ('%s')", "table", "table", "start", "end"),
		},
		{
			Type:        "DROP_PARTITION",
			Description: "Drop partitions beyond retention period",
			Schedule:    "Weekly at 02:00 UTC",
			SQL:         "DROP TABLE old_partition",
		},
		{
			Type:        "ANALYZE_PARTITIONS",
			Description: "Update partition statistics",
			Schedule:    "Daily at 06:00 UTC",
			SQL:         "ANALYZE partition_table",
		},
	}
}

// HashBasedPartitioning implements hash-based horizontal partitioning
// FAANG Interview Point: User data, distributed systems
type HashBasedPartitioning struct {
	hashFunction   string
	partitionCount int
}

func NewHashBasedPartitioning(partitionCount int) *HashBasedPartitioning {
	return &HashBasedPartitioning{
		hashFunction:   "CRC32",
		partitionCount: partitionCount,
	}
}

func (hbp *HashBasedPartitioning) DesignPartitions(table TableDesign, scale ScaleExpectation) (PartitioningScheme, error) {
	scheme := PartitioningScheme{
		Type:                "HASH_BASED",
		PartitionKey:        "user_id", // Common partition key
		PartitionCount:      hbp.partitionCount,
		EstimatedPartitions: int64(hbp.partitionCount),
	}

	recordsPerPartition := scale.RecordsPerDay / int64(hbp.partitionCount)

	for i := 0; i < hbp.partitionCount; i++ {
		partition := PartitionDefinition{
			Name:          fmt.Sprintf("%s_p%d", table.Name, i),
			HashModulus:   hbp.partitionCount,
			HashRemainder: i,
			EstimatedSize: recordsPerPartition,
		}
		scheme.Partitions = append(scheme.Partitions, partition)
	}

	return scheme, nil
}

func (hbp *HashBasedPartitioning) EstimatePartitionSize(scheme PartitioningScheme, timeFrame time.Duration) PartitionSizeEstimate {
	return PartitionSizeEstimate{
		TotalPartitions:      int64(hbp.partitionCount),
		AveragePartitionSize: 10000000 / int64(hbp.partitionCount), // Evenly distributed
		TotalStorageGB:       100,                                  // Total across all partitions
		QueryPerformance:     "Single partition access for point queries",
	}
}

func (hbp *HashBasedPartitioning) GeneratePartitionMaintenance(scheme PartitioningScheme) []MaintenanceTask {
	return []MaintenanceTask{
		{
			Type:        "REBALANCE_CHECK",
			Description: "Check for partition size imbalances",
			Schedule:    "Weekly",
			SQL:         "SELECT partition_name, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) FROM pg_tables",
		},
		{
			Type:        "PARTITION_STATISTICS",
			Description: "Update partition-level statistics",
			Schedule:    "Daily",
			SQL:         "ANALYZE ALL PARTITIONS",
		},
	}
}

// Data Archiving Patterns
// Interview Focus: Cost optimization and compliance

type DataArchivingStrategy struct {
	hotDataPeriod    time.Duration
	warmDataPeriod   time.Duration
	coldDataPeriod   time.Duration
	compressionRatio float64
}

func NewDataArchivingStrategy() *DataArchivingStrategy {
	return &DataArchivingStrategy{
		hotDataPeriod:    90 * 24 * time.Hour,      // 90 days
		warmDataPeriod:   365 * 24 * time.Hour,     // 1 year
		coldDataPeriod:   7 * 365 * 24 * time.Hour, // 7 years
		compressionRatio: 0.3,                      // 70% compression
	}
}

// DesignArchivingPlan creates data lifecycle management plan
// FAANG Interview Point: Compliance, cost optimization, query performance
func (das *DataArchivingStrategy) DesignArchivingPlan(table TableDesign, requirements SchemaRequirements) ArchivingPlan {
	plan := ArchivingPlan{
		TableName: table.Name,
		Tiers:     []DataTier{},
	}

	// Hot tier - frequently accessed data
	hotTier := DataTier{
		Name:             "HOT",
		StorageType:      "SSD",
		AccessPattern:    "Real-time queries",
		RetentionPeriod:  das.hotDataPeriod,
		CompressionRatio: 1.0,  // No compression for performance
		CostPerGB:        0.10, // High cost, high performance
		QueryLatency:     "< 100ms",
	}

	// Warm tier - occasionally accessed data
	warmTier := DataTier{
		Name:             "WARM",
		StorageType:      "Standard HDD",
		AccessPattern:    "Analytics queries",
		RetentionPeriod:  das.warmDataPeriod,
		CompressionRatio: 0.5,  // 50% compression
		CostPerGB:        0.03, // Medium cost
		QueryLatency:     "< 1s",
	}

	// Cold tier - rarely accessed data
	coldTier := DataTier{
		Name:             "COLD",
		StorageType:      "Archive Storage",
		AccessPattern:    "Compliance queries",
		RetentionPeriod:  das.coldDataPeriod,
		CompressionRatio: das.compressionRatio,
		CostPerGB:        0.01, // Low cost
		QueryLatency:     "< 10s",
	}

	plan.Tiers = []DataTier{hotTier, warmTier, coldTier}

	// Generate migration rules
	plan.MigrationRules = []MigrationRule{
		{
			FromTier:  "HOT",
			ToTier:    "WARM",
			Condition: fmt.Sprintf("created_at < NOW() - INTERVAL '%d days'", int(das.hotDataPeriod.Hours()/24)),
			Schedule:  "Daily at 02:00 UTC",
		},
		{
			FromTier:  "WARM",
			ToTier:    "COLD",
			Condition: fmt.Sprintf("created_at < NOW() - INTERVAL '%d days'", int(das.warmDataPeriod.Hours()/24)),
			Schedule:  "Weekly at 03:00 UTC",
		},
	}

	return plan
}

// Analytics Optimizations
// Interview Focus: OLAP vs OLTP separation, data warehousing

type AnalyticsOptimizer struct {
	separateOLAP      bool
	materializedViews bool
	columnStore       bool
}

func NewAnalyticsOptimizer() *AnalyticsOptimizer {
	return &AnalyticsOptimizer{
		separateOLAP:      true,
		materializedViews: true,
		columnStore:       true,
	}
}

// OptimizeForAnalytics creates analytics-optimized schema
// FAANG Interview Point: Data warehousing, business intelligence
func (ao *AnalyticsOptimizer) OptimizeForAnalytics(transactionalSchema SchemaDesign) AnalyticsSchema {
	schema := AnalyticsSchema{
		SourceSchema:      transactionalSchema,
		OptimizationType:  "COLUMN_STORE",
		MaterializedViews: []MaterializedView{},
		AggregationTables: []AggregationTable{},
	}

	// Create materialized views for common analytics queries
	schema.MaterializedViews = []MaterializedView{
		{
			Name: "daily_sales_summary",
			Definition: `
				SELECT 
					DATE(created_at) as sale_date,
					COUNT(*) as order_count,
					SUM(total_amount) as total_revenue,
					AVG(total_amount) as avg_order_value
				FROM orders 
				GROUP BY DATE(created_at)
			`,
			RefreshSchedule: "Every 1 hour",
			EstimatedSize:   "1GB per year",
		},
		{
			Name: "user_behavior_metrics",
			Definition: `
				SELECT 
					user_id,
					COUNT(DISTINCT DATE(created_at)) as active_days,
					COUNT(*) as total_orders,
					SUM(total_amount) as lifetime_value,
					FIRST_VALUE(created_at) OVER (PARTITION BY user_id ORDER BY created_at) as first_order_date
				FROM orders 
				GROUP BY user_id
			`,
			RefreshSchedule: "Daily at 01:00 UTC",
			EstimatedSize:   "500MB",
		},
	}

	// Create pre-aggregated tables for faster reporting
	schema.AggregationTables = []AggregationTable{
		{
			Name:          "monthly_revenue_by_category",
			Granularity:   "MONTHLY",
			Dimensions:    []string{"category", "region"},
			Metrics:       []string{"revenue", "order_count", "unique_customers"},
			EstimatedSize: "10MB per month",
		},
		{
			Name:          "hourly_user_activity",
			Granularity:   "HOURLY",
			Dimensions:    []string{"user_segment", "device_type"},
			Metrics:       []string{"page_views", "session_count", "conversion_rate"},
			EstimatedSize: "1GB per month",
		},
	}

	return schema
}

// Data structure definitions
type SchemaDesign struct {
	DatabaseType  string
	Tables        []TableDesign
	Indexes       []IndexDesign
	Constraints   []ConstraintDesign
	ModelingType  string
	Justification string
}

type TableDesign struct {
	Name              string
	Columns           []ColumnDesign
	EstimatedRowCount int64
}

type ColumnDesign struct {
	Name         string
	Type         string
	PrimaryKey   bool
	ForeignKey   string
	Unique       bool
	NotNull      bool
	DefaultValue string
}

type IndexDesign struct {
	Name    string
	Table   string
	Columns []string
	Type    string
	Unique  bool
}

type ConstraintDesign struct {
	Type             string
	SourceTable      string
	SourceColumn     string
	ReferencedTable  string
	ReferencedColumn string
}

type IndexRecommendation struct {
	Type            string
	Columns         []string
	Justification   string
	ExpectedGain    string
	StorageOverhead string
	Priority        string
}

type NormalizationAnalysis struct {
	CurrentLevel     int
	RecommendedLevel int
	Issues           []string
	Benefits         []string
	TradeOffs        []string
}

type StorageEstimate struct {
	TotalSizeGB      float64
	IndexSizeGB      float64
	CompressionRatio float64
	GrowthRate       float64
}

type PartitioningScheme struct {
	Type                string
	PartitionKey        string
	PartitionInterval   time.Duration
	PartitionCount      int
	RetentionPeriod     time.Duration
	EstimatedPartitions int64
	Partitions          []PartitionDefinition
}

type PartitionDefinition struct {
	Name          string
	StartRange    string
	EndRange      string
	HashModulus   int
	HashRemainder int
	EstimatedSize int64
}

type PartitionSizeEstimate struct {
	TotalPartitions      int64
	AveragePartitionSize int64
	TotalStorageGB       int64
	QueryPerformance     string
}

type MaintenanceTask struct {
	Type        string
	Description string
	Schedule    string
	SQL         string
}

type ArchivingPlan struct {
	TableName      string
	Tiers          []DataTier
	MigrationRules []MigrationRule
}

type DataTier struct {
	Name             string
	StorageType      string
	AccessPattern    string
	RetentionPeriod  time.Duration
	CompressionRatio float64
	CostPerGB        float64
	QueryLatency     string
}

type MigrationRule struct {
	FromTier  string
	ToTier    string
	Condition string
	Schedule  string
}

type AnalyticsSchema struct {
	SourceSchema      SchemaDesign
	OptimizationType  string
	MaterializedViews []MaterializedView
	AggregationTables []AggregationTable
}

type MaterializedView struct {
	Name            string
	Definition      string
	RefreshSchedule string
	EstimatedSize   string
}

type AggregationTable struct {
	Name          string
	Granularity   string
	Dimensions    []string
	Metrics       []string
	EstimatedSize string
}

// Validation and Analysis Methods
// Interview Focus: Systematic evaluation of design decisions

func (ndm *NormalizedDataModeling) ValidateNormalization(schema SchemaDesign) NormalizationAnalysis {
	analysis := NormalizationAnalysis{
		CurrentLevel:     3, // Assume 3NF
		RecommendedLevel: 3,
		Issues:           []string{},
		Benefits:         []string{},
		TradeOffs:        []string{},
	}

	analysis.Benefits = []string{
		"Strong data consistency",
		"Minimal data redundancy",
		"Easier to maintain data integrity",
		"Efficient storage utilization",
	}

	analysis.TradeOffs = []string{
		"Complex queries require multiple JOINs",
		"Higher query latency for read operations",
		"Increased application complexity",
		"Potential for N+1 query problems",
	}

	return analysis
}

func (ddm *DenormalizedDataModeling) ValidateNormalization(schema SchemaDesign) NormalizationAnalysis {
	analysis := NormalizationAnalysis{
		CurrentLevel:     1, // Heavily denormalized
		RecommendedLevel: 1,
		Issues:           []string{},
		Benefits:         []string{},
		TradeOffs:        []string{},
	}

	analysis.Benefits = []string{
		"Fast read queries - single table access",
		"Reduced JOIN complexity",
		"Better cache locality",
		"Simplified application logic for reads",
	}

	analysis.TradeOffs = []string{
		"Data redundancy and inconsistency risk",
		"Increased storage costs",
		"Complex update operations",
		"Higher maintenance overhead",
	}

	return analysis
}

func (ndm *NormalizedDataModeling) EstimateStorage(schema SchemaDesign, recordCount int64) StorageEstimate {
	var totalSize float64

	// Calculate base table sizes
	for _, table := range schema.Tables {
		avgRowSize := ndm.estimateRowSize(table.Columns)
		tableSize := float64(table.EstimatedRowCount*avgRowSize) / (1024 * 1024 * 1024) // Convert to GB
		totalSize += tableSize
	}

	// Add index overhead (typically 20-30% for normalized schemas)
	indexOverhead := totalSize * 0.25

	return StorageEstimate{
		TotalSizeGB:      totalSize,
		IndexSizeGB:      indexOverhead,
		CompressionRatio: 0.7, // 30% compression
		GrowthRate:       0.1, // 10% monthly growth
	}
}

func (ddm *DenormalizedDataModeling) EstimateStorage(schema SchemaDesign, recordCount int64) StorageEstimate {
	var totalSize float64

	// Calculate base table sizes (larger due to denormalization)
	for _, table := range schema.Tables {
		avgRowSize := ddm.estimateRowSize(table.Columns) * 2 // 2x due to redundancy
		tableSize := float64(table.EstimatedRowCount*avgRowSize) / (1024 * 1024 * 1024)
		totalSize += tableSize
	}

	// Less index overhead but larger base size
	indexOverhead := totalSize * 0.15

	return StorageEstimate{
		TotalSizeGB:      totalSize,
		IndexSizeGB:      indexOverhead,
		CompressionRatio: 0.8,  // Better compression due to redundancy
		GrowthRate:       0.15, // Higher growth rate due to duplication
	}
}

func (ndm *NormalizedDataModeling) estimateRowSize(columns []ColumnDesign) int64 {
	var size int64

	for _, col := range columns {
		switch {
		case strings.Contains(col.Type, "BIGINT"):
			size += 8
		case strings.Contains(col.Type, "INTEGER"):
			size += 4
		case strings.Contains(col.Type, "VARCHAR"):
			// Extract length or use default
			size += 50 // Average VARCHAR size
		case strings.Contains(col.Type, "TEXT"):
			size += 200 // Average TEXT size
		case strings.Contains(col.Type, "TIMESTAMP"):
			size += 8
		case strings.Contains(col.Type, "DECIMAL"):
			size += 16
		default:
			size += 10 // Default estimate
		}
	}

	return size
}

func (ddm *DenormalizedDataModeling) estimateRowSize(columns []ColumnDesign) int64 {
	// Similar to normalized but account for JSON and array types
	var size int64

	for _, col := range columns {
		switch {
		case strings.Contains(col.Type, "JSON"):
			size += 500 // JSON can be large
		case strings.Contains(col.Type, "TEXT"):
			size += 300 // Larger TEXT fields in denormalized
		default:
			size += 10 // Base estimate
		}
	}

	return size
}
