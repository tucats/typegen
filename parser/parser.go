package parser

import (
	"github.com/tucats/typegen/language"
)

type Parser struct {
	Types      map[string]*Type
	Camel      bool
	Omit       bool
	Debug      bool
	UseAliases bool
	Pretty     bool
	Package    string
	Name       string
	Type       *Type
	Target     language.Language
}

// New creates a new parser object with an empty intialized Type map.
func New() *Parser {
	return &Parser{
		Types: map[string]*Type{},
	}
}

// Language sets the target language for the parse. It returns a pointer
// to the same parser so this call can be chained.
func (p *Parser) Language(target language.Language) *Parser {
	p.Target = target

	return p
}

// OmitEmpty sets the omit flag for the parse. When set, all values of
// structures are marked as "omitempty" in the json tags for go. It
// returns a pointer to the same parser so this call can be chained.
func (p *Parser) OmitEmpty(flag bool) *Parser {
	p.Omit = flag

	return p
}

// CamelCase sets the camel case flag for the parse. It returns a pointer
// to the same parser so this call can be chained.
func (p *Parser) CamelCase(flag bool) *Parser {
	p.Camel = flag

	return p
}

// Named sets the base datatype name for the parse. It returns a pointer
// to the same parser so this call can be chained.
func (p *Parser) Named(name string) *Parser {
	p.Name = name

	return p
}
