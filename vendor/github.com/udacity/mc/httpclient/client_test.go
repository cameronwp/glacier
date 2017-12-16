package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	hoth "github.com/udacity/go-hoth"
	"github.com/udacity/mc/config"
)

type testRetriever struct{}

func (t testRetriever) New(URL string) (*hoth.Backend, error) {
	return hoth.GetBackend(URL, "", "", config.UserAgent, &hoth.DefaultLogger{})
}

func TestGetBackend(t *testing.T) {
	t.Run("does not retrieve a JWT", func(t *testing.T) {
		_, err := getBackend("http://youdacity.com", "asdfasdf")
		if err != nil {
			t.Error(err)
		}
	})
}

func TestPostGraphQL(t *testing.T) {
	query := `query user(uid: "asdf"){ key }`
	res := `{"data": {"user": {"key": "value"}}}`
	ts := TestServer{
		Dresponse: res,
	}
	testServerURL := ts.Open()
	defer ts.Close()

	type User struct {
		Key string `json:"key"`
	}
	type Data struct {
		User User `json:"user"`
	}
	data := Data{}

	err := PostGraphQL(testRetriever{}, fmt.Sprintf("http://%s", testServerURL), query, &data)
	if err != nil {
		t.Error(err)
	}

	if data.User.Key != "value" {
		t.Errorf("Expected 'value' but found %s", data.User.Key)
	}
}

func TestGet(t *testing.T) {
	t.Run("should run a GET", func(t *testing.T) {
		ts := TestServer{
			ReqCB: func(r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected a GET, found %s", r.Method)
				}
			},
		}
		testServerURL := ts.Open()
		defer ts.Close()

		out := testOut{}
		err := Get(testBackend{}, testServerURL, "", nil, &out)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("should marshal responses into an interface", func(t *testing.T) {
		expected := "data"
		ts := TestServer{
			Dresponse: fmt.Sprintf(`{"key": "%s"}`, expected),
		}
		testServerURL := ts.Open()
		defer ts.Close()

		out := testOut{}
		err := Get(testBackend{}, testServerURL, "", nil, &out)
		if err != nil {
			t.Error(err)
		}
		if out.Key != expected {
			t.Errorf("expected %s response, found %s", expected, out.Key)
		}
	})

	t.Run("should pass a body to request", func(t *testing.T) {
		body := []byte(`{"param": "value"}`)
		ts := TestServer{
			ReqCB: func(r *http.Request) {
				fbody, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
				}
				if !bytes.Equal(fbody, body) {
					t.Errorf("expected %s response, found %s", body, fbody)
				}
			},
		}
		testServerURL := ts.Open()
		defer ts.Close()

		out := testOut{}
		err := Get(testBackend{}, testServerURL, "", body, &out)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestPost(t *testing.T) {
	t.Run("should run a POST", func(t *testing.T) {
		ts := TestServer{
			ReqCB: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected a POST, found %s", r.Method)
				}
			},
		}
		testServerURL := ts.Open()
		defer ts.Close()

		out := testOut{}
		err := Post(testBackend{}, testServerURL, "", nil, &out)
		if err != nil {
			t.Error(err)
		}
	})
}
