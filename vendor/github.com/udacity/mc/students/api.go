package students

import (
	"fmt"

	"github.com/udacity/mc/httpclient"
)

// Response is what we get from students API.
type Response struct {
	List []User `json:"list,omitempty"`
	More bool   `json:"more,omitempty"`
}

// User is someone in students API.
type User struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	UID       string `json:"key,omitempty"`
	Email     string `json:"email_address,omitempty"`
}

var productionURL = "https://students-api.udacity.com/api/v1"

// Search looks up a user in students API. Limited to 10 responses.
func Search(email string) ([]User, error) {
	r := httpclient.Backend{}
	return search(r, email)
}

func search(r httpclient.Retriever, email string) ([]User, error) {
	path := fmt.Sprintf("accounts/search?q=%s", email)
	response := Response{}
	err := httpclient.Get(r, productionURL, path, nil, &response)
	if err != nil {
		return []User{}, err
	}
	return response.List, nil
}

// FetchAccount gets account info by uid.
func FetchAccount(uid string) (User, error) {
	r := httpclient.Backend{}
	return fetchAccount(r, uid)
}

func fetchAccount(r httpclient.Retriever, uid string) (User, error) {
	path := fmt.Sprintf("accounts/%s", uid)
	user := User{}
	err := httpclient.Get(r, productionURL, path, nil, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
