package lang

import "strings"

const (
	// NormalizedChinese is the BCP 47 for Chinese.
	NormalizedChinese = "zh-cn"
	// NormalizedEnglish is the BCP 47 for English.
	NormalizedEnglish = "en-us"
	// NormalizedPortuguese is the BCP 47 for Portuguese.
	NormalizedPortuguese = "pt-br"
	// NilLanguage represents a non-supported language.
	NilLanguage = ""
)

// Normalize parses something that may be a language into BCP 47 / INTL Team standards. If this is not a supported language (see Normalized*), returns "".
func Normalize(inlang string) string {
	lowerLang := strings.ToLower(inlang)
	if strings.Index(lowerLang, "zh") == 0 {
		return NormalizedChinese
	} else if strings.Index(lowerLang, "en") == 0 {
		return NormalizedEnglish
	} else if strings.Index(lowerLang, "pt") == 0 {
		return NormalizedPortuguese
	}

	return NilLanguage
}
