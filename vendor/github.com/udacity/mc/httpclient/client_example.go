package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ExampleNew shows how to perform a GET against a (an?) Udacity API.
func ExampleNew() {
	client, err := New("https://reviews-api.udacity.com/api/v1")
	if err != nil {
		// very likely an authentication error
	}

	type Opportunity struct{}

	out := Opportunity{}
	// just use the bare hoth.Backend to make the call
	err = client.Call(http.MethodGet, "opportunities/1", nil, out)
	if err != nil {
		// handle err
	}
}

// ExampleGet shows how to use the Get shorthand method
func ExampleGet() {
	type Project struct{}

	out := Project{}

	err := Get(Backend{}, "https://review-api.udacity.com/api/v1", "project/1", nil, &out)
	if err != nil {
		// handle err
	}

	fmt.Println(out)
}

// ExamplePost shows how to use the Post shorthand method
func ExamplePost() {
	type Project struct{}

	out := Project{}

	params := make(map[string]string)
	params["rubric"] = "changed rubric"
	p, _ := json.Marshal(params)

	err := Post(Backend{}, "https://review-api.udacity.com/api/v1", "project/1", p, &out)
	if err != nil {
		// handle err
	}

	fmt.Println(out)
}

// ExamplePostGraphQL shows how to perform a Hoth auth'ed GraphQL call.
func ExamplePostGraphQL() {
	type Mentor struct{}

	type Data struct {
		Mentor Mentor
	}

	data := Data{}
	// the retriever just needs a New method to get a client. This can either be
	// the default httpclient.Backend or any method that will return a
	// *hoth.Backend.
	retriever := Backend{}
	query := `query (uid: "asdf"){ email }`
	URL := "https://mentor-api.udacity.com/api/v1"

	err := PostGraphQL(retriever, URL, query, &data)
	if err != nil {
		// handle err
	}

	fmt.Println(data.Mentor)
}
