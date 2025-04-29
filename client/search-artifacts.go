package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ArtifactSearchResult struct {
	CreatedOn    time.Time `json:"createdOn"` // Use time.Time for dates
	Owner        string    `json:"owner"`
	ArtifactType string    `json:"artifactType"`
	ModifiedOn   time.Time `json:"modifiedOn"` // Use time.Time for dates
	ModifiedBy   string    `json:"modifiedBy"`
	GroupID      string    `json:"groupId"`
	ArtifactID   string    `json:"artifactId"`
	// Add other fields if needed (e.g., name, description, labels)
	Name        *string           `json:"name,omitempty"`        // Use pointer for optional fields
	Description *string           `json:"description,omitempty"` // Use pointer for optional fields
	Labels      map[string]string `json:"labels,omitempty"`
}

type SearchResponse struct {
	Artifacts []ArtifactSearchResult `json:"artifacts"`
	Count     int                    `json:"count"`
}

// Helper type for handling multiple label flags
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	var (
		name         = flag.String("name", "", "Filter by artifact name")
		offset       = flag.Int("offset", 0, "The number of artifacts to skip")
		limit        = flag.Int("limit", 20, "The number of artifacts to return")
		order        = flag.String("order", "", "Sort order: asc or desc")
		orderby      = flag.String("orderby", "", "Sort by: name, createdOn, modifiedOn, groupId, artifactId, artifactType")
		labels       stringSlice // Custom type for multiple label flags
		description  = flag.String("description", "", "Filter by description")
		groupID      = flag.String("groupId", "", "Filter by artifact group")
		globalID     = flag.Int64("globalId", 0, "Filter by globalId")
		contentID    = flag.Int64("contentId", 0, "Filter by contentId")
		artifactID   = flag.String("artifactId", "", "Filter by artifactId")
		artifactType = flag.String("artifactType", "", "Filter by artifact type (e.g., JSON, PROTOBUF)")
		registryURL  = flag.String("registry", "http://localhost:8080/apis/registry/v3", "Apicurio Registry API v3 URL (via proxy)")
	)
	flag.Var(&labels, "label", "Filter by label (key:value format). Can be used multiple times.")
	flag.Parse()

	baseURL := fmt.Sprintf("%s/search/artifacts", *registryURL)
	params := url.Values{}

	if *name != "" {
		params.Add("name", *name)
	}
	if *offset != 0 { // Only add if not default
		params.Add("offset", fmt.Sprintf("%d", *offset))
	}
	if *limit != 20 { // Only add if not default
		params.Add("limit", fmt.Sprintf("%d", *limit))
	}
	if *order != "" {
		params.Add("order", *order)
	}
	if *orderby != "" {
		params.Add("orderby", *orderby)
	}
	if len(labels) > 0 {
		for _, label := range labels {
			params.Add("labels", label) // Add each label separately
		}
	}
	if *description != "" {
		params.Add("description", *description)
	}
	if *groupID != "" {
		params.Add("groupId", *groupID)
	}
	if *globalID != 0 {
		params.Add("globalId", fmt.Sprintf("%d", *globalID))
	}
	if *contentID != 0 {
		params.Add("contentId", fmt.Sprintf("%d", *contentID))
	}
	if *artifactID != "" {
		params.Add("artifactId", *artifactID)
	}
	if *artifactType != "" {
		params.Add("artifactType", *artifactType)
	}

	searchURL := baseURL
	if len(params) > 0 {
		searchURL += "?" + params.Encode()
	}

	log.Printf("Searching artifacts using URL: %s\n", searchURL)

	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error searching artifacts: Received status code %d. Response: %s", resp.StatusCode, string(bodyBytes))
	}

	var searchResult SearchResponse
	err = json.Unmarshal(bodyBytes, &searchResult)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON response: %v\nResponse was: %s", err, string(bodyBytes))
	}

	log.Printf("Successfully searched artifacts. Found %d matching artifact(s).\n", searchResult.Count)

	if searchResult.Count > 0 {
		fmt.Println("--------------------------------------------------")
		for i, artifact := range searchResult.Artifacts {
			fmt.Printf("Artifact %d:\n", i+1)
			fmt.Printf("  Group ID:     %s\n", artifact.GroupID)
			fmt.Printf("  Artifact ID:  %s\n", artifact.ArtifactID)
			fmt.Printf("  Type:         %s\n", artifact.ArtifactType)
			if artifact.Name != nil {
				fmt.Printf("  Name:         %s\n", *artifact.Name)
			}
			if artifact.Description != nil {
				fmt.Printf("  Description:  %s\n", *artifact.Description)
			}
			fmt.Printf("  Created On:   %s\n", artifact.CreatedOn.Format(time.RFC3339))
			fmt.Printf("  Modified On:  %s\n", artifact.ModifiedOn.Format(time.RFC3339))
			if len(artifact.Labels) > 0 {
				fmt.Println("  Labels:")
				for k, v := range artifact.Labels {
					fmt.Printf("    %s: %s\n", k, v)
				}
			}
			fmt.Println("--------------------------------------------------")
		}
	}
}
