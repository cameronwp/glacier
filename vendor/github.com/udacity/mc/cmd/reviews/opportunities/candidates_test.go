package opportunities

import (
	"fmt"
	"testing"
	"time"

	"github.com/udacity/mc/cc"
	"github.com/udacity/mc/flagging"
	"github.com/udacity/mc/gae"
	"github.com/udacity/mc/mentor"
	"github.com/udacity/mc/reviews"
)

var mentorsT = []mentor.Mentor{
	mentor.Mentor{
		UID:       "1",
		Languages: []string{"en-us", "pt-br"},
		Application: mentor.Application{
			Nanodegrees: []string{"nd001", "nd002"},
			Services:    []string{"reviews", "guru"},
		},
	},
	mentor.Mentor{
		UID:       "2",
		Languages: []string{"EN", "zh-cn"},
		Application: mentor.Application{
			Nanodegrees: []string{"nd001", "nd003"},
			Services:    []string{"reviews"},
		},
	},
	mentor.Mentor{
		UID:       "3",
		Languages: []string{},
		Application: mentor.Application{
			Nanodegrees: []string{"nd001", "nd003"},
			Services:    []string{"reviews"},
		},
		CreatedAt: time.Now(),
	},
}

func TestMapToCandidates(t *testing.T) {
	candidates := mapToCandidates(mentorsT)

	if len(candidates) != 3 {
		t.Errorf("Expected to find 3 candidates, but found %d", len(candidates))
	}

	if candidates[0].UID != "1" {
		t.Errorf("Expected the first candidate to have a UID of 1, found %s", candidates[0].UID)
	}

	if candidates[2].CreatedAt != mentorsT[2].CreatedAt {
		t.Errorf("Expected created timestamp to carry over. got %s but expected %s", candidates[2].CreatedAt, mentorsT[2].CreatedAt)
	}
}

func TestFilterCandidates(t *testing.T) {
	candidates := filterCandidates(mentorsT, "nd001", "reviews", []string{"en-us"})
	if len(candidates) != 2 {
		t.Errorf("Expected 2 candidates to make it, found %d instead", len(candidates))
	}

	t.Run("when only one person has signed up for the ND", func(t *testing.T) {
		candidates := filterCandidates(mentorsT, "nd002", "reviews", []string{"en-us"})
		if len(candidates) != 1 {
			t.Fatalf("Expected to filter and get 1 candidate, found %d", len(candidates))
		}

		if candidates[0].UID != "1" {
			t.Errorf("Expected to find mentor 1, found mentor %s", candidates[0].UID)
		}
	})

	t.Run("when no one has signed up for the ND", func(t *testing.T) {
		candidates := filterCandidates(mentorsT, "nd999", "reviews", []string{"en-us"})
		if len(candidates) != 0 {
			t.Errorf("Expected not to find any candidates, found %d instead", len(candidates))
		}
	})

	t.Run("with more than one language", func(t *testing.T) {
		candidates := filterCandidates(mentorsT, "nd001", "reviews", []string{"pt-br", "zh-cn"})
		if len(candidates) != 2 {
			t.Errorf("Expected to find 2 candidates, found %d instead", len(candidates))
		}
	})
}

func TestCandidateBio(t *testing.T) {
	t.Run("When we have their email", func(t *testing.T) {
		c := candidate{}

		user := gae.User{
			Email: gae.Email{
				Address: "test@test.com",
			},
		}

		markCandidateBio(&c, user)
		if c.Email != user.Email.Address {
			t.Errorf("expected email %s, got %s", user.Email.Address, c.Email)
		}
	})
}

func TestCandidateCertificate(t *testing.T) {
	t.Run("With a good certs response", func(t *testing.T) {
		c := candidate{
			UID: mentorsT[0].UID,
		}

		certs := []reviews.Certification{
			{
				ProjectID: 190,
				Status:    "certified",
			},
		}

		projectID = "190"

		markCandidateCerts(&c, certs)
		if c.CertificationStatus != "certified" {
			t.Errorf("Expected certified status, found: %s", c.CertificationStatus)
		}
	})
}

func TestCandidateEnrollment(t *testing.T) {
	testCases := []struct {
		user                     cc.User
		expectedEnrolled         bool
		expectedGraduated        bool
		expectedCompletionAmount float64
	}{
		{
			cc.User{
				Nanodegrees: []cc.Nanodegree{},
			},
			false,
			false,
			0.0,
		},
		{
			cc.User{
				Nanodegrees: []cc.Nanodegree{
					{
						IsGraduated: true,
						AggregatedState: cc.AggregatedState{
							CompletionAmount: 1.0,
						},
					},
				},
			},
			true,
			true,
			1.0,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run: %d", i), func(t *testing.T) {
			c := candidate{}

			markCandidateEnrollment(&c, tc.user)

			if c.IsEnrolled != tc.expectedEnrolled {
				t.Errorf("Expected enrollment %t, found %t", tc.expectedEnrolled, c.IsEnrolled)
			}

			if c.IsGraduated != tc.expectedGraduated {
				t.Errorf("Expected graduated %t, found %t", tc.expectedGraduated, c.IsGraduated)
			}

			if c.CompletionAmount != tc.expectedCompletionAmount {
				t.Errorf("Expected completed %f, found %f", tc.expectedCompletionAmount, c.CompletionAmount)
			}
		})
	}
}

