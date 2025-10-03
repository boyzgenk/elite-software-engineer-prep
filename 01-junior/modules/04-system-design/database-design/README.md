# Database Design for System Design Interviews

## Overview

Database design is a critical component in system design interviews at FAANG companies. This module covers essential database concepts, trade-offs, and patterns that demonstrate scalable system architecture knowledge.

## Key Topics Covered

### 1. SQL vs NoSQL Decision Framework
- **ACID Properties**: Atomicity, Consistency, Isolation, Durability
- **CAP Theorem**: Consistency, Availability, Partition Tolerance
- **Use Case Analysis**: When to choose relational vs document vs key-value stores
- **Performance Characteristics**: Throughput, latency, scalability patterns

### 2. Database Sharding Strategies
- **Horizontal Partitioning**: Distributing data across multiple database instances
- **Shard Key Selection**: Critical decisions for even data distribution
- **Sharding Patterns**: Hash-based, range-based, directory-based
- **Cross-Shard Operations**: Handling queries that span multiple shards
- **Rebalancing**: Dynamic shard redistribution as system scales

### 3. Replication Patterns
- **Master-Slave**: Single writer, multiple readers for read scaling
- **Master-Master**: Multi-writer configurations with conflict resolution
- **Read Replicas**: Geographic distribution and read load distribution
- **Consistency Models**: Strong, eventual, and tunable consistency
- **Failover Strategies**: Automatic failover and split-brain prevention

### 4. Data Modeling for Scale
- **Schema Design**: Normalization vs denormalization trade-offs
- **Indexing Strategies**: B-tree, hash, bitmap, and composite indexes
- **Partitioning**: Time-based, hash-based, and range partitioning
- **Data Lifecycle**: Hot, warm, cold data management
- **Analytics Optimization**: OLTP vs OLAP workload separation

## Interview Discussion Points

### Database Selection Criteria
1. **Data Structure**: Structured (SQL) vs Semi-structured (Document) vs Unstructured (Key-Value)
2. **Scalability Requirements**: Read-heavy, write-heavy, or balanced workloads
3. **Consistency Requirements**: Strong consistency vs eventual consistency
4. **Query Patterns**: Simple lookups vs complex joins and analytics
5. **Operational Complexity**: Team expertise and maintenance overhead

### Common FAANG System Examples
- **Social Media Feed**: User timeline generation with denormalized data
- **E-commerce Product Catalog**: Product search and inventory management
- **Chat/Messaging**: Message storage and real-time delivery
- **Analytics Platform**: Time-series data and aggregation patterns
- **Content Management**: Media storage and metadata management

### Performance Considerations
- **Read vs Write Optimization**: Different patterns for different workloads
- **Cache Integration**: Redis, Memcached for hot data
- **Connection Pooling**: Managing database connections at scale
- **Query Optimization**: Index usage and query plan analysis
- **Monitoring**: Key metrics for database health and performance

## System Design Integration

### Typical Database Layers in Large Systems
1. **Application Database**: Primary transactional data store
2. **Cache Layer**: In-memory data for fast access
3. **Search Index**: Elasticsearch or Solr for complex queries
4. **Analytics Store**: Data warehouse for reporting and ML
5. **Archive Storage**: Cold storage for compliance and backup

### Scaling Patterns
- **Vertical Scaling**: CPU, memory, storage upgrades (limited scalability)
- **Horizontal Scaling**: Adding more database instances
- **Functional Partitioning**: Separating different features into different databases
- **CQRS**: Command Query Responsibility Segregation for read/write optimization

## Common Interview Questions

### Design Questions
1. "How would you design a database schema for a Twitter-like social media platform?"
2. "Design a sharding strategy for a global e-commerce platform"
3. "How would you handle database failover in a mission-critical system?"
4. "Design a data model for a real-time analytics dashboard"

### Trade-off Questions
1. "SQL vs NoSQL for a chat application - discuss the trade-offs"
2. "When would you choose master-master vs master-slave replication?"
3. "Normalize vs denormalize data for a social media feed"
4. "Consistent hashing vs range-based sharding trade-offs"

### Scaling Questions
1. "Your database is at 80% capacity - what are your options?"
2. "How do you handle hot shards in a distributed database?"
3. "Design a strategy to migrate from a monolithic to sharded database"
4. "How do you maintain data consistency across microservices?"

## Best Practices for Interviews

### Problem Approach
1. **Clarify Requirements**: Understand scale, consistency needs, query patterns
2. **Start Simple**: Begin with a single database, then discuss scaling
3. **Identify Bottlenecks**: CPU, memory, I/O, network limitations
4. **Propose Solutions**: Multiple approaches with trade-offs
5. **Consider Operations**: Monitoring, backup, disaster recovery

### Communication Tips
- Draw diagrams showing data flow and relationships
- Quantify decisions with estimates (QPS, storage, latency)
- Discuss monitoring and alerting strategies
- Consider failure scenarios and recovery plans
- Mention real-world examples from major tech companies

## Files in This Module

- `sql_vs_nosql.go`: Comparison framework and decision matrix
- `sharding_strategies.go`: Horizontal partitioning patterns
- `replication_patterns.go`: Database replication strategies
- `data_modeling.go`: Schema design and optimization patterns

Each file contains production-ready Go code examples that demonstrate the concepts in a practical, interview-relevant context.