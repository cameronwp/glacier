// Package hoth is an API Client for the Udacity Hoth centralized authentication service
package hoth

// https://github.com/udacity/hoth
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	// ProductionURL is the Production URL for the Hoth API
	ProductionURL = "https://hoth.udacity.com/authenticate"
	// StagingURL is the Staging URL for the Hoth API
	StagingURL = "https://hoth-staging.udacity.com/authenticate"
)

// GetJWT sends an email/password to the authenticate endpoint to get a JSON Web Token
func GetJWT(email string, password string) (string, error) {
	type user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	u := user{email, password}
	j, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	br := bytes.NewReader(j)

	req, err := http.NewRequest("POST", ProductionURL, br)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	body := string(resBody)

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP Status:%v Response: %v", res.StatusCode, body)
	}

	// Return the JWT if everything is good.
	return body, nil
}
