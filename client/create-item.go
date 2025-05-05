package main

import (
	"bytes"
	"encoding/json"
	"flag" // Import flag package
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath" // Import filepath package
	"strings"       // Import strings package
)

const (
	registryURL  = "http://localhost:8080"
	artifactType = "JSON"
	version      = "1.0.0"
	contentType  = "application/json"
	draftState   = "DRAFT"
)

// Structs remain the same
type Content struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

type VersionInfo struct {
	Version string  `json:"version"`
	State   string  `json:"state,omitempty"`
	Content Content `json:"content"`
}

type CreateArtifactPayload struct {
	ArtifactID   string      `json:"artifactId"`
	ArtifactType string      `json:"artifactType"`
	FirstVersion VersionInfo `json:"firstVersion"`
}

func main() {
	schemaFilePath := flag.String("schema", "", "Path to the schema file (e.g., ./client/example-schemas/Item/Item.json)")
	groupID := flag.String("group", "", "Group ID for the artifact")
	flag.Parse()

	if *schemaFilePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Schema file path must be provided using the -schema flag.")
		flag.Usage()
		os.Exit(1)
	}

	if *groupID == "" {
		fmt.Fprintln(os.Stderr, "Error: Group ID must be provided using the -group flag.")
		flag.Usage()
		os.Exit(1)
	}

	baseFilename := filepath.Base(*schemaFilePath)
	artifactID := strings.TrimSuffix(baseFilename, filepath.Ext(baseFilename))
	if artifactID == "" {
		fmt.Fprintf(os.Stderr, "Error: Could not derive artifact ID from filename: %s\n", baseFilename)
		os.Exit(1)
	}
	fmt.Printf("Using Schema File: %s\n", *schemaFilePath)
	fmt.Printf("Derived Artifact ID: %s\n", artifactID)

	schemaContentBytes, err := os.ReadFile(*schemaFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading schema file %s: %v\n", *schemaFilePath, err)
		os.Exit(1)
	}
	schemaContentString := string(schemaContentBytes)

	payload := CreateArtifactPayload{
		ArtifactID:   artifactID, // Use derived artifactID
		ArtifactType: artifactType,
		FirstVersion: VersionInfo{
			Version: version,
			State:   draftState, // Keep state if needed
			Content: Content{
				Content:     schemaContentString,
				ContentType: contentType,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling payload: %v\n", err)
		os.Exit(1)
	}

	apiURL := fmt.Sprintf("%s/apis/registry/v3/groups/%s/artifacts", registryURL, *groupID)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("Sending request to: %s\n", apiURL)
	fmt.Printf("Sending payload: %s\n", string(payloadBytes))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Printf("Response Status: %s\n", resp.Status)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, responseBody, "", "  "); err == nil {
		fmt.Printf("Response Body:\n%s\n", prettyJSON.String())
	} else {
		fmt.Printf("Response Body: %s\n", string(responseBody))
	}

	if resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "API request failed with status code %d\n", resp.StatusCode)
		os.Exit(1)
	}
}
