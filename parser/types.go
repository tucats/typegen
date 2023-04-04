package parser

import (
	"fmt"
	"strings"
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
func (t *Type) Matches(test *Type) bool {
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

		return t1.Type.Matches(t2.Type)
	}

	// @tomcole later, make this allow partial matches
	if t.Kind == StructType {
		if len(t.Fields) != len(test.Fields) {
			return false
		}

		for index, t1 := range t.Fields {
			t2 := test.Fields[index]
			if !t1.Type.Matches(t2.Type) {
				return false
			}
		}
	}

	if t.Kind == TypeType {
		if !t.Fields[0].Type.Matches(test.Fields[0].Type) {
			return false
		}
	}

	return true
}
