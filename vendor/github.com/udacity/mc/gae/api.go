package gae

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/udacity/mc/creds"

	uy "github.com/udacity/uy-go"
)

// Response is what we get from GAE.
type Response struct {
	Enrollments         []Enrollment `json:"enrollments,omitempty"`
	ModifiedEnrollments []Enrollment `json:"modified_enrollments,omitempty"`
	User                User         `json:"user,omitempty"`
}

// Enrollment describes an enrollment in GAE.
type Enrollment struct {
	ContentVersion string `json:"content_version,omitempty"`
	NodeKey        string `json:"node_key,omitempty"`
	ProductVariant string `json:"product_variant,omitempty"`
	State          string `json:"state,omitempty"`
}

// Payload is what we send to GAE to change an enrollment.
type Payload struct {
	ContentVersion string   `json:"content_version"`
	Fresh          bool     `json:"fresh"`
	NodeKeys       []string `json:"node_keys"`
	ProductVariant string   `json:"product_variant"`
	State          string   `json:"state"`
}

// User is a user in GAE.
type User struct {
	Email Email `json:"email,omitempty"`
}

// Email includes what we know about a user's email address.
type Email struct {
	Address string `json:"address,omitempty"`
}

var clients []*uy.APIClient

// EnrollInMentorshipND enrolls a uid in nd050. If it returns false, no action was taken, which may mean that the user was already enrolled.
func EnrollInMentorshipND(uid string) (bool, error) {
	payload := Payload{
		"1.0.0", true, []string{"nd050"}, "STANDARD", "enrolled",
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	client, err := getClient()
	if err != nil {
		return false, err
	}

	response := Response{}
	err = client.PostObj(enrollmentsURL(uid), bytes.NewReader(b), &response)
	if err != nil {
		return false, err
	}

	return checkEnrollments(response)
}

func checkEnrollments(r Response) (bool, error) {
	if len(r.ModifiedEnrollments) == 0 {
		return false, nil
	}

	if len(r.ModifiedEnrollments) > 1 {
		return false, fmt.Errorf("error: somehow %d (not 1) enrollments were changed", len(r.ModifiedEnrollments))
	}

	if r.ModifiedEnrollments[0].NodeKey != "nd050" {
		return false, fmt.Errorf("error: wrong node key: %s", r.ModifiedEnrollments[0].NodeKey)
	}

	return r.ModifiedEnrollments[0].State == "enrolled", nil
}

// UnenrollFromMentorshipND unenrolls a uid from nd050. If it returns false, no action was taken, which may mean that the user was not enrolled to begin with.
func UnenrollFromMentorshipND(uid string) (bool, error) {
	payload := Payload{
		"1.0.0", true, []string{"nd050"}, "STANDARD", "unenrolled",
	}

	p, err := json.Marshal(payload)

	client, err := getClient()
	if err != nil {
		return false, err
	}

	response := Response{}
	err = client.PostObj(enrollmentsURL(uid), bytes.NewReader(p), &response)
	if err != nil {
		return false, err
	}

	return checkUnenrollments(response)
}

func checkUnenrollments(r Response) (bool, error) {
	if len(r.ModifiedEnrollments) == 0 {
		return false, nil
	}

	if len(r.ModifiedEnrollments) > 1 {
		return false, fmt.Errorf("error: somehow %d (not 1) enrollments were changed", len(r.ModifiedEnrollments))
	}

	if r.ModifiedEnrollments[0].NodeKey != "nd050" {
		return false, fmt.Errorf("error: wrong node key: %s", r.ModifiedEnrollments[0].NodeKey)
	}

	return r.ModifiedEnrollments[0].State == "unenrolled", nil
}

// FetchEnrollments gets all of a user's ND enrollments, whether currently
// enrolled or not.
func FetchEnrollments(uid string) ([]Enrollment, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	response := Response{}
	err = client.GetObj(enrollmentsURL(uid)+"?fresh=true", &response)
	if err != nil {
		return nil, err
	}

	return response.Enrollments, nil
}

// FetchUser gets info about a user from GAE.
func FetchUser(uid string) (User, error) {
	client, err := getClient()
	if err != nil {
		return User{}, err
	}

	response := Response{}
	err = client.GetObj(userURL(uid), &response)
	if err != nil {
		return User{}, err
	}

	return response.User, nil
}

// creating a client requires an auth trip. reuse clients to avoid re-auth'ing
func getClient() (*uy.APIClient, error) {
	if len(clients) == 0 {
		client := uy.NewAPIClient()
		email, password, err := creds.Load()
		if err != nil {
			return nil, err
		}
		err = client.Auth(email, password)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	return clients[0], nil
}

func enrollmentsURL(uid string) string {
	return fmt.Sprintf("users/%s/enrollments", uid)
}

func userURL(uid string) string {
	return fmt.Sprintf("users/%s?fresh=false", uid)
}
