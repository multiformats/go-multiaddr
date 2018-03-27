// Multiaddr is a cross-protocol, cross-platform format for representing
// internet addresses. It emphasizes explicitness and self-description.
// Learn more here: https://github.com/multiformats/multiaddr
//
// Multiaddrs have both a binary and string representation.
//
//     import ma "github.com/multiformats/go-multiaddr"
//
//     addr, err := ma.NewMultiaddr("/ip4/1.2.3.4/tcp/80")
//     // err non-nil when parsing failed.
//

package multiaddr

import (
	"fmt"
	"log"
)

// Multiaddr is the data structure representing a Multiaddr
// Multiaddrs are immutable and safe to use as map keys.
type multiaddr string

// NewMultiaddr parses and validates an input string, returning a Multiaddr
func NewMultiaddr(s string) (a Multiaddr, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Panic in NewMultiaddr on input %q: %s", s, e)
			err = fmt.Errorf("%v", e)
		}
	}()
	return stringToMultiaddr(s)
}

// NewMultiaddrBytes initializes a Multiaddr from a byte representation.
// It validates it as an input string.
func NewMultiaddrBytes(b []byte) (a Multiaddr, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Panic in NewMultiaddrBytes on input %q: %s", b, e)
			err = fmt.Errorf("%v", e)
		}
	}()

	if err := validateBytes(b); err != nil {
		return nil, err
	}

	return multiaddr(b), nil
}

// Equal returns whether two Multiaddrs are exactly equal
func (m multiaddr) Equal(o Multiaddr) bool {
	return m.ByteString() == o.ByteString()
}

// Bytes returns the []byte representation of this Multiaddr
func (m multiaddr) ByteString() string {
	return string(m)
}

// Bytes returns the []byte representation of this Multiaddr
func (m multiaddr) Bytes() []byte {
	return []byte(m)
}

// String returns the string representation of this Multiaddr
// (may panic if internal state is corrupted)
func (m multiaddr) String() string {
	s, err := bytesToString(m.Bytes())
	if err != nil {
		panic("multiaddr failed to convert back to string. corrupted?")
	}
	return s
}

// Protocols returns the list of Protocols this Multiaddr includes
// will panic if protocol code incorrect (and bytes accessed incorrectly)
func (m multiaddr) Protocols() []Protocol {
	ps := make([]Protocol, 0, 8)
	b := m.Bytes()
	for len(b) > 0 {
		code, n, err := ReadVarintCode(b)
		if err != nil {
			panic(err)
		}

		p := ProtocolWithCode(code)
		if p.Code == 0 {
			// this is a panic (and not returning err) because this should've been
			// caught on constructing the Multiaddr
			panic(fmt.Errorf("no protocol with code %d", b[0]))
		}
		ps = append(ps, p)
		b = b[n:]

		size, err := sizeForAddr(p, b)
		if err != nil {
			panic(err)
		}

		b = b[size:]
	}
	return ps
}

// Encapsulate wraps this Multiaddr around another. For example:
//
//      /ip4/1.2.3.4 encapsulate /tcp/80 = /ip4/1.2.3.4/tcp/80
//
func (m multiaddr) Encapsulate(o Multiaddr) Multiaddr {
	return multiaddr(m.ByteString() + o.ByteString())
}

func (m multiaddr) Split(s Multiaddr) (Multiaddr, Multiaddr) {
	haystack := m.ByteString()
	needle := s.ByteString()
	length := len(needle)
	offset := 0

	if length == 0 {
		return nil, nil
	}

	b := m.Bytes()
	for len(haystack) >= offset+length {
		if haystack[offset:offset+length] == needle {
			return multiaddr(haystack[:offset]), multiaddr(haystack[offset+length:])
		}
		code, n, err := ReadVarintCode(b[offset:])
		if err != nil {
			panic(err)
		}

		p := ProtocolWithCode(code)
		if p.Code == 0 {
			// this is a panic (and not returning err) because this should've been
			// caught on constructing the Multiaddr
			panic(fmt.Errorf("no protocol with code %d", code))
		}
		offset += n

		size, err := sizeForAddr(p, b[offset:])
		if err != nil {
			panic(err)
		}
		offset += size
	}
	return nil, nil
}

// Decapsultate removes a Multiaddr wrapping. For example:
//
//      /ip4/1.2.3.4/tcp/80 decapsulate /tcp/80 = /ip4/1.2.3.4
//
func (m multiaddr) Decapsulate(o Multiaddr) Multiaddr {
	if a, _ := m.Split(o); a != nil {
		return a
	}
	return m
}

var ErrProtocolNotFound = fmt.Errorf("protocol not found in multiaddr")

// ValueForProtocol returns the value (if any) following the specified protocol
func (m multiaddr) ValueForProtocol(code int) (string, error) {
	target := ProtocolWithCode(code)
	if target.Code == 0 {
		return "", ErrProtocolNotFound
	}

	b := m.Bytes()
	for offset := 0; offset < len(b); {
		c, n, err := ReadVarintCode(b[offset:])
		if err != nil {
			panic(err)
		}

		p := ProtocolWithCode(c)
		if p.Code == 0 {
			// this is a panic (and not returning err) because this should've been
			// caught on constructing the Multiaddr
			panic(fmt.Errorf("no protocol with code %d", code))
		}
		offset += n

		size, err := sizeForAddr(p, b[offset:])
		if err != nil {
			panic(err)
		}
		if code == c {
			if target.Transcoder == nil {
				return "", nil
			}
			return target.Transcoder.BytesToString(b[offset : offset+size])
		}
		offset += size
	}
	return "", ErrProtocolNotFound
}
