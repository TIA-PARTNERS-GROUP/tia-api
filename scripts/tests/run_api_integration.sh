#!/bin/bash
echo "--- Running Api Integration Tests ---"
go test ./test/integration/api/... -v
