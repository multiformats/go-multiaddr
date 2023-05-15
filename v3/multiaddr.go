package multiaddrv3

import (
	"encoding/binary"
	"math"
)

// A Multiaddr is immutable. Think of it like a Go string or netip.AddrPort.
type Multiaddr struct {
	raw        string
	components string // really a []uint16, little endian
}

// Parsed says if the multiaddress could be parsed.
// A multiaddress can only be parsed if we know all the multicodec code points.
// When new multiaddress components are rolled out in the network, we might not know about them yet.
func (m Multiaddr) Parsed() bool { return len(m.components) > 0 }

// Bytes returns the byte representation of the multiaddress.
// The returned byte slice must not be modified.
func (m Multiaddr) Bytes() []byte { return []byte(m.raw) }

// NumComponents returns the number of components of this multiaddress.
// It returns -1 if the multiaddress couldn't be parsed (i.e. Parsed() == false).
// TODO: should it return 0 then?
func (m Multiaddr) NumComponents() int {
	if !m.Parsed() {
		return -1
	}
	return len(m.components) / 2
}

// getIndex converts to an int index to a uint16, checking for overflows
func (m Multiaddr) getIndex(index int) uint16 {
	if index < 0 || index > math.MaxUint16 {
		panic("invalid index")
	}
	idx := 2 * uint16(index)
	if idx >= uint16(len(m.components)) {
		panic("overflow")
	}
	return idx
}

func (m Multiaddr) Component(index int) Component {
	idx := m.getIndex(index)
	var offset uint16
	componentsSlice := []byte(m.components)
	raw := []byte(m.raw)
	offset = binary.LittleEndian.Uint16(componentsSlice[idx:])
	proto64, n := binary.Uvarint(raw[offset:])
	if proto64 > math.MaxUint32 {
		panic("protocol too long")
	}
	offset += uint16(n)
	size := ProtocolWithCode(int(proto64)).Size / 8
	if size > math.MaxInt8 {
		// TODO: add a check like this to the protocol adding function
		panic("size too long")
	}
	if size == LengthPrefixedVarSize {
		l64, n := binary.Uvarint(raw[offset:])
		offset += uint16(n)
		if size > math.MaxUint16 {
			// TODO: do we need a length check here?
			panic("size too long")
		}
		size = int(l64)
	}
	return Component{
		proto: ProtocolCode(proto64),
		val:   m.raw[offset : offset+uint16(size)],
	}
}
