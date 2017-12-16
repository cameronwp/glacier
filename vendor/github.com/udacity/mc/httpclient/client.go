package httpclient

import (
	hoth "github.com/udacity/go-hoth"
	"github.com/udacity/mc/config"
	"github.com/udacity/mc/creds"
	"github.com/udacity/mc/logging"
)

// Retriever represents an interface that can retrieve an auth'ed Hoth backend
// client.
type Retriever interface {
	New(string) (*hoth.Backend, error)
}

// Backend implements Retriever.
type Backend struct{}

// New gets an authenticated backend against an API.
func New(URL string) (*hoth.Backend, error) {
	jwt, err := creds.FetchJWT()
	if err != nil {
		return nil, err
	}

	return getBackend(URL, jwt)
}

// New guarantees an authenticated Hoth client or an error.
func (b Backend) New(URL string) (*hoth.Backend, error) {
	return New(URL)
}

func getBackend(URL string, jwt string) (*hoth.Backend, error) {
	b, err := hoth.GetBackend(URL, "", "", config.UserAgent, logging.FileLogger())
	if err != nil {
		return nil, err
	}

	b.SetJWT(jwt)
	return b, nil
}
