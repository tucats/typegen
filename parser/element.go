package parser

import "math"

func (p *Parser) element(item interface{}, depth int) (*Type, error) {
	var err error

	switch actual := item.(type) {
	case string:
		return newType(StringType), nil

	case bool:
		return newType(BoolType), nil

	case int:
		return newType(IntType), nil

	case int64:
		return newType(IntType), nil

	case float64:
		if actual == math.Floor(actual) {
			return newType(IntType), nil
		}

		return newType(FloatType), nil

	case []interface{}:
		return p.array(actual, depth)

	case map[string]interface{}:
		return p.structure(actual, depth)
	}

	return nil, err
}
