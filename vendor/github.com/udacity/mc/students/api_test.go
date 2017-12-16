package students

import "testing"
import "github.com/udacity/mc/httpclient"
import "fmt"

var searchRes = `{"more":false,"list":[{"key":"123456789","first_name":"ferris","last_name":"bueller","nickname":"noschool","email_address":"ferris@udacity.com","email_code_sent":true,"email_date_code_sent":"2017-11-22T23:24:03.999Z","verified":true,"city":null,"state":null,"country":null,"tags":[{"tag":"upgraded","namespace":"classroom-2.0"},{"tag":"visited","namespace":"classroom-2.0"}],"email_preferences":{"ok_user_research":true,"master_ok":true,"ok_course":true},"languages":"en-US","preferred_language":"en-US","created_at":"2017-02-13T21:30:45.324Z","updated_at":"2017-11-22T23:24:03.999Z"}]}`

func TestSearch(t *testing.T) {
	t.Run("Good response", func(t *testing.T) {
		ts := httpclient.TestServer{
			Epath:     "accounts/search",
			Dresponse: searchRes,
		}
		testURL := ts.Open()
		defer ts.Close()

		out := Response{}
		err := httpclient.Get(httpclient.TestBackend{}, fmt.Sprintf("http://%s", testURL), "accounts/search", nil, &out)
		if err != nil {
			t.Error(err)
		}

		if len(out.List) != 1 {
			t.Errorf("expected 1 user, found %d", len(out.List))
		}

		if out.List[0].FirstName != "ferris" {
			t.Errorf("expected user to be 'ferris', found %s", out.List[0].FirstName)
		}
	})

	t.Run("Bad response", func(t *testing.T) {
		ts := httpclient.TestServer{
			Epath: "accounts/search",
		}
		testURL := ts.Open()
		defer ts.Close()

		out := Response{}
		err := httpclient.Get(httpclient.TestBackend{}, fmt.Sprintf("http://%s", testURL), "accounts/search", nil, &out)
		if err != nil {
			t.Error(err)
		}

		if len(out.List) != 0 {
			t.Errorf("expected 1 user, found %d", len(out.List))
		}
	})
}
