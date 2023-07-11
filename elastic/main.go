package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	// username := "elastic"
	// password := "test@123432324"
	// username := "elastic"
	// password := "WkMqCuDgiaPXl4ZZAX=2"
	// password := "BwRCbJgoj9zUUj=1dpWH"

	// certPath := "es_ca.crt"
	// cert, err := ioutil.ReadFile(certPath)
	// if err != nil {
	// 	log.Fatalf("ioutil.ReadFile failed: %s", err.Error())
	// 	return
	// }
	// transport := http.Transport{
	// 	MaxIdleConnsPerHost:   10,
	// 	ResponseHeaderTimeout: time.Second,
	// 	TLSClientConfig: &tls.Config{
	// 		MinVersion:         tls.VersionTLS12,
	// 		InsecureSkipVerify: true,
	// 	},
	// }

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://103.161.38.247:19200",
			// "https://localhost:9200",
		},
		// Transport: &transport,
		// Username:  username,
		// Password:  password,
		// CACert:    cert,
	}
	es, _ := elasticsearch.NewClient(cfg)
	log.Printf("Elastic Version: %s", elasticsearch.Version)

	res, err := es.Info()
	if err != nil {
		log.Fatalf("es.Info failed: %s", err.Error())
		return
	}
	defer res.Body.Close()
	log.Printf("ES Info %+v", res)

	var wg sync.WaitGroup

	for i, title := range []string{"Test One", "Test Two"} {
		wg.Add(1)

		go func(i int, title string) {
			defer wg.Done()

			// Build the request body.
			data, err := json.Marshal(struct{ Title string }{Title: title})
			if err != nil {
				log.Fatalf("Error marshaling document: %s", err)
				return
			}

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "test",
				DocumentID: strconv.Itoa(i + 1),
				Body:       bytes.NewReader(data),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
				return
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
				return
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(i, title)
	}
	wg.Wait()

	log.Printf("Test: %s", strings.Repeat("-", 37))

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
		return
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Println(r)
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))
}
