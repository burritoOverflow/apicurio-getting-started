package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Represents the structure needed for the v3 API request body
type CreateArtifactRequest struct {
	ArtifactID   string        `json:"artifactId"`
	ArtifactType string        `json:"artifactType,omitempty"` // Type might be inferred by registry in some cases
	Content      SchemaContent `json:"content"`
	// Add other fields like Name, Description if needed
	// Name         string       `json:"name,omitempty"`
	// Description  string       `json:"description,omitempty"`
}

type SchemaContent struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

func main() {
	schemaFilePath := flag.String("schema", "", "Path to the JSON or Protobuf schema file")
	artifactID := flag.String("id", "", "Required: Artifact ID") // Made explicitly required
	groupID := flag.String("group", "default", "Optional: Group ID (defaults to 'default')")
	registryURL := flag.String("registry", "http://localhost:8080/apis/registry/v3", "Apicurio Registry API v3 URL (via proxy)")
	flag.Parse()

	if *schemaFilePath == "" {
		log.Println("Error: -schema flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if *artifactID == "" {
		log.Println("Error: -id flag is required")
		flag.Usage()
		os.Exit(1)
	}

	schemaBytes, err := ioutil.ReadFile(*schemaFilePath)
	if err != nil {
		log.Printf("Error reading schema file '%s': %v\n", *schemaFilePath, err)
		os.Exit(1)
	}
	schemaString := string(schemaBytes) // Content needs to be a string in the JSON body

	var contentType string
	var artifactType string // This might be optional if registry can infer, but good to provide
	fileExt := strings.ToLower(filepath.Ext(*schemaFilePath))
	switch fileExt {
	case ".json":
		contentType = "application/json"
		artifactType = "JSON"
	case ".proto":
		contentType = "application/x-protobuf"
		artifactType = "PROTOBUF"
	default:
		log.Printf("Error: Unsupported schema file extension '%s'. Only .json and .proto are currently supported.\n", fileExt)
		os.Exit(1)
	}

	requestPayload := CreateArtifactRequest{
		ArtifactID:   *artifactID,
		ArtifactType: artifactType,
		Content: SchemaContent{
			Content:     schemaString,
			ContentType: contentType,
		},
	}

	requestBodyBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error marshalling request body to JSON: %v\n", err)
		os.Exit(1)
	}

	// Target the correct v3 endpoint
	apiURL := fmt.Sprintf("%s/groups/%s/artifacts", *registryURL, *groupID)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Set Content-Type for the JSON payload
	req.Header.Set("Content-Type", "application/json")
	log.Printf("Creating artifact '%s' (type: %s) in group '%s' from file '%s' using v3 API...\n", *artifactID, artifactType, *groupID, *schemaFilePath)
	log.Printf("Target URL: %s\n", apiURL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to Apicurio Registry: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		// Continue to show status code if body reading fails
	}

	log.Printf("Response Status Code: %d\n", resp.StatusCode)

	// Handle common v3 status codes (200 OK is typical for successful creation)
	if resp.StatusCode == http.StatusOK {
		log.Printf("Artifact created/updated successfully. Response:\n%s\n", string(responseBodyBytes))
	} else if resp.StatusCode == http.StatusConflict {
		log.Printf("Conflict - Artifact possibly already exists or ID mismatch. Response:\n%s\n", string(responseBodyBytes))
	} else if resp.StatusCode == http.StatusBadRequest {
		log.Printf("Bad Request - Check artifact ID format or request structure. Response:\n%s\n", string(responseBodyBytes))
	} else if resp.StatusCode == http.StatusInternalServerError {
		log.Printf("Internal Server Error. Response:\n%s\n", string(responseBodyBytes))
	} else {
		log.Printf("Received unexpected status code %d. Response:\n%s\n", resp.StatusCode, string(responseBodyBytes))
	}
}
