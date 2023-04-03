package parser

import "encoding/json"

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
