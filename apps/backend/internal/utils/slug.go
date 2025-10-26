package utils

import (
	"strings"
	"unicode"
)

func GenerateSlug(input string) string {
	slug := strings.ToLower(strings.TrimSpace(input))
	slug = strings.ReplaceAll(slug, " ", "-")

	var cleaned strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			cleaned.WriteRune(r)
		}
	}
	return cleaned.String()
}
