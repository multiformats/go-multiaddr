package multiaddr

import (
	"bytes"
	"fmt"
)

// Split returns the sub-address portions of a multiaddr.
func Split(m Multiaddr) []Multiaddr {
	split, err := bytesSplit(m.Bytes())
	if err != nil {
		panic(fmt.Errorf("invalid multiaddr %s", m.String()))
	}

	addrs := make([]Multiaddr, len(split))
	for i, addr := range split {
		addrs[i] = &multiaddr{bytes: addr}
	}
	return addrs
}

// Join returns a combination of addresses.
func Join(ms ...Multiaddr) Multiaddr {
	var b bytes.Buffer
	for _, m := range ms {
		b.Write(m.Bytes())
	}
	return &multiaddr{bytes: b.Bytes()}
}

// Cast re-casts a byte slice as a multiaddr. will panic if it fails to parse.
func Cast(b []byte) Multiaddr {
	_, err := bytesToString(b)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return &multiaddr{bytes: b}
}

// StringCast like Cast, but parses a string. Will also panic if it fails to parse.
func StringCast(s string) Multiaddr {
	m, err := NewMultiaddr(s)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return m
}
