package parser

import (
	"fmt"
	"sort"
	"strings"
)

func (p *Parser) structure(data map[string]interface{}, depth int) (*Type, error) {
	var (
		err error
	)

	keys := []string{}
	for key := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	t := newType(StructType)

	for _, key := range keys {
		ft, _ := p.element(data[key], depth+1)

		t.Field(key, ft)
	}

	// If this is the root object, this is all we need.
	if depth == 0 {
		if p.Debug {
			fmt.Printf("[%2d] %s-> naked type\n", depth, strings.Repeat("| ", depth*2))
		}

		return t, err
	}

	// Is this a type we already know about? If so, use that.
	for name, definition := range p.Types {
		if t.Matches(definition) {
			if p.Debug {
				fmt.Printf("[%2d] %s-> apply type %s\n", depth, strings.Repeat("| ", depth*2), name)
			}

			return newType(TypeType).Named(name).Field("base type", t), err
		}
	}

	// Not a type we know about, make a new type
	name := p.generateTypeName()
	p.Types[name] = t
	t = newType(TypeType).Named(name).Field("typed field", t)

	if p.Debug {
		fmt.Printf("[%2d] %s-> created type %s\n", depth, strings.Repeat("| ", depth*2), t)
	}

	return t, err
}
