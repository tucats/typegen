package parser

import (
	"strconv"
	"sync/atomic"
)

var sequenceNumber int32

func (p *Parser) generateTypeName() string {
	name := p.Name
	if name == "" {
		name = "generated"
	}

	n := atomic.AddInt32(&sequenceNumber, 1)

	return name + "Type" + strconv.Itoa(int(n))
}
