package generate

import (
	"strings"
	"unicode"
)

func setCase(name string, lowercase bool) string {
	result := strings.Builder{}

	for index, ch := range name {
		if index == 0 {
			if lowercase {
				ch = unicode.ToLower(ch)
			} else {
				ch = unicode.ToUpper(ch)
			}
		}

		result.WriteRune(ch)
	}

	return result.String()
}

func indent(s string, d int) string {
	return strings.Repeat("  ", d) + s
}

func pad(s string, w int) string {
	for len(s) < w {
		s = s + " "
	}

	return s
}
