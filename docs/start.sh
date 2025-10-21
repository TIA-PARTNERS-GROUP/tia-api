#!/bin/sh

echo "Waiting for API to be ready..."
until curl -f http://api:8080/swagger/doc.json > /app/swagger.json 2>/dev/null; do
  echo "Waiting for API..."
  sleep 5
done

echo "API is ready! Generating documentation..."
swagger-markdown -i /app/swagger.json -o /app/output/API.md
echo "Markdown documentation generated!"
