#!/bin/bash

echo "TIA API Data Generator"
echo "========================="

if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Please create a .env file with the following variables:"
    echo "DATABASE_URL=\"root:tia-dev-password@tcp(127.0.0.1:3306)/tia-dev?charset=utf8mb4&parseTime=True&loc=Local\""
    echo "JWT_SECRET=\"your-secret-key-here\""
    exit 1
fi

export $(cat .env | grep -v '^#' | xargs)

if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL not set in .env file!"
    exit 1
fi

echo "Generating test data for TIA API..."
echo "Database: $DATABASE_URL"
echo ""

go run scripts/data_generator.go

echo ""
echo "Data generation completed!"
echo "You can now start your API server and test with the generated data."
