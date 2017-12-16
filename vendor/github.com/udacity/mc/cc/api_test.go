package cc

import (
	"fmt"
	"testing"

	"github.com/udacity/mc/httpclient"
)

const (
	simpleRes   = `{"data": {"user": {"id": "1234"}}}`
	noNDsRes    = `{"data": {"user": {"id": "1234", "nanodegrees": []}}}`
	oneNDRes    = `{"data": {"user": {"id": "1234", "nanodegrees": [{"key": "1"}]}}}`
	progressRes = `{"data": {"user": {"id": "1234", "nanodegrees": [{"key": "1", "aggregated_state": {"completion_amount": 0.5}}]}}}`
)

func TestFetchUser(t *testing.T) {
	ts := httpclient.TestServer{
		Dresponse: simpleRes,
	}
	testServerURL := ts.Open()
	defer ts.Close()

	r := httpclient.TestBackend{}
	user, err := fetchUser(r, fmt.Sprintf("http://%s", testServerURL), "1234")
	if err != nil {
		t.Errorf("Error in response: %s", err)
	}
	if user.ID != "1234" {
		t.Errorf("Expected id of 1234, found %s", user.ID)
	}
	if user.Email != "" {
		t.Errorf("Expected blank email, found %s", user.Email)
	}
}

func TestUserNoNDs(t *testing.T) {
	ts := httpclient.TestServer{
		Dresponse: noNDsRes,
	}
	testServerURL := ts.Open()
	defer ts.Close()

	r := httpclient.TestBackend{}
	user, err := fetchUserNanodegrees(r, fmt.Sprintf("http://%s", testServerURL), "1234")
	if err != nil {
		t.Errorf("Error in response: %s", err)
	}
	if len(user.Nanodegrees) != 0 {
		t.Errorf("Expected 0 NDs, found %d", len(user.Nanodegrees))
	}
}

func TestUserOneND(t *testing.T) {
	ts := httpclient.TestServer{
		Dresponse: oneNDRes,
	}
	testServerURL := ts.Open()
	defer ts.Close()

	r := httpclient.TestBackend{}
	user, err := fetchUserNanodegrees(r, fmt.Sprintf("http://%s", testServerURL), "1234")
	if err != nil {
		t.Errorf("Error in response: %s", err)
	}
	if len(user.Nanodegrees) != 1 {
		t.Errorf("Expected 1 ND, found %d", len(user.Nanodegrees))
	}
	if user.Nanodegrees[0].Key != "1" {
		t.Errorf("Expected key of 1, found %s", user.Nanodegrees[0].Key)
	}
	if user.Nanodegrees[0].AggregatedState.CompletionAmount != 0 {
		t.Errorf("Expected 0 completion amount, found %f", user.Nanodegrees[0].AggregatedState.CompletionAmount)
	}
}

func TestUserProgress(t *testing.T) {
	ts := httpclient.TestServer{
		Dresponse: progressRes,
	}
	testServerURL := ts.Open()
	defer ts.Close()

	r := httpclient.TestBackend{}
	user, err := fetchUserNanodegrees(r, fmt.Sprintf("http://%s", testServerURL), "1234")
	if err != nil {
		t.Errorf("Error in response: %s", err)
	}
	if user.Nanodegrees[0].AggregatedState.CompletionAmount != 0.5 {
		t.Errorf("Expected 0.5 completion amount, found %f", user.Nanodegrees[0].AggregatedState.CompletionAmount)
	}
}
