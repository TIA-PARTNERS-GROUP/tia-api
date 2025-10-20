#!/bin/bash
echo "--- Running Integration Tests ---"
go test ./test/integration/... -v
