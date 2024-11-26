package multiaddr

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

// Protocols is the list of multiaddr protocols supported by this module.
var Protocols = []Protocol{}
