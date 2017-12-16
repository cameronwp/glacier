package flagging

import (
	"fmt"
	"testing"

	"github.com/udacity/mc/httpclient"
)

const (
	errorFlagsRes   = `badresponse`
	nilFlagsRes     = `[]`
	suspectFlagsRes = `[{"state": "suspect"}]`
	cheatedFlagsRes = `[{"state": "cheated"}]`
	bothFlagsRes    = `[{"state": "cheated"}, {"state": "suspect"}]`
)

func TestGetFlags(t *testing.T) {
	testCases := []struct {
		name      string
		res       string
		checker   func(fs []Flag) (bool, string)
		shouldErr bool
	}{
		{"nil response",
			nilFlagsRes,
			func(fs []Flag) (bool, string) {
				return len(fs) == 0, fmt.Sprintf("Expected 0 flags, found %d", len(fs))
			},
			false},
		{"suspect",
			suspectFlagsRes,
			func(fs []Flag) (bool, string) {
				lenMsg := fmt.Sprintf("Expected 1 flags, found %d", len(fs))

				if len(fs) != 1 {
					return false, lenMsg
				}

				valueMsg := fmt.Sprintf("Expected state: suspect, found %s", fs[0].State)

				if fs[0].State != "suspect" {
					return false, valueMsg
				}

				return true, ""
			},
			false},
		{"cheated",
			cheatedFlagsRes,
			func(fs []Flag) (bool, string) {
				lenMsg := fmt.Sprintf("Expected 1 flags, found %d", len(fs))

				if len(fs) != 1 {
					return false, lenMsg
				}

				valueMsg := fmt.Sprintf("Expected state: cheated, found %s", fs[0].State)

				if fs[0].State != "cheated" {
					return false, valueMsg
				}

				return true, ""
			},
			false},
		{"both",
			bothFlagsRes,
			func(fs []Flag) (bool, string) {
				lenMsg := fmt.Sprintf("Expected 2 flags, found %d", len(fs))

				if len(fs) != 2 {
					return false, lenMsg
				}

				cheatedMsg := fmt.Sprintf("Expected state: cheated, found %s", fs[0].State)
				suspectMsg := fmt.Sprintf("Expected state: suspect, found %s", fs[1].State)

				if fs[0].State != "cheated" {
					return false, cheatedMsg
				}

				if fs[1].State != "suspect" {
					return false, suspectMsg
				}

				return true, ""
			},
			false},
		{"error",
			errorFlagsRes,
			func(fs []Flag) (bool, string) {
				return false, "Should have errored"
			},
			true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httpclient.TestServer{
				Dresponse: tc.res,
			}
			testServerURL := ts.Open()
			defer ts.Close()

			r := httpclient.TestBackend{}
			flags, err := fetchFlags(r, fmt.Sprintf("http://%s", testServerURL), "")

			if tc.shouldErr {
				if err == nil {
					t.Errorf("Something should have gone wrong in the request, but didn't.")
				}
			} else if err != nil {
				t.Errorf("Something went wrong: %s", err)
			}

			isCorrect, errorMsg := tc.checker(flags)

			if !isCorrect && !tc.shouldErr {
				t.Errorf(errorMsg)
			}
		})
	}
}
