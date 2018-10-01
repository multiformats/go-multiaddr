package multiaddr

import "fmt"

// Split returns the sub-address portions of a multiaddr.
func Split(m Multiaddr) []Multiaddr {
	split, err := bytesSplit(m.Bytes())
	if err != nil {
		panic(fmt.Errorf("invalid multiaddr %s", m.String()))
	}

	addrs := make([]Multiaddr, len(split))
	for i, addr := range split {
		addrs[i] = multiaddr{bytes: addr}
	}
	return addrs
}

// Join returns a combination of addresses.
func Join(ms ...Multiaddr) Multiaddr {
	switch len(ms) {
	case 0:
		// empty multiaddr, unfortunately, we have callers that rely on
		// this contract.
		return multiaddr{}
	case 1:
		return ms[0]
	}

	length := 0
	bs := make([][]byte, len(ms))
	for i, m := range ms {
		bs[i] = m.Bytes()
		length += len(bs[i])
	}

	bidx := 0
	b := make([]byte, length)
	for _, mb := range bs {
		bidx += copy(b[bidx:], mb)
	}
	return multiaddr{bytes: b}
}

// Cast re-casts a byte slice as a multiaddr. will panic if it fails to parse.
func Cast(b []byte) Multiaddr {
	m, err := NewMultiaddrBytes(b)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return m
}

// StringCast like Cast, but parses a string. Will also panic if it fails to parse.
func StringCast(s string) Multiaddr {
	m, err := NewMultiaddr(s)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return m
}
