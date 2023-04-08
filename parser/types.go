package parser

import (
	"fmt"
	"strings"
)

type BaseType int

const (
	NullType BaseType = iota
	InterfaceType
	BoolType
	IntType
	FloatType
	StringType
	ArrayType
	StructType
	TypeType
	GenericArrayType
	GenericStructType
)

type Field struct {
	Name string
	Type *Type
}

type Type struct {
	Kind     BaseType
	Name     string
	AltName  string
	Fields   []*Field
	BaseType *Type
	Mergable bool
	Omit     bool
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

	optional := ""

	if t.Omit {
		optional = ", optional"
	}

	switch t.Kind {
	case NullType:
		return "nil" + optional

	case GenericArrayType:
		return "[]any" + optional

	case GenericStructType:
		return "struct" + optional

	case BoolType:
		return "bool" + optional

	case IntType:
		return "int" + optional

	case InterfaceType:
		return "interface{}" + optional

	case FloatType:
		return "float64" + optional

	case StringType:
		return "string" + optional

	case ArrayType:
		bt := t.BaseType

		return "[]" + bt.String() + optional

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

		return text.String() + optional

	case TypeType:
		return t.Name

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

	if t.Kind == StructType {
		if t.Mergable || test.Mergable {
			// Make a list of fields in the base type with their types.
			fields := map[string]*Type{}

			for _, field := range t.Fields {
				found := false

				for _, newField := range test.Fields {
					if newField.Name == field.Name {
						found = true

						break
					}
				}

				if !found {
					field.Type.Omit = true
				}

				fields[field.Name] = field.Type
			}

			// For every field in the test type, see if it is in the base type. If so, the
			// types must match. Otherwise, mark is as optional and add it to the base field
			// type list.
			for _, field := range test.Fields {
				if t2, found := fields[field.Name]; found {
					if !t2.Matches(field.Type) {
						return false
					}
				} else {
					field.Type.Omit = true
					fields[field.Name] = field.Type
				}
			}

			// Lastly, copy the field definitions back to the base type.
			for name, fieldType := range fields {
				found := false

				// Try to find an existing definition by name. If found,
				// update the base type with the (possibly updated) omit
				// flag.
				for index, existingField := range t.Fields {
					if existingField.Name == name {
						omit := fieldType.Omit || existingField.Type.Omit
						t.Fields[index].Type.Omit = omit
						found = true

						break
					}
				}

				// If this type wasn't found in the original type, add it with the
				// optional flag set indicating it's not in all instances of the type.
				if !found {
					fieldType.Omit = true
					t.Field(name, fieldType)
				}
			}
		} else {
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
	}

	return true
}
