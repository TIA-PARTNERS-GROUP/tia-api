# TIA API

A comprehensive business networking and collaboration platform with real-time Change Data Capture (CDC) pipeline and intelligent connection recommendations.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [API Endpoints](#api-endpoints)
- [Development](#development)
- [Testing](#testing)
- [Documentation](#documentation)
- [CDC Pipeline & Knowledge Graph](#cdc-pipeline--knowledge-graph)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before you begin, make sure you have the following tools installed:

* **[Go](https://go.dev/doc/install)** (v1.20+)
* **[Docker](https://docs.docker.com/engine/install/)**
* **[Docker Compose](https://docs.docker.com/compose/install/)**

## Quick Start

### 1. Environment Setup

Create your `.env` file:

```bash
# .env
DATABASE_URL="root:tia-dev-password@tcp(127.0.0.1:3306)/tia-dev?charset=utf8mb4&parseTime=True&loc=Local"
JWT_SECRET="dbbf432d3d0205b3fdfb590cd6bd5dc2cb263e584b9cd403d06be3efade76e72"
```

### 2. Running the Application

**Basic Setup (API only)**
```bash
docker-compose up --build
```

**Full Setup (API + CDC Pipeline + Knowledge Graph)**
```bash
./setup-cdc.sh
```

**Local Development**
```bash
./scripts/run.sh
```

### 3. Access Services

- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Memgraph Lab**: http://localhost:7444
- **Kafka Connect**: http://localhost:8083

## Architecture

```
MariaDB → Debezium → Kafka → Knowledge Graph Builder → Memgraph → Connection Analyzer → API Endpoints
```

### Core Components

1. **Main API** - RESTful API with authentication and business logic
2. **Knowledge Graph Builder** - Processes CDC events and builds graph
3. **Connection Analyzer** - Provides intelligent connection recommendations
4. **Memgraph** - Graph database for knowledge graph storage
5. **Kafka** - Message streaming for CDC events
6. **Debezium** - Change Data Capture connector

## API Endpoints

### Connection Recommendations

- `GET /api/v1/connections/complementary/{userId}` - Complementary partners
- `GET /api/v1/connections/alliance/{userId}` - Alliance partners
- `GET /api/v1/connections/mastermind/{userId}` - Mastermind partners
- `GET /api/v1/connections/recommendations/{userId}` - All recommendations
- `GET /api/v1/connections/analysis/{userId}` - Connection analysis

### Core API Endpoints

- `GET /api/v1/users` - User management
- `GET /api/v1/businesses` - Business management
- `GET /api/v1/projects` - Project management
- `GET /api/v1/skills` - Skill management

## Development

### Install Essential Go Tools

```bash
# Visualization Tools
go install github.com/roblaszczak/vgt@latest
go install github.com/goptics/vizb@latest

# Documentation Tools
go install golang.org/x/tools/cmd/godoc@latest
go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

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

## Testing

### Configure Database for Integration Tests

Start the test database:

```bash
docker-compose up -d database
```

### Running Tests

**Unit Tests**
```bash
./scripts/tests/run_unit.sh
```

**Integration Tests**
```bash
./scripts/tests/run_service_integration.sh
./scripts/tests/run_api_integration.sh
```

**Test Visualization**
```bash
./scripts/tests/vis_integration.sh
```

**Code Coverage**
```bash
./scripts/code_coverage.sh
```

**Benchmarking**
```bash
./scripts/benchmark.sh
```

## Documentation

### Generate Swagger Documentation

```bash
swag init -g cmd/api/main.go
```

### View Go Documentation

```bash
godoc -http=:6060
```
Access: `http://localhost:6060/pkg/github.com/TIA-PARTNERS-GROUP/tia-api/`

### Generate Markdown Documentation

```bash
gomarkdoc ./internal/core/services -o docs/services.md
```

## CDC Pipeline & Knowledge Graph

The TIA API includes a Change Data Capture (CDC) pipeline that creates a real-time knowledge graph from your database and provides intelligent connection recommendations.

### Features

- **Real-time CDC**: Captures database changes using Debezium and Kafka
- **Knowledge Graph**: Builds a graph database using Memgraph
- **Smart Connections**: Three types of business connection recommendations

### Connection Types

#### 1. Complementary Partners
Businesses that offer products or services that align well with your own.

**Algorithm**: Business type compatibility analysis
- Technology + Consulting = High compatibility
- Manufacturing + Services = Medium compatibility

#### 2. Alliance Partners
Businesses you can collaborate with on projects.

**Algorithm**: Shared skills + project collaboration potential
- Analyzes skill overlap between users
- Considers active projects for collaboration

#### 3. Mastermind Partners
People who complement your skills through direct collaboration or supportive accountability.

**Algorithm**: Complementary skills + business phase progression
- Startup → Growth/Mature = High compatibility
- Analyzes skill complementarity

### Directory Structure

```
cdc/
├── kg-builder/
│   ├── main.go                   # Knowledge Graph Builder service
│   └── Dockerfile.kg-builder     # Dockerfile for KG Builder
├── connection-analyzer/
│   ├── main.go                   # Connection Analysis service
│   └── Dockerfile.connection-analyzer  # Dockerfile for Connection Analyzer
├── scripts/
│   ├── setup-debezium.sh         # Debezium connector setup
│   └── setup-cdc-pipeline.sh     # CDC pipeline setup
└── README.md                     # CDC documentation
```

### Services

#### Knowledge Graph Builder
- **Port**: Internal service
- **Dependencies**: Kafka, Memgraph
- **Function**: Processes CDC events and builds graph

#### Connection Analyzer
- **Port**: 8082
- **Dependencies**: Memgraph, Main API
- **Function**: Provides connection recommendations

### Testing CDC Pipeline

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

## Troubleshooting

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

### Additional Resources

- [Debezium Documentation](https://debezium.io/documentation/)
- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [Memgraph Documentation](https://memgraph.com/docs/)
- [Cypher Query Language](https://neo4j.com/docs/cypher-manual/)

## Contributing

When contributing to the TIA API:

1. Follow the established directory structure
2. Add comprehensive tests
3. Update documentation
4. Consider performance implications
5. Test with realistic data volumes
6. Follow Go best practices and conventions