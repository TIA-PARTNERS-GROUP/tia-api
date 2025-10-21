#!/bin/bash

# Script to replace all occurrences of {object} gin.H with {object} map[string]interface{}
# in Swagger annotations within Go files.
# WARNING: This script modifies files in place.

# Determine the operating system for correct sed syntax
OS=$(uname)

echo "Searching for and fixing Swagger comments (*.go) recursively..."

# Find all .go files and process them
find . -type f -name "*.go" -print0 | while IFS= read -r -d $'\0' file; do

  # Check for the target string before performing the substitution
  if grep -q "gin.H" "$file"; then
    echo "Updating: $file"

    # The substitution command:
    # s/{object} gin.H/{object} map[string]interface{}/g

    if [[ "$OS" == "Darwin" ]]; then
      # macOS/BSD syntax (requires backup extension)
      sed -i '' 's/{object} gin.H/{object} map[string]interface{}/g' "$file"
    else
      # GNU/Linux syntax (standard -i)
      sed -i 's/{object} gin.H/{object} map[string]interface{}/g' "$file"
    fi

    if [[ $? -ne 0 ]]; then
      echo "ERROR: Failed to process $file"
    fi
  fi
done

echo "--- Swagger Annotation Fix Complete ---"
echo "You must now run 'swag init' in your project root to regenerate documentation."
