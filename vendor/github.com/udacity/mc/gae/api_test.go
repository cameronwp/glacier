package gae

import (
	"fmt"
	"testing"
)

var (
	nochangeRes = Response{}
	enrolledRes = Response{
		ModifiedEnrollments: []Enrollment{
			{State: "enrolled", NodeKey: "nd050"},
		},
	}
	unenrolledRes = Response{
		ModifiedEnrollments: []Enrollment{
			{State: "unenrolled", NodeKey: "nd050"},
		},
	}
	errorRes = Response{
		ModifiedEnrollments: []Enrollment{
			{State: "enrolled", NodeKey: "nd505"},
		},
	}
)

func TestEnroll(t *testing.T) {
	testCases := []struct {
		res       Response
		expected  bool
		shouldErr bool
	}{
		{nochangeRes, false, false},
		{enrolledRes, true, false},
		{unenrolledRes, false, false},
		{errorRes, false, true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run: %d", i+1), func(t *testing.T) {
			enrolled, err := checkEnrollments(tc.res)
			if tc.shouldErr {
				if err == nil {
					t.Fatal("expected an error but did not find one")
				}
				return
			}

			if err != nil {
				t.Error(err)
			}

			if enrolled != tc.expected {
				t.Errorf("expected 'enrolled: %t', found 'enrolled: %t'", tc.expected, enrolled)
			}
		})
	}
}

func TestUnenroll(t *testing.T) {
	testCases := []struct {
		res       Response
		expected  bool
		shouldErr bool
	}{
		{nochangeRes, false, false},
		{enrolledRes, false, false},
		{unenrolledRes, true, false},
		{errorRes, false, true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run: %d", i+1), func(t *testing.T) {
			unenrolled, err := checkUnenrollments(tc.res)
			if tc.shouldErr {
				if err == nil {
					t.Fatal("expected an error but did not find one")
				}
				return
			}

			if err != nil {
				t.Error(err)
			}

			if unenrolled != tc.expected {
				t.Errorf("expected 'unenrolled: %t', found 'unenrolled: %t'", tc.expected, unenrolled)
			}
		})
	}
}
