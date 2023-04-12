package parser

import (
	"fmt"
	"sort"
	"strings"
)

// Process a JSON object as a structure type (a compound object with named fields).
// The members of the structure are processed recursively and their types are added
// to this structure definition. If the field references another structure, the
// field points to the type definition previously parsed.
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

		if ft.Kind == StructType {
			typeName := p.MakeTypeName(key)

			if t, found := p.Types[typeName]; !found {
				p.Types[typeName] = ft
				if p.Debug {
					fmt.Printf("[%2d] %screate struct type %s as %s\n", depth, strings.Repeat("| ", depth), typeName, ft)
				}

				ft = newType(TypeType).Named(typeName)
			} else {
				if p.Debug {
					fmt.Printf("[%2d] %susing struct type %s as %s\n", depth, strings.Repeat("| ", depth), typeName, t)
				}
				ft = t
			}
		}

		if ft.Kind == ArrayType {
			typeName := p.MakeTypeName(key)
			if t, found := p.Types[typeName]; !found {
				if p.Debug {
					fmt.Printf("[%2d] %screate array type %s as %s\n", depth, strings.Repeat("| ", depth), typeName, ft.BaseType)
				}

				// If the base type is a struct, save it as a type
				if ft.BaseType.Kind == StructType {
					p.Types[typeName] = ft.BaseType
					ft.BaseType = newType(TypeType).Named(typeName)
				}
			} else {
				if p.Debug {
					fmt.Printf("[%2d] %susing struct type %s as %s\n", depth, strings.Repeat("| ", depth), typeName, t)
				}
				ft.BaseType = t
			}
		}

		t.Field(key, ft)
	}

	return t, err
}
