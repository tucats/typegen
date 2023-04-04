package parser

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

func (p *Parser) element(item interface{}, depth int) (*Type, error) {
	var err error

	if p.Debug {
		fmt.Printf("[%2d] %sparsing %v\n", depth, strings.Repeat("| ", depth), desc(item))

		if desc(item) == "map[GlossDiv, title]" {
			fmt.Println("DEBUG")
		}
	}

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

func desc(item interface{}) string {
	switch actual := item.(type) {
	case bool:
		return fmt.Sprintf("bool %v", actual)
	case int:
		return fmt.Sprintf("int %v", actual)
	case int64:
		return fmt.Sprintf("int64 %v", actual)
	case float64:
		return fmt.Sprintf("float64 %v", actual)
	case string:
		return fmt.Sprintf("bool %v", strconv.Quote(actual))
	case []interface{}:
		return fmt.Sprintf("array[%d]", len(actual))
	case map[string]interface{}:
		keys := []string{}
		for key := range actual {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		text := strings.Builder{}
		text.WriteString("map[")

		for index, key := range keys {
			if index > 0 {
				text.WriteString(", ")
			}

			text.WriteString(key)
		}

		text.WriteString("]")

		return text.String()

	default:
		return "## Unknown type"
	}
}
