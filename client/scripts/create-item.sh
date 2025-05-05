#!/usr/bin/env bash

MY_REGISTRY_URL="http://localhost:8080"

# assuming script is run from project root
ITEM_SCHEMA_FILE="./client/example-schemas/Item/Item.json"
ARTIFACT_ID="Item"

# Check if jq is installed
if ! command -v jq &> /dev/null
then
    echo "jq could not be found, please install it."
    exit 1
fi

# Read and compact the JSON content from the file
ITEM_JSON_CONTENT=$(jq -c . "$ITEM_SCHEMA_FILE")

# Check if reading the file was successful
if [ $? -ne 0 ]; then
    echo "Error reading or parsing $ITEM_SCHEMA_FILE"
    exit 1
fi

echo "Creating Item schema in group my-group using content from $ITEM_SCHEMA_FILE"

# Construct the payload with artifactType set to JSON
PAYLOAD=$(printf '{"artifactId":"%s","artifactType":"JSON","firstVersion":{"version":"1.0.0","content":{"content":%s,"contentType":"application/json"}}}' "$ARTIFACT_ID" "$ITEM_JSON_CONTENT")

echo "Sending request to: $MY_REGISTRY_URL/apis/registry/v3/groups/my-group/artifacts"
echo "Payload: $PAYLOAD"

# Send the request using the constructed PAYLOAD
curl -v -X POST "$MY_REGISTRY_URL/apis/registry/v3/groups/my-group/artifacts" \
   -H "Content-Type: application/json" \
   --data "$PAYLOAD"
