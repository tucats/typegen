package generate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tucats/typegen/parser"
)

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
		definition.Name = name
		result.WriteString(swiftElement(p, definition, 1))
		result.WriteRune('\n')
	}

	name := p.Name
	if name == "" {
		name = "jsonData"
	}

	p.Type.Name = name
	result.WriteString(swiftElement(p, p.Type, 1))

	return result.String()
}

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
		return indent(def.Name, depth-1)

	case parser.ArrayType:
		return swiftArray(p, def, depth)

	case parser.StructType:
		return swiftStruct(p, def, depth)

	default:
		return fmt.Sprintf("###Unsupported type: %v", def.Kind)
	}
}

func swiftArray(p *parser.Parser, def *parser.Type, depth int) string {
	t := def.Fields[0].Type
	bt := strings.TrimSpace(swiftElement(p, t, depth))

	return indent("[]"+bt, depth-1)
}

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

	result.WriteString(fmt.Sprintf("struct %s: Codable {\n", def.Name))

	for _, field := range def.Fields {
		result.WriteString(pad("", depth*2))
		result.WriteString("var ")
		result.WriteString(pad(setCase(field.Name, p.Camel)+":", nameWidth+1))
		result.WriteString(" ")

		t := swiftElement(p, field.Type, depth+1)

		result.WriteString(pad(t, typeWidth))
		result.WriteRune('\n')
	}

	result.WriteString("}\n")

	return result.String()
}
