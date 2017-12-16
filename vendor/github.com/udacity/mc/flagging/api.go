package flagging

import (
	"fmt"

	"github.com/udacity/mc/httpclient"
)

// Flag represents a flag in reviews API.
type Flag struct {
	State string `json:"state,omitempty"`
}

// FetchFlags gets any cheating flags ("suspect" or "cheated") associated with the user.
func FetchFlags(uid string) ([]Flag, error) {
	r := httpclient.Backend{}
	return fetchFlags(r, productionURL, uid)
}

func fetchFlags(r httpclient.Retriever, baseURL string, uid string) ([]Flag, error) {
	var flags []Flag
	err := httpclient.Get(r, baseURL, reportURL(uid), nil, &flags)
	if err != nil {
		return nil, err
	}
	return flags, nil
}

// reportURL creates a URL to get a report of flags on a user.
func reportURL(uid string) string {
	return fmt.Sprintf("report/%s", uid)
}
