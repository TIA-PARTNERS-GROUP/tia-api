#!/bin/bash
echo "--- Running Service Integration Tests ---"
go test ./test/integration/services/... -v
