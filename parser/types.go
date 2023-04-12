package parser

import (
	"strings"
	"unicode"

	"github.com/tucats/typegen/language"
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

// Field describes a field in a structure (a JSON object). This includes
// the name of the field, and it's type designation.
type Field struct {
	Name string
	Type *Type
}

// Type describes everything known about a type designation. This could be an
// elemental type like a string or a bool.  It could also be an array or
// structure, or a type designation from the named type dictionary.
type Type struct {
	// Kind is the base type of this type, from the list of possible types.
	Kind BaseType

	// Name is the name of the type. This is only used for TypeType designations.
	// For all other types, this value is an emptystring.
	Name string

	// Fields is a list of the fields if this is a structure/object designation.
	// For all other types, this value is nil.
	Fields []*Field

	// BaseType of this type. This is used only for ArrayType, and tells the type
	// of each member of the array (if it is a homogeneous array) or GeneritArrayType
	// if the array is heterogeneous.  For all other types, this value is nil.
	BaseType *Type

	// Flag set indicating that this type can accept additional fields added to it
	// when it is a struct/object that is part of an array. The match operation will
	// merge types when this flag is set when comparing compatable structure/object
	// designations.
	Mergable bool

	// Indicates that, for Go language generation, this type should include the
	// "omitempty" flag in the json tag.
	Omit bool
}

// AliasTypeSuffix is the string added to the end of type names that are inferred from
// structure/object field names.
var AliasTypeSuffix = "Type"

// Create a new type object of the given base type.
func newType(kind BaseType) *Type {
	return &Type{
		Kind: kind,
	}
}

// Named sets the name of the type designation. This returns a pointer to the type
// value so this call can be chained.
func (t *Type) Named(name string) *Type {
	t.Name = name

	return t
}

// Field adds a field signation (name and type) to the current type's list of fields.
// This returns a pointer to the type value so this call can be chained.
func (t *Type) Field(name string, field *Type) *Type {
	t.Fields = append(t.Fields, &Field{
		Name: name,
		Type: field,
	})

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
		t1 := t.BaseType
		t2 := test.BaseType

		if t1 == nil || t2 == nil {
			return false
		}

		return t1.Matches(t2)
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

// MakeTypeName generates a type name from a field name. This includes modifying
// the field name to ensure it is a valid name in the target language, and adding
// the suffix to the type name.
func (p *Parser) MakeTypeName(name string) string {
	result := strings.Builder{}
	escape := false

	for index, ch := range name {
		if index == 0 && p.Target == language.GoLang {
			ch = unicode.ToUpper(ch)
			if !unicode.IsLetter(ch) {
				ch = 'X'
			}
		}

		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			switch p.Target {
			case language.GoLang:
				ch = '_'

			case language.Swift:
				escape = true
			}
		}

		result.WriteRune(ch)
	}

	result.WriteString(AliasTypeSuffix)

	text := result.String()
	if escape {
		text = "`" + text + "`"
	}

	return text
}
