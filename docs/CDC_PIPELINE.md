# TIA CDC Pipeline & Knowledge Graph

This document describes the Change Data Capture (CDC) pipeline that creates a knowledge graph from your MySQL/MariaDB database and provides intelligent connection recommendations.

## 🏗️ Architecture Overview

```
MariaDB → Debezium → Kafka → Knowledge Graph Builder → Memgraph → Connection Analyzer → API Endpoints
```

### Components

1. **MariaDB**: Source database with binlog enabled
2. **Debezium**: CDC connector that captures database changes
3. **Kafka**: Message streaming platform
4. **Knowledge Graph Builder**: Service that processes CDC events and builds the graph
5. **Memgraph**: Graph database storing the knowledge graph
6. **Connection Analyzer**: Service that analyzes the graph and provides recommendations
7. **API Endpoints**: REST endpoints integrated with your main API

## Quick Start

### 1. Setup CDC Pipeline

```bash
# Start the CDC pipeline
./scripts/setup-cdc-pipeline.sh
```

### 2. Generate Test Data

```bash
# Generate realistic test data
./scripts/generate_data.sh
```

### 3. Test Connection Endpoints

```bash
# Get complementary partners for user ID 1
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/connections/complementary/1

# Get alliance partners for user ID 1
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/connections/alliance/1

# Get mastermind partners for user ID 1
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/connections/mastermind/1

# Get all recommendations
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/connections/recommendations/1

# Get connection analysis
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/connections/analysis/1
```

## Connection Types

### 1. COMPLEMENTARY PARTNERS

**Definition**: Businesses that offer products or services that align well with your own.

**Algorithm**:
- Analyzes business type compatibility
- Technology + Consulting = High compatibility
- Manufacturing + Services = Medium compatibility
- Considers business phase alignment

**Example Response**:
```json
{
  "type": "COMPLEMENTARY_PARTNERS",
  "count": 5,
  "recommendations": [
    {
      "type": "COMPLEMENTARY_PARTNER",
      "score": 0.9,
      "reason": "Complementary business types: Tech + Consulting",
      "business": {
        "id": 2,
        "name": "Creative Design Studio",
        "businessType": "Services",
        "businessCategory": "B2B",
        "businessPhase": "Mature"
      },
      "user": {
        "id": 2,
        "firstName": "Sarah",
        "lastName": "Johnson",
        "email": "sarah.johnson@example.com"
      },
      "compatibilityFactors": [
        "Business type compatibility: 90%",
        "Current connection status: not_connected"
      ]
    }
  ]
}
```

### 2. ALLIANCE PARTNERS

**Definition**: Businesses you can collaborate with on projects.

**Algorithm**:
- Analyzes shared skills between users
- Considers project collaboration potential
- Weights management and strategic skills higher
- Looks for active projects in both businesses

**Example Response**:
```json
{
  "type": "ALLIANCE_PARTNERS",
  "count": 3,
  "recommendations": [
    {
      "type": "ALLIANCE_PARTNER",
      "score": 0.8,
      "reason": "Shared 4 skills with potential for project collaboration",
      "business": {
        "id": 3,
        "name": "Green Energy Co",
        "businessType": "Manufacturing",
        "businessPhase": "Startup"
      },
      "user": {
        "id": 3,
        "firstName": "Michael",
        "lastName": "Brown",
        "email": "michael.brown@example.com"
      },
      "skills": [
        {
          "skillId": 11,
          "skillName": "Project Management",
          "skillCategory": "Business",
          "userProficiency": "intermediate",
          "matchScore": 1.0
        }
      ],
      "compatibilityFactors": [
        "Shared skills: 4",
        "Alliance potential: 80%",
        "Current connection status: not_connected"
      ]
    }
  ]
}
```

### 3. MASTERMIND PARTNERS

**Definition**: People who complement your skills through direct collaboration or supportive accountability.

**Algorithm**:
- Analyzes complementary skill sets
- Considers business phase progression
- Startup → Growth/Mature = High compatibility
- Weights business and strategic skills

**Example Response**:
```json
{
  "type": "MASTERMIND_PARTNERS",
  "count": 2,
  "recommendations": [
    {
      "type": "MASTERMIND_PARTNER",
      "score": 0.9,
      "reason": "Complementary business phase (Startup → Growth) and skill sets",
      "business": {
        "id": 4,
        "name": "HealthTech Innovations",
        "businessType": "Technology",
        "businessPhase": "Growth"
      },
      "user": {
        "id": 4,
        "firstName": "Emily",
        "lastName": "Davis",
        "email": "emily.davis@example.com"
      },
      "compatibilityFactors": [
        "Phase compatibility: 90%",
        "Skill complementarity: 70%",
        "Current connection status: not_connected"
      ]
    }
  ]
}
```

## Configuration

### Environment Variables

