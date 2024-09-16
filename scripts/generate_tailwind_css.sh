#!/bin/bash

# Get the directory of the script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Change to the project root directory
cd "$SCRIPT_DIR/.."

# Generate Tailwind CSS
npx tailwindcss -i ./content/assets/input.css -o ./public/main.css --minify