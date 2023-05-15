package multiaddrv3

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"unsafe"
)

// StringCast like Cast, but parses a string. Will also panic if it fails to parse.
func StringCast(s string) Multiaddr {
	m, err := NewMultiaddr(s)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return m
}

// NewMultiaddr parses and validates an input string, returning a *Multiaddr
func NewMultiaddr(s string) (a Multiaddr, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Panic in NewMultiaddr on input %q: %s", s, e)
			err = fmt.Errorf("%v", e)
		}
	}()
	return parseString(s)
}

// NewMultiaddrBytes creates a new multiaddress from a byte slice.
// There's no guarantee that we can actually parse the multiaddress,
// as we might not know all the multicodecs used in the address.
// This function *does not* return an error when that happens.
// The Parsed method can be used to determine if parsing was successful.
// It is safe to modify the byte slice after this function has returned.
func NewMultiaddrBytes(b []byte) (Multiaddr, error) {
	m := Multiaddr{raw: string(b)}

	components := make([]byte, 0, 32)
	var offset uint16
	for len(b) > 0 {
		components = binary.LittleEndian.AppendUint16(components, offset)
		proto64, n := binary.Uvarint(b)
		b = b[n:]
		offset += uint16(n)
		if proto64 > math.MaxInt {
			return m, nil
		}
		proto, ok := protocolsByCode[int(proto64)]
		if !ok {
			// We don't understand this component.
			return m, nil
		}
		var size int
		if proto.Size == LengthPrefixedVarSize {
			size64, n := binary.Uvarint(b)
			b = b[n:]
			offset += uint16(n)
			if size64 > math.MaxUint16 {
				return Multiaddr{}, errors.New("multiaddr: component length overflow")
			}
			size = int(size64)
		} else {
			size = proto.Size / 8
		}
		if size > len(b) {
			return Multiaddr{}, errors.New("multiaddr: not enough bytes")
		}
		offset += uint16(size)
		b = b[size:]
	}
	m.components = byteSliceToString(components)
	return m, nil
}

func parseString(s string) (Multiaddr, error) {
	// consume trailing slashes
	s = strings.TrimRight(s, "/")
	// TODO: don't split here. This allocates like crazy
	sp := strings.Split(s, "/")

	m := Multiaddr{}
	if sp[0] != "" {
		return Multiaddr{}, fmt.Errorf("failed to parse multiaddr %q: must begin with /", s)
	}
	raw := make([]byte, 0, 256)
	components := make([]byte, 0, 32)

	// consume first empty elem
	sp = sp[1:]
	if len(sp) == 0 {
		return Multiaddr{}, fmt.Errorf("failed to parse multiaddr %q: empty multiaddr", s)
	}

	for len(sp) > 0 {
		name := sp[0]
		p := ProtocolWithName(name)
		if p.Code == 0 {
			return Multiaddr{}, fmt.Errorf("failed to parse multiaddr %q: unknown protocol %s", s, sp[0])
		}
		components = binary.LittleEndian.AppendUint16(components, uint16(len(raw)))
		raw = append(raw, p.VCode...)
		sp = sp[1:]

		if p.Size == 0 { // no length.
			continue
		}
		if len(sp) < 1 {
			return Multiaddr{}, fmt.Errorf("failed to parse multiaddr %q: unexpected end of multiaddr", s)
		}

		if p.Path {
			// it's a path protocol (terminal).
			// consume the rest of the address as the next component.
			sp = []string{"/" + strings.Join(sp, "/")}
		}

		a, err := p.Transcoder.StringToBytes(sp[0])
		if err != nil {
			return Multiaddr{}, fmt.Errorf("failed to parse multiaddr %q: invalid value %q for protocol %s: %s", s, sp[0], p.Name, err)
		}
		if p.Size == LengthPrefixedVarSize { // varint size.
			raw = binary.AppendUvarint(raw, uint64(len(a)))
		}
		raw = append(raw, a...)
		sp = sp[1:]
	}
	m.raw = byteSliceToString(raw)
	m.components = byteSliceToString(components)
	return m, nil
}

// byteSliceToString converts a byte slice to a string without allocating memory.
func byteSliceToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
