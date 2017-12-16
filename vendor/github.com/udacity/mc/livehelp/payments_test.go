package livehelp

import (
	"fmt"
	"testing"

	"github.com/udacity/mc/payments"
)

func TestLiveHelpCalculateAmountTier1(t *testing.T) {
	meta := map[string]string{
		"description": "Amounts are per conversation (i.e., 'session'), per activity-level grouping. Activity levels are grouped by their base number of conversations.",
	}
	actLevels := map[string]int{
		"active":           1,
		"very-active":      21,
		"extremely-active": 51,
	}
	amounts := map[string]map[string]float32{
		"beginner": map[string]float32{
			"active":           5,
			"very-active":      7,
			"extremely-active": 9,
		},
		"intermediate": map[string]float32{
			"active":           8,
			"very-active":      10,
			"extremely-active": 12,
		},
		"advanced": map[string]float32{
			"active":           12,
			"very-active":      14,
			"extremely-active": 16,
		},
	}
	lhSessions := &payments.LiveHelpSessions{
		Meta:           meta,
		ActivityLevels: actLevels,
		Amounts:        amounts,
	}

	testCases := []struct {
		ndLevel          string
		sessionCount     int
		livehelpSessions *payments.LiveHelpSessions
		amount           float32
	}{
		{"beginner", 0, lhSessions, 0},
		{"beginner", 1, lhSessions, 5.0},
		{"beginner", 20, lhSessions, 100.0},
		{"beginner", 21, lhSessions, 107.0},
		{"beginner", 50, lhSessions, 310.0},
		{"beginner", 51, lhSessions, 319.0},
		{"beginner", 100, lhSessions, 760.0},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("%s level with %d sessions", tc.ndLevel, tc.sessionCount)
		t.Run(testName, func(t *testing.T) {
			amount := calculateAmount(tc.ndLevel, tc.sessionCount, tc.livehelpSessions)
			if amount != tc.amount {
				t.Errorf("got %f; expected %f", amount, tc.amount)
			}
		})
	}
}

func TestLiveHelpCalculateAmountTier3(t *testing.T) {
	meta := map[string]string{
		"description": "Amounts are per conversation (i.e., 'session'), per activity-level grouping. Activity levels are grouped by their base number of conversations.",
	}
	actLevels := map[string]int{
		"active":           1,
		"very-active":      21,
		"extremely-active": 51,
	}
	amounts := map[string]map[string]float32{
		"beginner": map[string]float32{
			"active":           3.5,
			"very-active":      4.5,
			"extremely-active": 6.2999997,
		},
		"intermediate": map[string]float32{
			"active":           5.6,
			"very-active":      7,
			"extremely-active": 8.4,
		},
		"advanced": map[string]float32{
			"active":           8.4,
			"very-active":      9.8,
			"extremely-active": 11.2,
		},
	}
	lhSessions := &payments.LiveHelpSessions{
		Meta:           meta,
		ActivityLevels: actLevels,
		Amounts:        amounts,
	}

	testCases := []struct {
		ndLevel          string
		sessionCount     int
		livehelpSessions *payments.LiveHelpSessions
		amount           float32
	}{
		{"intermediate", 0, lhSessions, 0},
		{"intermediate", 1, lhSessions, 5.6},
		{"intermediate", 20, lhSessions, 111.999977},
		{"intermediate", 21, lhSessions, 118.999977},
		{"intermediate", 50, lhSessions, 321.999969},
		{"intermediate", 51, lhSessions, 330.399963},
		{"intermediate", 100, lhSessions, 742.000488},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("%s level with %d sessions", tc.ndLevel, tc.sessionCount)
		t.Run(testName, func(t *testing.T) {
			amount := calculateAmount(tc.ndLevel, tc.sessionCount, tc.livehelpSessions)
			if amount != tc.amount {
				t.Errorf("got %f; expected %f", amount, tc.amount)
			}
		})
	}
}
