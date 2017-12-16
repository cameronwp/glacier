package hoth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// Nearly infinite. Wait forever verses bad data, we choose: forever.
	httpRetries           = 100
	maxRetrySleep float64 = 5 * float64(time.Minute)
)

var (
	// ErrRetryLimit is returned when a HTTP request has been tried httpRetries times.
	ErrRetryLimit = errors.New("hoth.backend: HTTP request retry limit reached")
)

var httpClient = &http.Client{
	Timeout: time.Second * 30,
}

// HTTPError is returned for HTTP >=400 statuses. Many APIs use 404 when there are no results, so we don't
// want a fatal error.
type HTTPError struct {
	Msg        string
	StatusCode int
	URL        *url.URL
}

// Error returns the msg as a string.
func (e *HTTPError) Error() string { return e.Msg }

// GetBackend creates a Backend struct and gets the required JWT if email has
// been specified. Leave email blank if you want to set the JWT manually.
func GetBackend(apiURL string, email string, password string, userAgent string, log Logger) (*Backend, error) {
	var jwt string
	if email != "" {
		var err error
		log.Debugf("getting JWT for: %v", email)
		// Get the JSON Web Token required to access the Udacity APIs
		jwt, err = GetJWT(email, password)
		if err != nil {
			return nil, err
		}
	}

	return &Backend{apiURL: apiURL, jwt: jwt, userAgent: userAgent, Log: log}, nil
}

// SetJWT sets a backend's JWT.
func (b *Backend) SetJWT(jwt string) {
	b.jwt = jwt
}

// Backend performs HTTP requests against Udacity APIs needing a Hoth JWT.
type Backend struct {
	apiURL     string
	jwt        string
	userAgent  string
	Log        Logger
	LogVerbose bool
}

// Call creates and executes the request. Called by the endpoint clients.
// Fills the provided interface with the JSON result of the call.
func (b Backend) Call(method, path string, body []byte, v interface{}) error {

	for i := 0; i < httpRetries; i++ {
		// Build the request.
		bodyReader := bytes.NewBuffer(body)
		req, err := b.NewRequest(method, path, bodyReader)
		if err != nil {
			return err
		}
		// Close connection after request. Default cached connections will get
		// failures in the event of server closing idle connections.
		// TODO: Fixed in Go 1.7: https://github.com/golang/go/issues/4677
		req.Close = true

		// Run the request.
		err = b.Do(req, v)
		if err == nil {
			// Return quickly if everything worked.
			return nil
		}

		// Handle all of the errors, anything not returned will retry.
		switch e := err.(type) {
		case *HTTPError:
			if e.StatusCode == 404 {
				// 404 isn't a problem, don't log or retry, let caller handle.
				return err
			}
			// Log everything else >= 400.
			b.Log.Warningf("HTTP status %v url:%v body:%v", e.StatusCode, e.URL, e.Msg)
			// TODO: Stop retrying on 4?? client errors, when the APIs stop sending them incorrectly.
		case *net.DNSError:
			b.Log.Warningf("DNSError: %v", e)
		case *url.Error:
			if !e.Temporary() {
				b.Log.Errorf("permanent: %v", e)
				return e
			}
			// Retry temporary URL errors
			b.Log.Warningf("temporary: %v", e)
		default:
			// Log other errors and return.
			b.Log.Errorf("unhandled error %T: %v url:%v", err, err, req.URL)
			return err
		}
		// Sleep a bit between retry loop.
		retrySleep := getRetrySleep(i)

		b.Log.Infof("sleep - HTTP request retry in %v (%d/%d)", retrySleep, i+1, httpRetries)
		time.Sleep(retrySleep)
	}
	// Loop has run httpRetries times.
	return ErrRetryLimit
}

// NewRequest is used by Call to generate an http.Request. It handles encoding
// parameters and attaching the appropriate headers.
func (b Backend) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := b.apiURL + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", b.userAgent)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", b.jwt))

	return req, nil
}

// Do is used by Call to execute an API request and parse the response. It uses
// the backend's HTTP client to execute the request and unmarshals the response
// into v. It also handles unmarshaling errors returned by the API.
func (b Backend) Do(req *http.Request, v interface{}) error {
	b.Log.Debugf("%v request %v", req.Method, req.URL)

	start := time.Now()

	res, err := httpClient.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()
	if b.LogVerbose {
		b.Log.Debugf("HTTP status %v completed in %v", res.StatusCode, time.Since(start))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if b.LogVerbose {
		b.Log.Debugf(string(resBody))
	}

	if res.StatusCode >= 400 {
		return b.responseToError(res, resBody)
	}

	if v != nil {
		return json.Unmarshal(resBody, v)
	}

	return nil
}

func (b Backend) responseToError(res *http.Response, resBody []byte) error {
	return &HTTPError{
		Msg:        string(resBody),
		StatusCode: res.StatusCode,
		URL:        res.Request.URL,
	}
}

// getRetrySleep returns a time to sleep using an exponential backoff algorithm
func getRetrySleep(retryCount int) time.Duration {
	// TODO: Jitter?
	mili := float64(time.Millisecond)
	retry := math.Pow(2, float64(retryCount)) * mili * 100
	return time.Duration(math.Min(retry, maxRetrySleep))
}
