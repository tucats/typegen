package generate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tucats/typegen/parser"
)

// Given a parser, generate the Swift version of the definition tree. This
// generates the type elements, and then the root type value.
func generateSwift(p *parser.Parser) string {
	result := strings.Builder{}

	// Generate all the type definitions
	keys := []string{}
	for name := range p.Types {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	for _, name := range keys {
		definition := p.Types[name]

		if p.UseAliases {
			if alias, found := p.Aliases[name]; found {
				oldName := name
				name = setCase(p, alias) + "Type"

				if p.Debug {
					name = name + " /* " + oldName + " */"
				}
			}
		}

		definition.Name = upcase(name)

		result.WriteString(swiftElement(p, definition, 1))
		result.WriteRune('\n')
	}

	name := p.Name
	if name == "" {
		name = "JsonData"
	}

	p.Type.Name = name
	result.WriteString(swiftElement(p, p.Type, 1))

	return result.String()
}

// Generate the output for a single element. This will generate defintiions for
// scalar types and recursively generate references to structure fields and array
// types.
func swiftElement(p *parser.Parser, def *parser.Type, depth int) string {
	switch def.Kind {
	case parser.BoolType:
		return indent("Boolean", depth-1)

	case parser.StringType:
		return indent("String", depth-1)

	case parser.IntType:
		return indent("Int", depth-1)

	case parser.FloatType:
		return indent("Double", depth-1)

	case parser.TypeType:
		name := def.Name
		if def.AltName != "" {
			name = def.AltName
		}

		return indent(upcase(setCase(p, name)), depth-1)

	case parser.ArrayType:
		return swiftArray(p, def, depth)

	case parser.StructType:
		return swiftStruct(p, def, depth)

	default:
		return fmt.Sprintf("###Unsupported type: %v", def.Kind)
	}
}

// Show an array definition.
func swiftArray(p *parser.Parser, def *parser.Type, depth int) string {
	t := def.Fields[0].Type
	bt := strings.TrimSpace(swiftElement(p, t, depth))

	if depth == 1 {
		return "var jsonData: [" + bt + "]()"
	}

	return indent("["+bt+"]()", depth-1)
}

// Show a structure definition.
func swiftStruct(p *parser.Parser, def *parser.Type, depth int) string {
	result := strings.Builder{}

	nameWidth := 0
	typeWidth := 0

	if def == nil {
		return "## NIL DEF ##"
	}

	for n, field := range def.Fields {
		if field == nil {
			return fmt.Sprintf("## NIL FIELD %d ##", n)
		}

		if len(field.Name) > nameWidth {
			nameWidth = len(field.Name)
		}

		t := swiftElement(p, field.Type, depth+1)
		if len(t) > typeWidth {
			typeWidth = len(t)
		}
	}

	result.WriteString(fmt.Sprintf("class %s: Codable {\n", setCase(p, def.Name)))

	for _, field := range def.Fields {
		result.WriteString(pad("", depth*2))
		result.WriteString("var ")
		result.WriteString(pad(setCase(p, field.Name)+":", nameWidth+3))
		result.WriteString(" ")

		t := swiftElement(p, field.Type, depth+1)

		result.WriteString(pad(t, typeWidth))
		result.WriteRune('\n')
	}

	result.WriteString("}\n")

	return result.String()
}
