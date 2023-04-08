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

		result.WriteString(fmt.Sprintf("class %s: Codable ", name+parser.AliasTypeSuffix))
		result.WriteString(swiftElement(p, definition, 1))
		result.WriteRune('\n')
	}

	name := p.Name
	if name == "" {
		name = "jsonData"
	}

	result.WriteString(fmt.Sprintf("class %s: Codable ", name))
	result.WriteString(swiftElement(p, p.Type, 1))

	return result.String()
}

// Generate the output for a single element. This will generate defintiions for
// scalar types and recursively generate references to structure fields and array
// types.
func swiftElement(p *parser.Parser, def *parser.Type, depth int) string {
	switch def.Kind {
	case parser.BoolType:
		return indent("Bool", depth-1)

	case parser.StringType:
		return indent("String", depth-1)

	case parser.IntType:
		return indent("Int", depth-1)

	case parser.FloatType:
		return indent("Double", depth-1)

	case parser.InterfaceType, parser.GenericArrayType:
		return "### NON-HOMOGENEOUS ARRAY MEMBERS ###"

	case parser.ArrayType:
		return swiftArray(p, def, depth)

	case parser.StructType, parser.GenericStructType:
		return swiftStruct(p, def, depth)

	default:
		return fmt.Sprintf("###Unsupported type: %v", def.Kind)
	}
}

// Show an array definition.
func swiftArray(p *parser.Parser, def *parser.Type, depth int) string {
	t := def.BaseType

	if t.Kind == parser.InterfaceType {
		return "  Bool? /* Invalid heterogeneous array in input data */"
	}

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

	result.WriteString(" {\n")

	for _, field := range def.Fields {
		result.WriteString(pad("", depth*2))
		result.WriteString("var ")

		if p.Pretty {
			result.WriteString(pad(setCase(p, field.Name)+":", nameWidth+3))
		} else {
			result.WriteString(setCase(p, field.Name) + ": ")
		}

		text := ""
		if t := p.Types[field.Name]; t != nil {
			text = field.Name + parser.AliasTypeSuffix
			if field.Type.Kind == parser.ArrayType {
				text = "[" + text + "]()"
			}
		} else {
			text = swiftElement(p, field.Type, depth+1)
		}

		if field.Type.Omit {
			text = text + "?" // Mark as optional
		}

		if !p.Pretty {
			text = strings.TrimSpace(text)
		}

		result.WriteString(pad(text, typeWidth))
		result.WriteRune('\n')
	}

	result.WriteString("}\n")

	return result.String()
}
