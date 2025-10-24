#!/bin/bash

echo "TIA CDC Pipeline Setup"
echo "========================="

if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Please create a .env file with the following variables:"
    echo "DATABASE_URL=\"root:tia-dev-password@tcp(127.0.0.1:3306)/tia-dev?charset=utf8mb4&parseTime=True&loc=Local\""
    echo "JWT_SECRET=\"your-secret-key-here\""
    exit 1
fi

echo "Building Docker images..."
docker-compose build

echo "Starting CDC pipeline services..."
docker-compose up -d

echo "Waiting for services to be ready..."

echo "Waiting for Kafka Connect..."
until curl -f http://localhost:8083/connectors 2>/dev/null; do
    echo "Still waiting for Kafka Connect..."
    sleep 10
done

echo "Kafka Connect is ready!"

echo "Waiting for Memgraph..."
until curl -f http://localhost:7444/status 2>/dev/null; do
    echo "Still waiting for Memgraph..."
    sleep 5
done

echo "Memgraph is ready!"

echo "Setting up Debezium connector..."
./cdc/scripts/setup-debezium.sh

echo ""
echo "CDC Pipeline Setup Complete!"
echo ""
echo "Services Status:"
echo "   - MariaDB: localhost:3306"
echo "   - Kafka: localhost:9092"
echo "   - Kafka Connect: http://localhost:8083"
echo "   - Memgraph: bolt://localhost:7687"
echo "   - Memgraph Lab: http://localhost:7444"
echo "   - Connection Analyzer: http://localhost:8082"
echo ""
echo "Check connector status:"
echo "   curl http://localhost:8083/connectors/tia-mariadb-connector/status"
echo ""
echo "View topics:"
echo "   docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list"
echo ""
echo "Query knowledge graph:"
echo "   Connect to Memgraph Lab at http://localhost:7444"
echo ""
echo "Test connection endpoints:"
echo "   GET http://localhost:8082/api/v1/connections/complementary/{userId}"
echo "   GET http://localhost:8082/api/v1/connections/alliance/{userId}"
echo "   GET http://localhost:8082/api/v1/connections/mastermind/{userId}"
echo "   GET http://localhost:8082/api/v1/connections/recommendations/{userId}"
