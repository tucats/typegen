package parser

import "github.com/tucats/typegen/language"

type Parser struct {
	Types  map[string]*Type
	Camel  bool
	Omit   bool
	Debug  bool
	Name   string
	Type   *Type
	Target language.Language
}

func New() *Parser {
	return &Parser{
		Types: map[string]*Type{},
	}
}

func (p *Parser) Language(target language.Language) *Parser {
	p.Target = target

	return p
}
func (p *Parser) OmitEmpty(flag bool) *Parser {
	p.Omit = flag

	return p
}

func (p *Parser) CamelCase(flag bool) *Parser {
	p.Camel = flag

	return p
}

func (p *Parser) Named(name string) *Parser {
	p.Name = name

	return p
}
