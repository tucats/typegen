package parser

import "encoding/json"

// Parse parses a JSON data element, and stores the definition in the current
// parser as the Type element. If the JSON is not valid, that is returned as
// an error condition.
func (p *Parser) Parse(data []byte) error {
	var (
		err     error
		element interface{}
	)

	err = json.Unmarshal(data, &element)
	if err != nil {
		return err
	}

	p.Type, err = p.element(element, 0)

	return err
}
