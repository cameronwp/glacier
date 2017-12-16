package uygo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// WrapError wraps a generic error in the uyclient
type WrapError struct {
	Op   string
	Path string
	Err  error
}

func (e *WrapError) Error() string {
	return e.Op + " " + e.Path + ": " + e.Err.Error()
}

// JSONError records an error with the resulting JSON
type JSONError struct {
	Op   string
	Path string
	JSON string
}

func (e *JSONError) Error() string {
	return e.Op + " " + e.Path + ": " + e.JSON
}

type APIClient struct {
	http       http.Client
	BaseURL    string
	Authorized bool
	XSRFToken  string
	AccountKey string
	SessionID  string
}

var ProdURL = "https://www.udacity.com/api/"
var AbcinthURL = "https://mirror-dot-udacity-abcinth.appspot.com/api/"
var LocalURL = "http://localhost:8080/api/"

// NewAPIClient creates a new Udacity API client at the production URL
func NewAPIClient() *APIClient {
	return CreateAPIClient(ProdURL)
}

// CreateAPIClient creates a new Udacity API client at the given host URL
//
// This can be used to create a client that will work on a dev version of the site, e.g. abcinth
func CreateAPIClient(baseURL string) *APIClient {
	client := &APIClient{BaseURL: baseURL, Authorized: false}
	// the udacity API loves cookies, so...
	client.http.Jar, _ = cookiejar.New(nil)
	return client
}

// NewRequest creates a new `http.Request` (see `http.NewRequest`) with the XSRF-TOKEN and
// Accept headers already set up.
func (c *APIClient) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	if !c.Authorized {
		return nil, fmt.Errorf("Authorize APIClient first")
	}
	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-XSRF-TOKEN", c.XSRFToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Get performs an HTTP `GET` request on path (where path is appended to
// `https://www.udacity.com/api/`)
func (c *APIClient) Get(path string) (*http.Response, error) {

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return c.http.Do(req)
}

// GetObj performs an HTTP `GET` request on `path` (where path is appended
// to baseAPIURL) and unmarshals the repsonse json into `obj`
func (c *APIClient) GetObj(path string, obj interface{}) error {
	resp, err := c.Get(path)
	if err != nil {
		return &WrapError{Path: path, Op: "GetObj", Err: err}
	}
	err = getJSON(resp, obj)
	if err != nil {
		return &WrapError{Path: path, Op: "GetObj", Err: err}
	}
	return nil
}

// Post performs an HTTP `POST` request on `path` with the data in `body`
// (where `path` is appended to `https://www.udacity.com/api/`)
func (c *APIClient) Post(path string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	return c.http.Do(req)
}

// PostObj performs an HTTP `POST` request on `path` with the data in `body`
// and unmarshals the response json into `obj`
func (c *APIClient) PostObj(path string, body io.Reader, obj interface{}) error {
	resp, err := c.Post(path, body)
	if err != nil {
		return err
	}
	err = getJSON(resp, obj)
	return err
}

// Auth authorizes the client using `username` and `password`. You should
// probably call this before anything else.
func (c *APIClient) Auth(username, password string) error {
	type user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	m := make(map[string]user)
	m["udacity"] = user{username, password}
	b, err := json.Marshal(m)
	if err != nil {
		return &WrapError{Op: "Auth-Marshal", Err: err}
	}
	br := bytes.NewReader(b)
	resp, err := c.http.Post(c.BaseURL+"session", "application/json", br)
	if err != nil {
		return &WrapError{Op: "Auth-Post", Err: err}
	}
	defer resp.Body.Close()
	var j map[string]interface{}
	err = getJSON(resp, &j)
	if err != nil {
		return &WrapError{Op: "Auth-getJSON", Err: err}
	}
	// if everything's cool grab the XSRF-TOKEN cookie
	apiURL, err := url.Parse(c.BaseURL + "session")
	if err != nil {
		return &WrapError{Op: "Auth-XSRF-TOKEN", Err: err}
	}
	for _, cookie := range c.http.Jar.Cookies(apiURL) {
		if cookie.Name == "XSRF-TOKEN" {
			c.XSRFToken = cookie.Value
		}
	}
	account := j["account"].(map[string]interface{})
	session := j["session"].(map[string]interface{})
	c.AccountKey = account["key"].(string)
	c.SessionID = session["id"].(string)
	c.Authorized = true
	return nil
}

func getJSON(res *http.Response, obj interface{}) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// Strip first five characters which exist as "security" to make the default
	// JSON response invalid.
	cleanBody := body[5:]
	if res.StatusCode != 200 {
		return &JSONError{Op: "Auth-StatusNot200", JSON: string(cleanBody)}
	}
	return json.Unmarshal(cleanBody, obj)
}
