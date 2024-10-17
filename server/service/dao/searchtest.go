package dao

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"

// 	"github.com/elastic/go-elasticsearch/v7"
// 	"github.com/elastic/go-elasticsearch/v7/esapi"
// )

// // Search performs a search request with the given parameters
// func SearchTest(es *elasticsearch.Client, index string, queryBody io.Reader) (string, error) {
// 	// Prepare the search request
// 	req := esapi.SearchRequest{
// 		Index:  []string{index},
// 		Body:   queryBody, // Use the io.Reader directly
// 		Pretty: true,      // pretty print the JSON response
// 		TrackTotalHits: &esapi.TrackTotalHits{
// 			Relation: "eq",
// 			Value:    true,
// 		},
// 	}

// 	// Perform the search request
// 	res, err := req.Do(context.Background(), es)
// 	if err != nil {
// 		log.Fatalf("Error getting response: %s", err)
// 		return "", err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		var e map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
// 			log.Fatalf("Error parsing the response body: %s", err)
// 			return "", err
// 		}
// 		// Print the response status and error information.
// 		log.Fatalf("[%s] %s: %s",
// 			res.Status(),
// 			e["error"].(map[string]interface{})["type"],
// 			e["error"].(map[string]interface{})["reason"],
// 		)
// 		return "", fmt.Errorf("Elasticsearch error: %s", e["error"].(map[string]interface{})["reason"])
// 	}

// 	// Read and return the entire response body
// 	var buf []byte
// 	if _, err := io.ReadAll(res.Body, &buf); err != nil {
// 		log.Fatalf("Error reading the response body: %s", err)
// 		return "", err
// 	}

// 	return string(buf), nil
// }
// ;