```bash
# Database
DATABASE_URL="root:tia-dev-password@tcp(127.0.0.1:3306)/tia-dev?charset=utf8mb4&parseTime=True&loc=Local"
JWT_SECRET="your-secret-key-here"

# Kafka
KAFKA_BROKERS="kafka:29092"

# Memgraph
MEMGRAPH_HOST="memgraph"
MEMGRAPH_PORT="7687"

# Services
CONNECTION_ANALYZER_HOST="connection-analyzer"
CONNECTION_ANALYZER_PORT="8082"
```

### MariaDB Configuration

Ensure your MariaDB has binlog enabled:

```sql
-- In my.cnf or my.ini
[mysqld]
log-bin=mysql-bin
binlog-format=ROW
server-id=1
```

## 📈 Knowledge Graph Schema

### Nodes

- **User**: `{id, firstName, lastName, email, active, emailVerified}`
- **Business**: `{id, name, businessType, businessCategory, businessPhase, active}`
- **Project**: `{id, name, description, projectStatus, startDate, targetEndDate}`
- **Skill**: `{id, category, name, description, active}`

### Relationships

- **User** `-[:OPERATES]->` **Business**
- **User** `-[:MANAGES]->` **Project**
- **User** `-[:HAS_SKILL]->` **Skill**
- **Business** `-[:CONNECTS_TO]->` **Business**
- **Project** `-[:REQUIRES_SKILL]->` **Skill**
- **Business** `-[:HAS_PROJECT]->` **Project**

## Monitoring & Debugging

### Check Service Status

```bash
# Check all services
docker-compose -f docker-compose.cdc.yml ps

# Check Kafka topics
docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list

# Check Debezium connector status
curl http://localhost:8083/connectors/tia-mariadb-connector/status

# Check Memgraph
curl http://localhost:7444/status
```

### View Knowledge Graph

1. Open Memgraph Lab: http://localhost:7444
2. Connect to: `bolt://localhost:7687`
3. Run queries:

```cypher
// View all users and their businesses
MATCH (u:User)-[:OPERATES]->(b:Business)
RETURN u, b

// Find complementary business connections
MATCH (b1:Business)-[r:CONNECTS_TO]->(b2:Business)
WHERE r.connectionType = 'Partnership'
RETURN b1, r, b2

// Analyze skill distribution
MATCH (u:User)-[:HAS_SKILL]->(s:Skill)
RETURN s.category, count(s) as skill_count
ORDER BY skill_count DESC
```

### Debug CDC Events

```bash
# View Kafka messages
docker exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic tia-db.tia-dev.users \
  --from-beginning
```

## 🚨 Troubleshooting

### Common Issues

1. **Debezium Connector Fails**
   - Check MariaDB binlog configuration
   - Verify database user permissions
   - Check connector logs: `docker logs kafka-connect`

2. **Knowledge Graph Not Updating**
   - Check Kafka topics for messages
   - Verify Knowledge Graph Builder logs: `docker logs kg-builder`
   - Check Memgraph connection

3. **Connection Recommendations Empty**
   - Verify data exists in Memgraph
   - Check Connection Analyzer logs: `docker logs connection-analyzer`
   - Ensure user has associated business and skills

4. **API Endpoints Return Errors**
   - Check authentication tokens
   - Verify user permissions
   - Check Connection Analyzer service status

### Performance Optimization

1. **Kafka Tuning**
   - Increase `max.batch.size` and `max.queue.size`
   - Adjust `poll.interval.ms`

2. **Memgraph Optimization**
   - Create appropriate indexes
   - Use query optimization techniques

3. **Connection Analysis**
   - Cache frequent queries
   - Implement pagination for large result sets

## 🔮 Future Enhancements

### Planned Features

1. **Machine Learning Integration**
   - Train models on connection success rates
   - Improve recommendation accuracy
   - Add predictive analytics

2. **Real-time Notifications**
   - WebSocket support for live updates
   - Push notifications for new connections
   - Real-time collaboration features

3. **Advanced Analytics**
   - Network analysis metrics
   - Business growth predictions
   - Market opportunity identification

4. **Integration Enhancements**
   - CRM system integration
   - Social media data enrichment
   - External API connections

### API Extensions

```go
// Future endpoint examples
GET /api/v1/connections/trending
GET /api/v1/connections/analytics/network-strength
GET /api/v1/connections/predictions/growth-potential
POST /api/v1/connections/introductions/request
```

## 📚 Additional Resources

- [Debezium Documentation](https://debezium.io/documentation/)
- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [Memgraph Documentation](https://memgraph.com/docs/)
- [Cypher Query Language](https://neo4j.com/docs/cypher-manual/)
- [TIA API Documentation](http://localhost:8080/swagger/index.html)

## 🤝 Contributing

When contributing to the CDC pipeline:

1. Follow the established patterns
2. Add comprehensive tests
3. Update documentation
4. Consider performance implications
5. Test with realistic data volumes

## 📄 License

This CDC pipeline is part of the TIA API project and follows the same licensing terms.
