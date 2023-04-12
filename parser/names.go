package parser

import (
	"strconv"
	"sync/atomic"
)

var sequenceNumber int32

// Generate a sequence number that is guaranteed to be unique and thread-safe.
func (p *Parser) GenerateTypeName() string {
	name := p.Name
	if name == "" {
		name = "jsonData"
	}

	n := atomic.AddInt32(&sequenceNumber, 1)

	return name + AliasTypeSuffix + strconv.Itoa(int(n))
}
