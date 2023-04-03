package parser

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
	Kind   BaseType
	Name   string
	Fields []*Field
}

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

	return true
}
