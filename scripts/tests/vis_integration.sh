#!/bin/bash
echo "--- Visualising Integration Tests ---"
go test -json -v ./test/integration/... | vgt
