package nd

import (
	"fmt"
	"testing"

	"github.com/udacity/mc/gae"
)

func TestFilterForND(t *testing.T) {
	testCases := []struct {
		enrollments []gae.Enrollment
		nodeKey     string
		expectedLen int
	}{
		{
			[]gae.Enrollment{
				{NodeKey: "nd050"},
			},
			"nd050", 1,
		},
		{
			[]gae.Enrollment{
				{NodeKey: "nd050"},
				{NodeKey: "nd001"},
			},
			"nd050", 1,
		},
		{
			[]gae.Enrollment{
				{NodeKey: "nd050"},
				{NodeKey: "nd050"},
			},
			"nd050", 2,
		},
		{
			[]gae.Enrollment{
				{NodeKey: "nd001"},
				{NodeKey: "nd002"},
			},
			"nd050", 0,
		},
		{
			[]gae.Enrollment{},
			"nd050", 0,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			filterForND(&tc.enrollments, tc.nodeKey)

			if len(tc.enrollments) != tc.expectedLen {
				t.Errorf("expected %d enrollments left, found %d", tc.expectedLen, len(tc.enrollments))
			}

			for _, e := range tc.enrollments {
				if e.NodeKey != tc.nodeKey {
					t.Errorf("expected node key of %s, found %s", tc.nodeKey, e.NodeKey)
				}
			}
		})
	}
}
