package parser

type Parser struct {
	Types map[string]*Type
	Camel bool
	Omit  bool
	Name  string
	Type  *Type
}

func New() *Parser {
	return &Parser{
		Types: map[string]*Type{},
	}
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
