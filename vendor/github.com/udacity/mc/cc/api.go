package cc

import (
	"fmt"

	"github.com/udacity/mc/httpclient"
)

// Data represents a collection of data provided by the Response.
type Data struct {
	User User `json:"user,omitempty"`
}

// User represents a mentor in Classroom Content.
type User struct {
	ID          string       `json:"id,omitempty"`
	Email       string       `json:"email,omitempty"`
	FirstName   string       `json:"first_name,omitempty"`
	LastName    string       `json:"last_name,omitempty"`
	Nanodegrees []Nanodegree `json:"nanodegrees,omitempty"`
}

// Nanodegree represents a degree and its attributes affiliated with a User.
type Nanodegree struct {
	Key             string          `json:"key"`
	Title           string          `json:"title"`
	IsGraduated     bool            `json:"is_graduated"`
	Enrollment      Enrollment      `json:"enrollment"`
	AggregatedState AggregatedState `json:"aggregated_state"`
}

// AggregatedState represents the state of a user's progress in an ND.
type AggregatedState struct {
	CompletionAmount float64 `json:"completion_amount"`
}

// Enrollment represents a users enrollment status in a specific degree
type Enrollment struct {
	Status         string `json:"status"`
	ProductVariant string `json:"product_variant"`
}

// FetchUser fetches a user from prod Classroom-Content API using a UID.
func FetchUser(uid string) (User, error) {
	r := httpclient.Backend{}
	return fetchUser(r, ccAPIURL, uid)
}

func fetchUser(r httpclient.Retriever, baseURL string, uid string) (User, error) {
	query := fmt.Sprintf(userQuery, uid)
	data := Data{}
	err := httpclient.PostGraphQL(r, baseURL, query, &data)

	if err != nil {
		return User{}, err
	}

	return data.User, nil
}

// FetchUserNanodegrees fetches a user and that user's Nanodegrees from prod Classroom-Content API using a UID.
func FetchUserNanodegrees(uid string) (User, error) {
	r := httpclient.Backend{}
	return fetchUserNanodegrees(r, ccAPIURL, uid)
}

func fetchUserNanodegrees(r httpclient.Retriever, baseURL string, uid string) (User, error) {
	query := fmt.Sprintf(userNanodegreesQuery, uid)
	data := Data{}
	err := httpclient.PostGraphQL(r, baseURL, query, &data)

	if err != nil {
		return User{}, err
	}

	return data.User, nil
}

// FetchUserNDProgress gets a user's progress in an ND.
func FetchUserNDProgress(uid string, ndkey string) (User, error) {
	r := httpclient.Backend{}
	return fetchUserNDProgress(r, ccAPIURL, uid, ndkey)
}

func fetchUserNDProgress(r httpclient.Retriever, baseURL string, uid string, ndkey string) (User, error) {
	query := fmt.Sprintf(UserNanodegreeProgressQuery, uid, ndkey)
	data := Data{}
	err := httpclient.PostGraphQL(r, baseURL, query, &data)

	if err != nil {
		return User{}, err
	}

	return data.User, nil
}
