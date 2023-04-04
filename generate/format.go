package generate

import (
	"strings"
	"unicode"

	"github.com/tucats/typegen/language"
	"github.com/tucats/typegen/parser"
)

func setCase(p *parser.Parser, name string) string {
	result := strings.Builder{}
	camelNext := false
	backTick := ""

	for index, ch := range name {
		switch p.Target {
		case language.Swift:
			// If it's a risky name in Swift, let's escape it in the resulting code.
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
				backTick = "`"
			}

			result.WriteRune(ch)

		case language.GoLang:
			if index == 0 {
				if p.Camel {
					ch = unicode.ToLower(ch)
				} else {
					ch = unicode.ToUpper(ch)
				}
			}

			// If it is an invalid character in Go, drop it and camel case the next part.
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
				camelNext = true

				continue
			}

			if camelNext {
				ch = unicode.ToUpper(ch)
				camelNext = false
			}

			result.WriteRune(ch)
		}
	}

	return backTick + result.String() + backTick
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
