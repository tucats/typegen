package generate

import (
	"strings"
	"unicode"

	"github.com/tucats/typegen/language"
	"github.com/tucats/typegen/parser"
)

// Handle name modifications appropriate to the target language and
// settings. This includes setting camel-case if requested for Go code,
// and also handling invalid names in the various languages, which are
// converted to valid name syntax.
//
// For Go, invalid characters are removed and the name is camel-cased
// at each term in the name. For Swift, if the name contains invalid
// characters, it is escaped using the back-tick character.
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

// Indent the given string by the number of "tab stops" in the d parameter. If this
// value is less than 1, no indentention is done and the string is returned unchanged.
// If the value is greater-than or equal-to 1, then the string is returned, prefaced
// by two spaces for each indentation number. That is, a d value of 3 indents by 6
// spaces, etc.
func indent(s string, d int) string {
	if d <= 0 {
		return s
	}

	return strings.Repeat("  ", d) + s
}

// Pad the string to the given width. If the string is already longer than the
// width, the string is returned as-is.  Otherwise, the string has as many blanks
// added to the right side of tne string needed to make the string be the given
// width in characters.
func pad(s string, w int) string {
	for len(s) < w {
		s = s + " "
	}

	return s
}
