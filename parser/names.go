package parser

import (
	"strconv"
	"sync/atomic"
)

var sequenceNumber int32

// Generate a sequence number that is guaranteed to be unique and thread-safe.
func (p *Parser) generateTypeName() string {
	name := p.Name
	if name == "" {
		name = "generated"
	}

	n := atomic.AddInt32(&sequenceNumber, 1)

	return name + "Type" + strconv.Itoa(int(n))
}
