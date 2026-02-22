#!/bin/bash

set -e

TEMPLATE_FILE=".env.example"
OUTPUT_FILE=".env.docker"

if [ ! -f "$TEMPLATE_FILE" ]; then
  echo "$TEMPLATE_FILE not found!"
  exit 1
fi

if [ -f "$OUTPUT_FILE" ]; then
  echo "$OUTPUT_FILE already exists."
  read -p "Overwrite? (y/n): " confirm
  if [ "$confirm" != "y" ]; then
    echo "Aborted."
    exit 0
  fi
fi

cp "$TEMPLATE_FILE" "$OUTPUT_FILE"

echo "Created $OUTPUT_FILE"
echo "Now edit $OUTPUT_FILE and insert your real keys."
