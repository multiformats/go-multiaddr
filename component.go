package multiaddr

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// component is a single multiaddr component.
type component struct {
	bytes    []byte
	protocol Protocol
	offset   int
}

func (c *component) Bytes() []byte {
	return c.bytes
}

func (c *component) Equal(o Multiaddr) bool {
	return bytes.Equal(c.bytes, o.Bytes())
}

func (c *component) Protocols() []Protocol {
	return []Protocol{c.protocol}
}

func (c *component) Decapsulate(o Multiaddr) Multiaddr {
	if c.Equal(o) {
		return nil
	}
	return c
}

func (c *component) Encapsulate(o Multiaddr) Multiaddr {
	m := multiaddr{bytes: c.bytes}
	return m.Encapsulate(o)
}

func (c *component) ValueForProtocol(code int) (string, error) {
	if c.protocol.Code != code {
		return "", ErrProtocolNotFound
	}
	return c.Value(), nil
}

func (c *component) ForEach(cb func(c Component) bool) {
	cb(c)
}

func (c *component) Protocol() Protocol {
	return c.protocol
}

func (c *component) RawValue() []byte {
	return c.bytes[c.offset:]
}

func (c *component) Value() string {
	if c.protocol.Transcoder == nil {
		return ""
	}
	value, err := c.protocol.Transcoder.BytesToString(c.bytes[c.offset:])
	if err != nil {
		// This component must have been checked.
		panic(err)
	}
	return value
}

func (c *component) String() string {
	str := "/" + c.protocol.Name

	value := c.Value()
	if len(value) == 0 {
		return str
	}

	if !(c.protocol.Path && value[0] == '/') {
		str += "/"
	}

	return str + value
}

// NewComponent constructs a new multiaddr component
func NewComponent(protocol, value string) (Component, error) {
	p := ProtocolWithName(protocol)
	if p.Code == 0 {
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
	if p.Transcoder != nil {
		bts, err := p.Transcoder.StringToBytes(value)
		if err != nil {
			return nil, err
		}
		return newComponent(p, bts), nil
	} else if value != "" {
		return nil, fmt.Errorf("protocol %s doesn't take a value", p.Name)
	}
	return newComponent(p, nil), nil
	// TODO: handle path /?
}

func newComponent(protocol Protocol, bvalue []byte) *component {
	size := len(bvalue)
	size += len(protocol.VCode)
	if protocol.Size < 0 {
		size += VarintSize(len(bvalue))
	}
	maddr := make([]byte, size)
	var offset int
	offset += copy(maddr[offset:], protocol.VCode)
	if protocol.Size < 0 {
		offset += binary.PutUvarint(maddr[offset:], uint64(len(bvalue)))
	}
	copy(maddr[offset:], bvalue)

	// For debugging
	if len(maddr) != offset+len(bvalue) {
		panic("incorrect length")
	}

	return &component{
		bytes:    maddr,
		protocol: protocol,
		offset:   offset,
	}
}
