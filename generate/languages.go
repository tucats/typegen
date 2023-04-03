package generate

import (
	"fmt"

	"github.com/tucats/typegen/parser"
)

type Language int

const (
	GoLang Language = iota
	Swift
)

func Generate(p *parser.Parser, language Language) string {
	switch language {
	case GoLang:
		return generateGo(p)

	case Swift:
		return generateSwift(p)

	default:
		return fmt.Sprintf("Unsupported language: %v", language)
	}
}
