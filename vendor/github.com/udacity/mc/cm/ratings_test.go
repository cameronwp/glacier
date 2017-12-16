package cm

import (
	"fmt"
	"testing"

	"github.com/udacity/mc/payments"
)

func TestClassMentorCalculateRatingAmountTier1(t *testing.T) {
	meta := map[string]string{
		"description": "Amounts are per rating level (based on number of stars), per student.",
	}
	ratingLevels := map[string]int{
		"highest": 5,
	}
	amounts := map[string]float32{
		"highest": 4,
	}
	cmRatings := &payments.ClassroomMentorshipRatings{
		Meta:         meta,
		RatingLevels: ratingLevels,
		Amounts:      amounts,
	}

	testCases := []struct {
		ratingCount      int
		classMentRatings *payments.ClassroomMentorshipRatings
		amount           float32
	}{
		{0, cmRatings, 0},
		{1, cmRatings, 4.0},
		{10, cmRatings, 40.0},
		{11, cmRatings, 44.0},
		{20, cmRatings, 80.0},
		{100, cmRatings, 400.0},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("%d 5 ratings", tc.ratingCount)
		t.Run(testName, func(t *testing.T) {
			amount := calculateRatingAmount(tc.ratingCount, tc.classMentRatings)
			if amount != tc.amount {
				t.Errorf("got %f; expected %f", amount, tc.amount)
			}
		})
	}
}
