Add artifact:

```bash
./bin/create-artifact -schema ./client/example-schemas/example.json -id my-custom-id3 -group my-group
```

Search artifact

```bash
./bin/search-artifacts -artifactId my-custom-id3 -groupId my-group
```


List all artifacts:

```bash
curl -v localhost:8080/apis/artifacts
```