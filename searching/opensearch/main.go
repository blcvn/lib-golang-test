package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

const IndexName = "go-test-index1"

func main() {
	// Initialize the client with SSL/TLS enabled.
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{"https://localhost:9200"},
		Username:  "admin", // For testing only. Don't store credentials in code.
		Password:  "admin",
	})
	if err != nil {
		fmt.Println("cannot initialize", err)
		os.Exit(1)
	}

	// Print OpenSearch version information on console.
	fmt.Println(client.Info())

	// Define index mapping.
	mapping := strings.NewReader(`{
     'settings': {
       'index': {
            'number_of_shards': 4
            }
          }
     }`)

	// Create an index with non-default settings.
	indicesCreateRequest := opensearchapi.IndicesCreateRequest{
		Index: IndexName,
		Body:  mapping,
	}
	fmt.Println("creating index", indicesCreateRequest)

	indicesResponse, err := indicesCreateRequest.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to create index ", err)
		os.Exit(1)
	}
	fmt.Println(indicesResponse)

	// Add a document to the index.
	document := strings.NewReader(`{
        "title": "Moneyball",
        "director": "Bennett Miller",
        "year": "2011"
    }`)

	docId := "1"
	req := opensearchapi.IndexRequest{
		Index:      IndexName,
		DocumentID: docId,
		Body:       document,
	}
	insertResponse, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to insert document ", err)
		os.Exit(1)
	}
	fmt.Println(insertResponse)

	// Search for the document.
	content := strings.NewReader(`{
       "size": 5,
       "query": {
           "multi_match": {
           "query": "miller",
           "fields": ["title^2", "director"]
           }
      }
    }`)

	search := opensearchapi.SearchRequest{
		Body: content,
	}

	searchResponse, err := search.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to search document ", err)
		os.Exit(1)
	}
	fmt.Println(searchResponse)

	// Delete the document.
	delete := opensearchapi.DeleteRequest{
		Index:      IndexName,
		DocumentID: docId,
	}

	deleteResponse, err := delete.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to delete document ", err)
		os.Exit(1)
	}
	fmt.Println("deleting document")
	fmt.Println(deleteResponse)

	// Delete previously created index.
	deleteIndex := opensearchapi.IndicesDeleteRequest{
		Index: []string{IndexName},
	}

	deleteIndexResponse, err := deleteIndex.Do(context.Background(), client)
	if err != nil {
		fmt.Println("failed to delete index ", err)
		os.Exit(1)
	}
	fmt.Println("deleting index", deleteIndexResponse)
}
