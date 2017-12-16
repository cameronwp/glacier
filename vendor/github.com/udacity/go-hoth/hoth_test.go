package hoth_test

import (
	"fmt"
	"testing"

	hoth "github.com/udacity/go-hoth"
)

// ExampleAuth shows how to get a JWT from Hoth.
func ExampleAuth() {
	jwt, err := hoth.GetJWT("user@udacity.com", "example")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(jwt)
}

// ExampleClientInit shows how to init the client and use the backend.
func ExampleClientInit() {
	var users []interface{}

	log := &hoth.DefaultLogger{}
	userAgent := "Udacity Example API Go Bindings"
	var err error
	backend, err := hoth.GetBackend("ProductionURL", "email@udacity.com", "hunter2", userAgent, log)
	if err != nil {
		return
	}

	err = backend.Call("GET", "/users", nil, users)
	if err != nil {
		return
	}
	return
}

// ExampleClientInitWithManualJWT demonstrates how to set a JWT after a backend
// has been instantiated.
func ExampleClientInitWithManualJWT() {
	var users []interface{}

	log := &hoth.DefaultLogger{}
	userAgent := "Udacity Example API Go Bindings"
	var err error
	// note the empty email string
	backend, err := hoth.GetBackend("ProductionURL", "", "", userAgent, log)
	if err != nil {
		return
	}
	backend.SetJWT("jwtpulledfromanothersource")

	err = backend.Call("GET", "/users", nil, users)
	if err != nil {
		return
	}
	return
}

func TestBadAccount(t *testing.T) {
	_, err := hoth.GetJWT("user@udacity.com", "example")
	if err == nil {
		t.Error("Expected bad account to fail")
	}
	t.Logf("Found: %v", err)
}
