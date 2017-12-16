package httpclient

import (
	"encoding/json"
	"net/http"
)

// Response represents the response object received from a GraphQL API.
type Response struct {
	Data interface{} `json:"data,omitempty"`
}

// PostGraphQL will get the backend client for a GraphQL API and perform the
// POST. out should be a reference. Note that this call will retry if the query
// string is malformed. You'll need to check the logs to see what's wrong.
func PostGraphQL(r Retriever, URL string, query string, out interface{}) error {
	client, err := r.New(URL)
	if err != nil {
		return err
	}

	payload := struct {
		Query string `json:"query"`
	}{
		Query: query,
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	res := Response{}
	res.Data = out
	// hoth.Backend will retry requests on 400 responses, which is what comes back
	// when the query string itself has issues.
	return client.Call(http.MethodPost, "", p, &res)
}

// Get GETs from an endpoint.
func Get(r Retriever, URL string, path string, body []byte, out interface{}) error {
	client, err := r.New(URL)
	if err != nil {
		return err
	}

	return client.Call(http.MethodGet, path, body, out)
}

// Post POSTs to an endpoint.
func Post(r Retriever, URL string, path string, body []byte, out interface{}) error {
	client, err := r.New(URL)
	if err != nil {
		return err
	}

	return client.Call(http.MethodPost, path, body, out)
}
