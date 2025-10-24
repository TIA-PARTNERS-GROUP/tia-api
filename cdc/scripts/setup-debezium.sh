#!/bin/bash

# Debezium Connector Setup Script
# This script configures the Debezium connector to capture changes from MariaDB

echo "Setting up Debezium Connector for MariaDB CDC..."

# Wait for Kafka Connect to be ready
echo "Waiting for Kafka Connect to be ready..."
until curl -f http://localhost:8083/connectors; do
  echo "Waiting for Kafka Connect..."
  sleep 5
done

echo "Kafka Connect is ready!"

# Create Debezium connector configuration
CONNECTOR_CONFIG='{
  "name": "tia-mariadb-connector",
  "config": {
    "connector.class": "io.debezium.connector.mysql.MySqlConnector",
    "tasks.max": "1",
    "database.hostname": "database",
    "database.port": "3306",
    "database.user": "root",
    "database.password": "tia-dev-password",
    "database.server.id": "184054",
    "database.server.name": "tia-db",
    "database.include.list": "tia-dev",
    "database.history.kafka.bootstrap.servers": "kafka:29092",
    "database.history.kafka.topic": "dbhistory.tia",
    "include.schema.changes": "true",
    "transforms": "route",
    "transforms.route.type": "org.apache.kafka.connect.transforms.RegexRouter",
    "transforms.route.regex": "([^.]+)\\.([^.]+)\\.([^.]+)",
    "transforms.route.replacement": "$3",
    "key.converter": "org.apache.kafka.connect.json.JsonConverter",
    "value.converter": "org.apache.kafka.connect.json.JsonConverter",
    "key.converter.schemas.enable": "false",
    "value.converter.schemas.enable": "false",
    "snapshot.mode": "initial",
    "binlog.buffer.size": "8192",
    "max.batch.size": "2048",
    "max.queue.size": "8192",
    "poll.interval.ms": "1000"
  }
}'

echo "Creating Debezium connector..."
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d "$CONNECTOR_CONFIG"

echo ""
echo "Debezium connector created!"
echo "You can check connector status at: http://localhost:8083/connectors/tia-mariadb-connector/status"
echo "Topics will be created as: tia-db.tia-dev.{table_name}"
