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

		result.WriteString(fmt.Sprintf("type %s ", name+parser.AliasTypeSuffix))
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
	comment := ""

	switch def.Kind {
	case parser.NullType, parser.InterfaceType:
		return comment + indent("interface{}", depth)

	case parser.BoolType:
		return comment + indent("bool", depth)

	case parser.StringType:
		return comment + indent("string", depth)

	case parser.IntType:
		return comment + indent("int", depth)

	case parser.FloatType:
		return comment + indent("float64", depth)

	case parser.ArrayType:
		return comment + goArray(p, def, depth)

	case parser.TypeType:
		return comment + indent(def.Name, depth+1)

	case parser.StructType:
		return comment + goStruct(p, def, depth+1)

	case parser.GenericArrayType:
		return comment + indent("[]interface{}", depth)

	default:
		return comment + fmt.Sprintf("###Unsupported type: %v", def.Kind)
	}
}

// Generate an array declaration in Go syntax.
func goArray(p *parser.Parser, def *parser.Type, depth int) string {
	t := def.BaseType
	bt := strings.TrimSpace(goElement(p, t, depth+1))

	return indent("[]"+bt, depth+1)
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

		text := ""

		if t := p.Types[field.Name]; t != nil {
			text = field.Name + parser.AliasTypeSuffix
			if field.Type.Kind == parser.ArrayType {
				text = "[]" + text
			}
		} else {
			text = goElement(p, field.Type, depth+1)
		}

		if len(text) > typeWidth {
			typeWidth = len(text)
		}
	}

	result.WriteString("struct {\n")

	for _, field := range def.Fields {
		result.WriteString(pad("", depth*2))

		result.WriteString(pad(setCase(p, field.Name), nameWidth))
		result.WriteString(" ")

		text := ""

		if t := p.Types[field.Name]; t != nil {
			text = field.Name + parser.AliasTypeSuffix
			if field.Type.Kind == parser.ArrayType {
				text = "[]" + text
			}
		} else {
			text = goElement(p, field.Type, depth+1)
		}

		result.WriteString(pad(strings.TrimSpace(text), typeWidth))
		result.WriteString(tag(p, field))
		result.WriteRune('\n')
	}

	result.WriteString("}\n")

	return result.String()
}

// Generate a JSON tag in Go syntax.
func tag(p *parser.Parser, field *parser.Field) string {
	name := field.Name
	omit := field.Type.Omit

	omitempty := ""
	if p.Omit || omit {
		omitempty = ",omitempty"
	}

	t := fmt.Sprintf("`json:\"%s%s\"`", name, omitempty)

	return "  " + t
}
