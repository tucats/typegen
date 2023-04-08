package parser

// Process an array. If all the types of the items in the array match, that becomes
// the type of the array.  Otherwise, the type is specified as interface{}.
func (p *Parser) array(data []interface{}, depth int) (*Type, error) {
	var (
		err error
		t   *Type
	)

	for _, value := range data {
		if t == nil {
			t, err = p.element(value, depth+1)
			t.Mergable = true
		} else {
			t2, _ := p.element(value, depth+1)
			t2.Mergable = true
			if !t.Matches(t2) {
				return newType(GenericArrayType), nil
			}
		}
	}

	if t == nil {
		return newType(GenericArrayType), nil
	}

	array := newType(ArrayType)
	array.BaseType = t

	return array, err
}
