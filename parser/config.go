package parser

import (
	"fmt"
	"strings"

	"github.com/tucats/typegen/language"
)

type Parser struct {
	Types      map[string]*Type
	Aliases    map[string]string
	Camel      bool
	Omit       bool
	Debug      bool
	UseAliases bool
	Name       string
	Type       *Type
	Target     language.Language
}

func New() *Parser {
	return &Parser{
		Types:   map[string]*Type{},
		Aliases: map[string]string{},
	}
}

func (p *Parser) Alias(t *Type, newname string, depth int) *Parser {
	if p.UseAliases {
		oldname := t.Name
		if oldname == "" && t.Kind == ArrayType {
			oldname = t.Fields[0].Type.Name
		}

		if oldname == "" && t.AltName != "" {
			oldname = t.AltName
		}

		p.Aliases[oldname] = newname

		if p.Debug {
			fmt.Printf("[%2d] %s-> Creating alias %s for %s\n", depth, strings.Repeat("| ", depth*2), oldname, newname)
		}
	}

	return p
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