func TestCandidateFlags(t *testing.T) {
	testCases := []struct {
		flags    []flagging.Flag
		expected string
		name     string
	}{
		{[]flagging.Flag{flagging.Flag{}}, "none", "no flags"},
		{
			[]flagging.Flag{
				flagging.Flag{
					State: "suspect",
				},
			},
			"suspect", "only suspected",
		},
		{
			[]flagging.Flag{
				flagging.Flag{
					State: "cheated",
				},
			},
			"cheated", "only cheated",
		},
		{
			[]flagging.Flag{
				flagging.Flag{
					State: "cheated",
				},
				flagging.Flag{
					State: "suspect",
				},
			},
			"cheated", "suspected and cheated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := candidate{}
			markCandidateFlags(&c, tc.flags)

			result := c.Flagged
			if result != tc.expected {
				t.Errorf("Expected %s, found %s", tc.expected, result)
			}
		})
	}
}

func TestCandidateOpportunities(t *testing.T) {
	testCases := []struct {
		opportunities []reviews.Opportunity
		expectedCount int
		expectedOpen  bool
	}{
		{
			[]reviews.Opportunity{},
			0, false,
		},
		{
			[]reviews.Opportunity{
				reviews.Opportunity{
					ExpiresAt: time.Now().AddDate(0, 0, 1),
				},
			},
			1, true,
		},
		{
			[]reviews.Opportunity{
				reviews.Opportunity{
					ExpiresAt: time.Now().AddDate(0, 0, -2),
				},
				reviews.Opportunity{
					ExpiresAt: time.Now().AddDate(0, 0, -1),
				},
			},
			2, false,
		},
		{
			[]reviews.Opportunity{
				reviews.Opportunity{
					ExpiresAt: time.Now().AddDate(0, 0, -2),
				},
				reviews.Opportunity{
					ExpiresAt: time.Now().AddDate(0, 0, 1),
				},
			},
			2, true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			c := candidate{}

			markCandidateOpportunities(&c, tc.opportunities)

			if c.PrevOpportunityCount != tc.expectedCount {
				t.Errorf("Expected count %d, found %d", tc.expectedCount, c.PrevOpportunityCount)
			}

			if c.HasOpenOpportunity != tc.expectedOpen {
				t.Errorf("Expected open %t, found %t", tc.expectedOpen, c.HasOpenOpportunity)
			}
		})
	}
}

func TestCandidateSubmissions(t *testing.T) {
	testCases := []struct {
		submissions    []reviews.Submission
		expectedCount  int
		expectedPassed bool
	}{
		{
			[]reviews.Submission{},
			0, false,
		},
		{
			[]reviews.Submission{
				reviews.Submission{
					Result: "passed",
				},
			},
			1, true,
		},
		{
			[]reviews.Submission{
				reviews.Submission{
					Result: "passed",
				},
				reviews.Submission{
					Result: "failed",
				},
			},
			2, true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			c := candidate{}

			markCandidateSubmissions(&c, tc.submissions)

			if tc.expectedCount != c.SubmissionCount {
				t.Errorf("Expected count %d, found %d", tc.expectedCount, c.SubmissionCount)
			}

			if tc.expectedPassed != c.PassedProject {
				t.Errorf("Expected passed %t, found %t", tc.expectedPassed, c.PassedProject)
			}
		})
	}
}

func TestPruneCertedCandidates(t *testing.T) {
	testCases := []struct {
		candidates []candidate
		expected   int
	}{
		{
			[]candidate{
				{CertificationStatus: "certified"},
				{CertificationStatus: "training"},
				{},
			},
			1,
		},
		{
			[]candidate{
				{CertificationStatus: "blocked"},
				{},
			},
			1,
		},
		{
			[]candidate{
				{},
				{},
			},
			2,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run %d", i+1), func(t *testing.T) {
			candidates := pruneCertedCandidates(tc.candidates)
			if len(candidates) != tc.expected {
				t.Errorf("expected %d candidates left, found %d", tc.expected, len(candidates))
			}
		})
	}
}

func TestBuildCSVRows(t *testing.T) {
	// a few spot checks to sanity check that the columns are in the right order
	now := time.Now()
	candidates := []candidate{
		candidate{
			ProjectID:  "123",
			UID:        "1",
			IsEnrolled: true,
			CreatedAt:  now,
			Flagged:    "cheated",
		},
		candidate{
			UID:             "2",
			UpdatedAt:       now,
			SubmissionCount: 5,
		},
	}

	testCases := []struct {
		row      int
		col      int
		name     string
		expected string
	}{
		{0, 0, "project id", "123"},
		{0, 1, "uid", "1"},
		{0, 7, "flagged", "cheated"},
		{0, 16, "created_at", now.Local().Format(time.RFC822)},
		{1, 17, "updated_at", now.Local().Format(time.RFC822)},
		{1, 12, "submission_count", "5"},
	}

	_, rows := buildCSVRows(candidates)

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run %d", i+1), func(t *testing.T) {
			if rows[tc.row][tc.col] != tc.expected {
				t.Errorf("expected %s in (%d, %d) to be %s, got %s", tc.name, tc.row, tc.col, tc.expected, rows[tc.row][tc.col])
			}
		})
	}
}
