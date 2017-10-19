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
	"strings"
)

// Convenience so we don't do stupid things with byte strings.
type bstr string

// Multiaddr is the data structure representing a Multiaddr
// Multiaddrs are immutable and safe to use as map keys.
type Multiaddr struct {
	bytes bstr
}

// NewMultiaddr parses and validates an input string, returning a Multiaddr
func NewMultiaddr(s string) (a Multiaddr, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Panic in NewMultiaddr on input %q: %s", s, e)
			err = fmt.Errorf("%v", e)
		}
	}()
	b, err := stringToBytes(s)
	if err != nil {
		return Multiaddr{}, err
	}
	return Multiaddr{bytes: b}, nil
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
		return Multiaddr{}, err
	}

	return Multiaddr{bytes: bstr(b)}, nil
}

// Equal returns whether two Multiaddrs are exactly equal
func (m Multiaddr) Equal(m2 Multiaddr) bool {
	return m.bytes == m2.bytes
}

// Bytes returns the []byte representation of this Multiaddr
func (m Multiaddr) Bytes() []byte {
	return []byte(m.bytes)
}

// String returns the string representation of this Multiaddr
// (may panic if internal state is corrupted)
func (m Multiaddr) String() string {
	s, err := bytesToString(m.Bytes())
	if err != nil {
		panic("multiaddr failed to convert back to string. corrupted?")
	}
	return s
}

// Protocols returns the list of Protocols this Multiaddr includes
// will panic if protocol code incorrect (and bytes accessed incorrectly)
func (m Multiaddr) Protocols() []Protocol {
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
func (m Multiaddr) Encapsulate(o Multiaddr) Multiaddr {
	return Multiaddr{bytes: m.bytes + o.bytes}
}

// Decapsultate removes a Multiaddr wrapping. For example:
//
//      /ip4/1.2.3.4/tcp/80 decapsulate /ip4/1.2.3.4 = /tcp/80
//
func (m Multiaddr) Decapsulate(o Multiaddr) Multiaddr {
	s1 := m.String()
	s2 := o.String()
	i := strings.LastIndex(s1, s2)
	if i < 0 {
		// Immutable!
		return o
	}

	ma, err := NewMultiaddr(s1[:i])
	if err != nil {
		panic("Multiaddr.Decapsulate incorrect byte boundaries.")
	}
	return ma
}

var ErrProtocolNotFound = fmt.Errorf("protocol not found in multiaddr")

// ValueForProtocol returns the value (if any) following the specified protocol
func (m Multiaddr) ValueForProtocol(code int) (string, error) {
	for _, sub := range Split(m) {
		p := sub.Protocols()[0]
		if p.Code == code {
			if p.Size == 0 {
				return "", nil
			}
			return strings.SplitN(sub.String(), "/", 3)[2], nil
		}
	}

	return "", ErrProtocolNotFound
}
