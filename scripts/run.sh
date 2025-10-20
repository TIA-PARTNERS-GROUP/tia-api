#!/bin/bash
echo "--- Running Main ---"
swag init -g cmd/api/main.go
go run cmd/api/main.go
