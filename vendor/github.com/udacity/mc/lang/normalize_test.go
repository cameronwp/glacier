package lang

import (
	"fmt"
	"testing"
)

func TestNormalizeSupportedLangs(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"zh-HANS", "zh-cn"},
		{"en-us", "en-us"},
		{"en", "en-us"},
		{"PT", "pt-br"},
		{"PT-BR", "pt-br"},
		{"wtf", ""},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Testing input %s", tc.input), func(t *testing.T) {
			result := Normalize(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, found %s", tc.expected, result)
			}
		})
	}

}
