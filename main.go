package main

import (
	"log"
	"net/http"
	"restAPI/apiclient"
	"restAPI/generator"
)

func main() {
	httpClient := apiclient.NewHTTPClient()
	gen := &generator.IncidentGenerator{
		APIClient: &apiclient.Client{
			HTTPClient:  httpClient,
			EndpointURL: "https://task.appranix.net/issues.json",
		},
		FieldsData: generator.FieldData{
			ID:   "1",
			Name: "Component",
			Value: []string{
				"Others",
			},
		},
	}

	http.HandleFunc("/incident", gen.CreateIncident)
	log.Fatal(http.ListenAndServe(":1234", nil))
}
