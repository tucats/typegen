package generate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tucats/typegen/parser"
)

// Given a parser, generate the Go version of the definition tree. This
// generates the type elements, and then the root type value.
func generateGo(p *parser.Parser) string {
	result := strings.Builder{}

	// Generate all the type definitions
	keys := []string{}
	for name := range p.Types {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	for _, name := range keys {
		definition := p.Types[name]
		result.WriteString(fmt.Sprintf("type %s ", name))
		result.WriteString(goElement(p, definition, 1))
		result.WriteRune('\n')
	}

	name := p.Name
	if name == "" {
		name = "jsonData"
	}

	result.WriteString(fmt.Sprintf("type %s ", name))
	result.WriteString(goElement(p, p.Type, 1))

	return result.String()
}

// Generate the output for a single element. This will generate defintiions for
// scalar types and recursively generate references to structure fields and array
// types.
func goElement(p *parser.Parser, def *parser.Type, depth int) string {
	switch def.Kind {
	case parser.InterfaceType:
		return indent("interface{}", depth-1)

	case parser.BoolType:
		return indent("bool", depth-1)

	case parser.StringType:
		return indent("string", depth-1)

	case parser.IntType:
		return indent("int", depth-1)

	case parser.FloatType:
		return indent("float64", depth-1)

	case parser.TypeType:
		return indent(def.Name, depth-1)

	case parser.ArrayType:
		return goArray(p, def, depth)

	case parser.StructType:
		return goStruct(p, def, depth)

	default:
		return fmt.Sprintf("###Unsupported type: %v", def.Kind)
	}
}

// Generate an array declaration in Go syntax.
func goArray(p *parser.Parser, def *parser.Type, depth int) string {
	t := def.Fields[0].Type
	bt := strings.TrimSpace(goElement(p, t, depth))

	return indent("[]"+bt, depth-1)
}

// Generate a structure declaration in Go syntax.
func goStruct(p *parser.Parser, def *parser.Type, depth int) string {
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

		t := goElement(p, field.Type, depth+1)
		if len(t) > typeWidth {
			typeWidth = len(t)
		}
	}

	result.WriteString("struct {\n")

	for _, field := range def.Fields {
		result.WriteString(pad("", depth*2))
		result.WriteString(pad(setCase(p, field.Name), nameWidth))
		result.WriteString(" ")

		t := goElement(p, field.Type, depth+1)

		result.WriteString(pad(t, typeWidth))
		result.WriteString(tag(p, field.Name))
		result.WriteRune('\n')
	}

	result.WriteString("}\n")

	return result.String()
}

// Generate a JSON tag in Go syntax.
func tag(p *parser.Parser, name string) string {
	omitempty := ""
	if p.Omit {
		omitempty = ",omitempty"
	}

	t := fmt.Sprintf("`json:\"%s%s\"`", name, omitempty)

	return "  " + t
}
