package parser

import (
	"fmt"
	"strings"

	"github.com/tucats/typegen/language"
)

type BaseType int

const (
	InterfaceType BaseType = iota
	BoolType
	IntType
	FloatType
	StringType
	ArrayType
	StructType
	TypeType
)

type Field struct {
	Name string
	Type *Type
}

type Type struct {
	Kind    BaseType
	Name    string
	AltName string
	Fields  []*Field
	Array   *Type
	Omit    bool
}

var AliasTypeSuffix = "Type"

func newType(kind BaseType) *Type {
	return &Type{
		Kind: kind,
	}
}

func (t *Type) Named(name string) *Type {
	t.Name = name

	return t
}

func (t *Type) Field(name string, field *Type) *Type {
	t.Fields = append(t.Fields, &Field{
		Name: name,
		Type: field,
	})

	return t
}

func (t *Type) Alias(name string, depth int, debug bool) *Type {
	oldName := t.Name
	if oldName == "" && t.Kind == ArrayType {
		oldName = t.Fields[0].Type.Name
	}

	if debug {
		fmt.Printf("[%2d] %s-> using alias %s for %s\n", depth, strings.Repeat("| ", depth*2), name, oldName)
	}

	if t.AltName == "" {
		t.AltName = name + AliasTypeSuffix
	}

	return t
}

// String generates a human-readable form of the Type object. This is used
// for debugging purposes only; the output is not generated as any given
// language representation.
func (t *Type) String() string {
	if t == nil {
		return "nil"
	}

	switch t.Kind {
	case BoolType:
		return "bool"

	case IntType:
		return "int"

	case InterfaceType:
		return "interface{}"

	case FloatType:
		return "float64"

	case StringType:
		return "string"

	case ArrayType:
		bt := t.Fields[0].Type

		return "[]" + bt.String()

	case StructType:
		text := strings.Builder{}
		text.WriteString("struct {")

		for index, field := range t.Fields {
			if index > 0 {
				text.WriteString(", ")
			}

			text.WriteString(field.Name)
		}

		text.WriteString("}")

		return text.String()

	case TypeType:
		return fmt.Sprintf("type %s [%s]", t.Name, t.Fields[0].Type)

	default:
		return "unknown type"
	}
}

// Matches determines if the current type and the test type match exactly. This is
// used to determine if the item is a type we have seen before or not. This will
// recursively process compound structure types.
func (t *Type) Matches(test *Type, target language.Language) bool {
	if test == nil {
		return false
	}

	if t.Kind != test.Kind {
		return false
	}

	if t.Kind == ArrayType {
		if len(t.Fields) != 1 || len(test.Fields) != 1 {
			return false
		}

		t1 := t.Fields[0]
		t2 := test.Fields[0]

		if t1 == nil || t2 == nil {
			return false
		}

		return t1.Type.Matches(t2.Type, target)
	}

	if t.Kind == StructType {
		if len(t.Fields) < len(test.Fields) {
			return false
		}

		switch target {
		case language.GoLang, language.Swift:
			// Make a list of fields in the primary structure with their types.
			fields := map[string]*Type{}

			for _, field := range t.Fields {
				fields[field.Name] = field.Type
			}

			// For every field in the test type, it must exist in the current
			// primary type, and have a matching type. If the field in the test
			// does not exist in the primary then it does not match.
			for _, field := range test.Fields {
				if t2, found := fields[field.Name]; found {
					if !t2.Matches(field.Type, target) {
						return false
					} else {
						// Found successfully, remove from field list.
						delete(fields, field.Name)
					}
				} else {
					return false
				}
			}

			// Now, go over the fields that are left (that is, fields that are
			// in the primary but not in the test) and mark them as being able
			// to be omitted.
			for _, fieldType := range fields {
				fieldType.Omit = true
			}

		// By default, to match the fields must be identical in other languages.
		default:
			if len(t.Fields) != len(test.Fields) {
				return false
			}

			for index, t1 := range t.Fields {
				t2 := test.Fields[index]
				if !t1.Type.Matches(t2.Type, target) {
					return false
				}
			}
		}
	}

	if t.Kind == TypeType {
		if !t.Fields[0].Type.Matches(test.Fields[0].Type, target) {
			return false
		}
	}

	return true
}
