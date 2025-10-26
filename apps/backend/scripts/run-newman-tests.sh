#!/bin/bash

# Resolve script's directory
SCRIPT_DIR="$(dirname "$0")"
POSTMAN_DIR="$SCRIPT_DIR/../postman"

# Run Newman in Docker
docker run --rm \
  -v "$POSTMAN_DIR:/etc/newman" \
  postman/newman \
  run /etc/newman/FastaBiz-API.postman_collection.json \
  -e /etc/newman/dev_environment.json
