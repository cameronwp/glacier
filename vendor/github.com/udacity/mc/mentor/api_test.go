package mentor

import (
	"fmt"
	"regexp"
	"testing"
)

func TestBuildUpdateMentorQuery(t *testing.T) {
	testCases := []struct {
		fields     map[string]string
		expectedRe string
	}{
		{
			map[string]string{
				"uid": "1234asdf",
			},
			`\(uid: "1234asdf"\)`,
		},
		{
			map[string]string{
				"uid":     "1234asdf",
				"country": "CA",
			},
			`\(uid: "1234asdf", country: "CA"\)`,
		},
		{
			map[string]string{
				"uid":          "1234asdf",
				"paypal_email": "payme@now.com",
				"country":      "CA",
			},
			`\(uid: "1234asdf", country: "CA", paypal_email: "payme@now.com"\)`,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			query := buildUpdateMentorQuery(tc.fields)
			re := regexp.MustCompile(tc.expectedRe)

			result := re.FindStringSubmatch(query)

			if result == nil {
				t.Errorf("Expected to find the query variables, but did not in:\n%s", query)
			}
		})
	}
}
