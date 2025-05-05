#!/usr/bin/env bash

curl 'http://localhost:8080/apis/registry/v3/groups/item-group/artifacts/Item'

echo "Getting ItemID"
curl 'http://localhost:8080/apis/registry/v3/groups/item-group/artifacts/ItemId/versions/1.0.0/' | jq
curl 'http://localhost:8080/apis/registry/v3/groups/item-group/artifacts/ItemId/versions/1.0.0/content' | jq

echo "Getting Item"
curl 'http://localhost:8080/apis/registry/v3/groups/item-group/artifacts/Item/versions/1.0.0/' | jq
curl 'http://localhost:8080/apis/registry/v3/groups/item-group/artifacts/Item/versions/1.0.0/content' | jq

GROUP_ID="new-items"
echo "Searching by group ID: ${GROUP_ID}"
curl "http://localhost:8080/apis/registry/v3/search/artifacts?groupId=${GROUP_ID}" | jq

NAMESPACE="com.example.items"
echo "Searching by namespace: ${NAMESPACE}"
curl "http://localhost:8080/apis/registry/v3/search/artifacts?properties=namespace=${NAMESPACE}" | jq

echo "Searching by fields in JSON"
curl "http://localhost:8080/apis/registry/v3/search/artifacts?fields=fields.name" | jq