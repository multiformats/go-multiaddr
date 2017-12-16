package multiaddr

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

// multiaddr is the data structure representing a Multiaddr
type Multiaddr []byte

// NewMultiaddr parses and validates an input string, returning a *Multiaddr
func NewMultiaddr(s string) (a Multiaddr, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Panic in NewMultiaddr on input %q: %s", s, e)
			err = fmt.Errorf("%v", e)
		}
	}()
	b, err := stringToBytes(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// NewMultiaddrBytes initializes a Multiaddr from a byte representation.
// It validates it as an input string but *does not* copy it.
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

	return Multiaddr(b), nil
}

// Equal tests whether two multiaddrs are equal
func (m Multiaddr) Equal(m2 Multiaddr) bool {
	return bytes.Equal(m, m2)
}

// Bytes returns the []byte representation of this Multiaddr
func (m Multiaddr) Bytes() []byte {
	cpy := make([]byte, len(m))
	copy(cpy, m)
	return cpy
}

// String returns the string representation of a Multiaddr
func (m Multiaddr) String() string {
	s, err := bytesToString(m)
	if err != nil {
		panic("multiaddr failed to convert back to string. corrupted?")
	}
	return s
}

// Protocols returns the list of protocols this Multiaddr has.
// will panic in case we access bytes incorrectly.
func (m Multiaddr) Protocols() []Protocol {
	ps := make([]Protocol, 0, 8)
	b := m
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

// Encapsulate wraps a given Multiaddr, returning the resulting joined Multiaddr
func (m Multiaddr) Encapsulate(o Multiaddr) Multiaddr {
	b := make([]byte, len(m)+len(o))
	copy(b, m)
	copy(b[len(m):], o)
	return b
}

// Decapsulate unwraps Multiaddr up until the given Multiaddr is found.
func (m Multiaddr) Decapsulate(o Multiaddr) Multiaddr {
	i := bytes.LastIndex(m, o)
	if i < 0 {
		return m
	}
	return m[:i]
}

var ErrProtocolNotFound = fmt.Errorf("protocol not found in multiaddr")

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
