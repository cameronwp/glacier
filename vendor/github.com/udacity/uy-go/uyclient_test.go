package uygo_test

import (
	"fmt"
	uy "github.com/udacity/uy-go"
	"io/ioutil"
	"os"
)

func ExampleAuth() {
	client := uy.NewAPIClient()
	err := client.Auth(os.Getenv("UY_USERNAME"), os.Getenv("UY_PASSWORD"))
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("Success!")
	// Output: Success!
}

func ExampleGetMe() {
	client := uy.NewAPIClient()
	err := client.Auth(os.Getenv("UY_USERNAME"), os.Getenv("UY_PASSWORD"))
	if err != nil {
		fmt.Printf("Auth Error: %v", err)
	}
	resp, err := client.Get("users/me")
	if err != nil {
		fmt.Printf("Get me Error: %v", err)
	}
	defer resp.Body.Close()
	jb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read response body Error: %v", err)
	}
	// check the response status
	if resp.StatusCode != 200 {
		fmt.Printf("Error: %d : %s {%s}", resp.StatusCode, resp.Status, string(jb))
	}
	// we're all good. the 'me' payload is in `jb`, you might want to parse it
	// with `encoding/json`
	fmt.Printf("Success!")
	// Output: Success!
}

// ExampleGetMeObj demonstrates how to get json payloads back as structs from
// the api client.
func ExampleGetMeObj() {
	client := uy.NewAPIClient()
	err := client.Auth(os.Getenv("UY_USERNAME"), os.Getenv("UY_PASSWORD"))
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	// you can automatically unmarshal into structs...
	type email struct {
		VerificationCode string `json:"_verification_code"`
		Address          string
	}
	type user struct {
		Nickname  string
		Email     email
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	type payload struct {
		User user
	}
	var p payload
	err = client.GetObj("users/me", &p)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("Success!")
	// Output: Success!
}
