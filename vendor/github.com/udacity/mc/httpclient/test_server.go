package httpclient

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/udacity/go-hoth"
)

// TestBackend is a Hoth Backend without auth.
type TestBackend struct{}

// New creates an unauthed backend.
func (t TestBackend) New(URL string) (*hoth.Backend, error) {
	return hoth.GetBackend(URL, "", "", "", &hoth.DefaultLogger{})
}

// TestServer represents a server meant to mock an API.
type TestServer struct {
	// Port is the port of the test server. Defaults to 1234. This value is not
	// guaranteed to be used! The final port in use may be different if an address
	// collision occurs. Use the string returned from Open() to get the actual
	// address:port in use.
	Port int
	// Epath is the expected path for requests. If the request path does not match
	// the Epath, the server responds 404. Needs to start with `/`.
	Epath string
	// Dresponse is the dumb response when the Epath matches the request (or if
	// no Epath has been set).
	Dresponse string
	// ReqCB is a callback to run against a request. Great for closures in tests.
	// Does not affect the response.
	ReqCB func(*http.Request)
	// ResCB is a callback tun run against a request to generate a response. Note
	// that you can use r.RequestURI to get the incoming path.
	ResCB func(*http.Request) (string, int)
	// Close the server.
	Close func()
}

// Open runs a live test server. Don't forget to close it! Returns the
// address:port in use. See TestServer struct for server behavior.
func (t *TestServer) Open() string {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond := func(body string, status int) {
			w.WriteHeader(status)
			_, err := w.Write([]byte(body))
			if err != nil {
				log.Fatal(err)
			}
		}

		if t.ReqCB != nil {
			t.ReqCB(r)
		}

		if t.ResCB != nil {
			body, status := t.ResCB(r)
			respond(body, status)
			return
		}

		// default the response to {}
		response := "{}"
		if t.Dresponse != "" {
			response = t.Dresponse
		}

		if t.Epath == "" {
			respond(response, http.StatusOK)
			return
		}

		// r.RequestURI starts with a "/", so add one if necessary
		expectedPath := ""
		if t.Epath != "" {
			if strings.Index(t.Epath, "/") != 0 {
				expectedPath = "/" + t.Epath
			} else {
				expectedPath = t.Epath
			}
		}

		if r.RequestURI == expectedPath {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(response))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte("404 page not found"))
			if err != nil {
				log.Fatal(err)
			}
		}
	}))

	port := 1234

	if t.Port != 0 {
		port = t.Port
	}

	l, listenURL := listen(port)
	ts.Listener = l
	ts.Start()
	t.Close = ts.Close

	return listenURL
}

func listen(port int) (net.Listener, string) {
	localhost := "127.0.0.1"
	listenURL := fmt.Sprintf("%s:%d", localhost, port)
	l, err := net.Listen("tcp", listenURL)
	if err != nil {
		// we sometimes get a race condition with ports in use (as a net.OpError).
		// increment the port and try again
		switch e := err.(type) {
		case *net.OpError:
			port = port + 1
			return listen(port)
		default:
			log.Fatal(e)
		}
	}

	return l, listenURL
}
