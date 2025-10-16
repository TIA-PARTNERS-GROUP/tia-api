#!/bin/bash

# Remove all // comments from TypeScript/JavaScript files
find src -name "*.ts" -o -name "*.js" | while read file; do
  echo "Processing $file"

  # Create a temporary file
  temp_file=$(mktemp)

  # Remove // comments but preserve URLs and keep lines with only comments
  # This sed command:
  # 1. Removes // comments at the end of lines
  # 2. Removes lines that are only // comments (but not empty lines)
  # 3. Preserves URLs that contain //
  sed -E 's|([^:]//.*)$||; /^[[:space:]]*\/\/[[:space:]]*$/d; /^[[:space:]]*\/\/[[:space:]]*[^[:space:]].*$/d' "$file" >"$temp_file"

  # Check if the file changed
  if ! cmp -s "$file" "$temp_file"; then
    mv "$temp_file" "$file"
    echo "  -> Comments removed from $file"
  else
    rm "$temp_file"
  fi
done

echo "Comment removal complete!"
