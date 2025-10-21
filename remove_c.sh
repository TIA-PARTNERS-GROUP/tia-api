#!/bin/bash

# Script to remove all single-line comments (//) from .go files.
# WARNING: This script modifies files in place. Use with caution and back up your code.

# Determine the operating system for correct sed syntax
OS=$(uname)

# Function to process a single Go file
process_file() {
    local file="$1"

    if [[ ! -f "$file" ]]; then
        echo "Error: File not found or not a regular file: $file"
        return 1
    fi

    echo "Processing $file..."

    # The sed commands:
    # 1. /^\s*\/\//d          -> Deletes lines that start only with whitespace followed by //.
    # 2. s/\/\/[^\n\r]*$//   -> Removes '//' and everything that follows until the end of the line.

    if [[ "$OS" == "Darwin" ]]; then
        # macOS/BSD syntax (requires backup extension)
        sed -i '' -E '
            /^\s*\/\//d
            s/\/\/[^\n\r]*$//
        ' "$file"
    else
        # GNU/Linux syntax (does not require backup extension, uses 'i')
        sed -i -E '
            /^\s*\/\//d
            s/\/\/[^\n\r]*$//
        ' "$file"
    fi

    if [[ $? -eq 0 ]]; then
        echo "Successfully cleaned $file."
    else
        echo "Error cleaning $file."
        return 1
    fi
}

# Main script logic: Finds all files ending in .go recursively and processes them
echo "Searching for all .go files recursively..."
echo "Running under OS: $OS"

# Use 'find' to get all .go files recursively
find . -type f -name "*.go" -print0 | while IFS= read -r -d $'\0' file; do
    process_file "$file"
done

echo "--- Comment Removal Complete ---"
