#!/bin/bash
set -e

OUTPUT_DIR="benchmark"
OUTPUT_FILE="$OUTPUT_DIR/benchmark_results.txt"
HTML_REPORT="$OUTPUT_DIR/benchmark_report.html"

echo "--- Running Go Benchmarks ---"

mkdir -p "$OUTPUT_DIR"

echo "  Saving raw results to: $OUTPUT_FILE"

go test -bench . -run=^$ -benchmem ./... > "$OUTPUT_FILE"

echo "  Raw benchmark run complete. Generating HTML chart via 'vizb'."
echo "  HTML report saved to: $HTML_REPORT"

vizb "$OUTPUT_FILE" -o "$HTML_REPORT"

if command -v open >/dev/null 2>&1; then
    open "$HTML_REPORT" # macOS
elif command -v xdg-open >/dev/null 2>&1; then
    xdg-open "$HTML_REPORT" # Linux
else
    echo "Could not automatically open $HTML_REPORT. Please open it manually in your browser."
fi

echo "--- Benchmark Visualization Complete ---"
