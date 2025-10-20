#!/bin/bash
set -e

OUTPUT_DIR="code_coverage"
COVERAGE_PROFILE="$OUTPUT_DIR/code_coverage.out"

echo "--- Generating Go Test Coverage Profile ---"

mkdir -p "$OUTPUT_DIR"

go test -coverprofile="$COVERAGE_PROFILE" ./...

echo "  Coverage profile generated: $COVERAGE_PROFILE"
echo "  Opening interactive HTML report in browser..."

go tool cover -html="$COVERAGE_PROFILE"

echo "--- Coverage Visualization Complete ---"
