package generate

import (
	"fmt"

	"github.com/tucats/typegen/language"
	"github.com/tucats/typegen/parser"
)

// Generate produces a text representation of the declaration, given
// a target language such as language.GoLang or language.Swift.
func Generate(p *parser.Parser, target language.Language) string {
	switch target {
	case language.GoLang:
		return generateGo(p)

	case language.Swift:
		return generateSwift(p)

	default:
		return fmt.Sprintf("Unsupported language: %v", target)
	}
}
