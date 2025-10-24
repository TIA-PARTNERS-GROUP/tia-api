# TIA CDC Pipeline & Knowledge Graph

This directory contains all components for the Change Data Capture (CDC) pipeline that creates a knowledge graph from the MySQL/MariaDB database and provides intelligent connection recommendations.

## 📁 Directory Structure

```
cdc/
├── configs/
│   └── docker-compose.cdc.yml    # CDC-specific Docker Compose (legacy)
├── kg-builder/
│   ├── main.go                   # Knowledge Graph Builder service
│   └── Dockerfile.kg-builder     # Dockerfile for KG Builder
├── connection-analyzer/
│   ├── main.go                   # Connection Analysis service
│   └── Dockerfile.connection-analyzer  # Dockerfile for Connection Analyzer
├── scripts/
│   ├── setup-debezium.sh         # Debezium connector setup
│   └── setup-cdc-pipeline.sh     # CDC pipeline setup (legacy)
└── README.md                     # This file
```

## Quick Start

### Option 1: Use the main setup script (Recommended)
```bash
# From the project root
./setup-cdc.sh
```

### Option 2: Use Docker Compose directly
```bash
# Start all services (API + CDC Pipeline)
docker-compose up -d

# Setup Debezium connector
./cdc/scripts/setup-debezium.sh
```

## 🏗️ Architecture

```
MariaDB → Debezium → Kafka → Knowledge Graph Builder → Memgraph → Connection Analyzer → API Endpoints
```

### Components

1. **Knowledge Graph Builder** (`cdc/kg-builder/`)
   - Consumes CDC events from Kafka
   - Builds and maintains the knowledge graph in Memgraph
   - Handles real-time updates to the graph

2. **Connection Analyzer** (`cdc/connection-analyzer/`)
   - Analyzes the knowledge graph
   - Provides intelligent connection recommendations
   - Implements three connection types:
     - Complementary Partners
     - Alliance Partners
     - Mastermind Partners

3. **Configuration** (`cdc/configs/`)
   - Docker Compose configuration for CDC services
   - Service definitions and environment variables

4. **Scripts** (`cdc/scripts/`)
   - Setup and deployment scripts
   - Debezium connector configuration

## Services

### Knowledge Graph Builder
- **Port**: Internal service
- **Dependencies**: Kafka, Memgraph
- **Function**: Processes CDC events and builds graph

### Connection Analyzer
- **Port**: 8082
- **Dependencies**: Memgraph, Main API
- **Function**: Provides connection recommendations

## Connection Types

### 1. COMPLEMENTARY PARTNERS
Businesses that offer products or services that align well with your own.

**Algorithm**: Business type compatibility analysis
- Technology + Consulting = High compatibility
- Manufacturing + Services = Medium compatibility

### 2. ALLIANCE PARTNERS
Businesses you can collaborate with on projects.

**Algorithm**: Shared skills + project collaboration potential
- Analyzes skill overlap between users
- Considers active projects for collaboration

### 3. MASTERMIND PARTNERS
People who complement your skills through direct collaboration or supportive accountability.

**Algorithm**: Complementary skills + business phase progression
- Startup → Growth/Mature = High compatibility
- Analyzes skill complementarity

## API Endpoints

The connection analyzer provides these endpoints:

- `GET /api/v1/connections/complementary/{userId}` - Complementary partners
- `GET /api/v1/connections/alliance/{userId}` - Alliance partners
- `GET /api/v1/connections/mastermind/{userId}` - Mastermind partners
- `GET /api/v1/connections/recommendations/{userId}` - All recommendations
- `GET /api/v1/connections/analysis/{userId}` - Connection analysis

## 🛠️ Development

### Building Services
```bash
# Build KG Builder
docker build -f cdc/kg-builder/Dockerfile.kg-builder -t tia-kg-builder .

# Build Connection Analyzer
docker build -f cdc/connection-analyzer/Dockerfile.connection-analyzer -t tia-connection-analyzer .
```

### Running Locally
```bash
# KG Builder
cd cdc/kg-builder
go run main.go

# Connection Analyzer
cd cdc/connection-analyzer
go run main.go
```

### Testing
```bash
# Test connection endpoints
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/connections/complementary/1
```

## Monitoring

### Check Service Status
```bash
# All services
docker-compose ps

# Specific services
docker logs kg-builder
docker logs connection-analyzer
docker logs memgraph
```

### View Knowledge Graph
1. Open Memgraph Lab: http://localhost:7444
2. Connect to: `bolt://localhost:7687`
3. Run Cypher queries to explore the graph

### Debug CDC Events
```bash
# View Kafka topics
docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list

# View messages
docker exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic tia-db.tia-dev.users \
  --from-beginning
```

## 🚨 Troubleshooting

### Common Issues

1. **Services not starting**
   - Check Docker and Docker Compose installation
   - Verify .env file exists and has correct values
   - Check port conflicts

2. **Knowledge Graph not updating**
   - Verify Debezium connector is running
   - Check Kafka topics for messages
   - Review KG Builder logs

3. **Connection recommendations empty**
   - Ensure data exists in the database
   - Check Memgraph for graph data
   - Verify Connection Analyzer service is running

4. **API endpoints returning errors**
   - Check authentication tokens
   - Verify user permissions
   - Review service logs

## 📚 Additional Resources

- [Main CDC Documentation](../docs/CDC_PIPELINE.md)
- [Debezium Documentation](https://debezium.io/documentation/)
- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [Memgraph Documentation](https://memgraph.com/docs/)
- [Cypher Query Language](https://neo4j.com/docs/cypher-manual/)

## 🤝 Contributing

When contributing to the CDC pipeline:

1. Follow the established directory structure
2. Add comprehensive tests
3. Update documentation
4. Consider performance implications
5. Test with realistic data volumes

