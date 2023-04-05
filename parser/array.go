package parser

// Process an array. If all the types of the items in the array match, that becomes
// the type of the array.  Otherwise, the type is specified as interface{}.
func (p *Parser) array(data []interface{}, depth int) (*Type, error) {
	var (
		err error
		t   *Type
	)

	array := newType(ArrayType)

	for _, value := range data {
		if t == nil {
			t, err = p.element(value, depth+1)
		} else {
			t2, _ := p.element(value, depth+1)
			if !t.Matches(t2, p.Target) {
				t = newType(InterfaceType)

				break
			}
		}
	}

	if t == nil {
		t = newType(InterfaceType)
	}

	t.Array = array

	return array.Field("arrayType", t), err
}
