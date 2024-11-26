package multiaddr

import (
	"fmt"
	"strings"
)

// These are special sizes
const (
	LengthPrefixedVarSize = -1
)

// Protocol is a Multiaddr protocol description structure.
type Protocol struct {
	// Name is the string representation of the protocol code. E.g., ip4,
	// ip6, tcp, udp, etc.
	Name ProtocolName

	// Code is the protocol's multicodec (a normal, non-varint number).
	Code ProtocolCode

	// VCode is a precomputed varint encoded version of Code.
	VCode []byte

	// Size is the size of the argument to this protocol.
	//
	// * Size == 0 means this protocol takes no argument.
	// * Size >  0 means this protocol takes a constant sized argument.
	// * Size <  0 means this protocol takes a variable length, varint
	//             prefixed argument.
	Size int // a size of -1 indicates a length-prefixed variable size

	// Path indicates a path protocol (e.g., unix). When parsing multiaddr
	// strings, path protocols consume the remainder of the address instead
	// of stopping at the next forward slash.
	//
	// Size must be LengthPrefixedVarSize.
	Path bool

	// Transcoder converts between the byte representation and the string
	// representation of this protocol's argument (if any).
	//
	// This should only be non-nil if Size != 0
	Transcoder Transcoder
}

var protocolsByName = map[ProtocolName]Protocol{}
var protocolsByCode = map[ProtocolCode]Protocol{}

// Protocols is the list of multiaddr protocols supported by this module.
var Protocols = []Protocol{}

func AddProtocol(p Protocol) error {
	if _, ok := protocolsByName[p.Name]; ok {
		return fmt.Errorf("protocol by the name %q already exists", p.Name)
	}

	if _, ok := protocolsByCode[p.Code]; ok {
		return fmt.Errorf("protocol code %d already taken by %q", p.Code, p.Code)
	}

	if p.Size != 0 && p.Transcoder == nil {
		return fmt.Errorf("protocols with arguments must define transcoders")
	}
	if p.Path && p.Size >= 0 {
		return fmt.Errorf("path protocols must have variable-length sizes")
	}

	Protocols = append(Protocols, p)
	protocolsByName[p.Name] = p
	protocolsByCode[p.Code] = p
	return nil
}

// protocolWithName returns the Protocol description with given protocol name.
func protocolWithName(s string) Protocol {
	return protocolsByName[ProtocolName(s)]
}

// protocolWithCode returns the Protocol description with given protocol code.
func protocolWithCode(c int) Protocol {
	return protocolsByCode[ProtocolCode(c)]
}

// protocolsWithString returns a slice of protocols matching given string.
func protocolsWithString(s string) ([]Protocol, error) {
	s = strings.Trim(s, "/")
	sp := strings.Split(s, "/")
	if len(sp) == 0 {
		return nil, nil
	}

	t := make([]Protocol, len(sp))
	for i, name := range sp {
		p := protocolWithName(name)
		if p.Code == 0 {
			return nil, fmt.Errorf("no protocol with name: %s", name)
		}
		t[i] = p
	}
	return t, nil
}
