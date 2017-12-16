package httpclient

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/udacity/go-hoth"
)

type testBackend struct{}

func (t testBackend) New(URL string) (*hoth.Backend, error) {
	return hoth.GetBackend(fmt.Sprintf("http://%s", URL), "", "", "", &hoth.DefaultLogger{})
}

type testOut struct {
	Key string `json:"key"`
}

func TestOpen(t *testing.T) {
	testCases := []struct {
		description string
		endpoint    string
		expected    string
		ts          TestServer
	}{
		{
			"should handle requests without expected paths",
			"api/v1/anything",
			"value",
			TestServer{
				Dresponse: `{"key": "value"}`,
			},
		},
		{
			"should handle requests with expected paths",
			"api/v1/all",
			"value",
			TestServer{
				Epath:     "/api/v1/all",
				Dresponse: `{"key": "value"}`,
			},
		},
		{
			"should handle requests with non-slashed expected paths",
			"api/v1/all",
			"value",
			TestServer{
				Epath:     "api/v1/all",
				Dresponse: `{"key": "value"}`,
			},
		},
		{
			"should handle requests with root paths",
			"",
			"value",
			TestServer{
				Dresponse: `{"key": "value"}`,
			},
		},
		{
			"should handle requests with slashed root paths",
			"/",
			"value",
			TestServer{
				Dresponse: `{"key": "value"}`,
			},
		},
		{
			"should handle requests with slashed expected paths against root",
			"",
			"value",
			TestServer{
				Epath:     "/",
				Dresponse: `{"key": "value"}`,
			},
		},
		{
			"should run request callbacks",
			"api/v1/all",
			"value",
			TestServer{
				Dresponse: `{"key": "value"}`,
				ReqCB: func(r *http.Request) {
					if r.RequestURI != "/api/v1/all" {
						panic("this shouldn't happen")
					}
				},
			},
		},
		{
			"should return a custom response based on request",
			"api/v1/1",
			"1",
			TestServer{
				ResCB: func(r *http.Request) (string, int) {
					// pretending like the last int is an ID for a resource
					id := r.RequestURI[len(r.RequestURI)-1 : len(r.RequestURI)]
					res := fmt.Sprintf(`{"key": "%s"}`, id)
					return res, http.StatusOK
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			testServerURL := tc.ts.Open()
			defer tc.ts.Close()

			out := testOut{}
			err := Get(testBackend{}, testServerURL, tc.endpoint, nil, &out)

			if err != nil {
				t.Error(err)
			}

			if out.Key != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, out.Key)
			}
		})
	}
}

func TestError(t *testing.T) {
	ts := TestServer{
		Epath: "/api/v1/resource",
	}
	testServerURL := ts.Open()
	defer ts.Close()

	out := testOut{}
	// note that the endpoint is not the same as the Epath
	err := Get(testBackend{}, testServerURL, "api/v1/differentresource", nil, &out)
	if err == nil {
		t.Errorf("expected an error but found none")
	}
}
