package multiaddr

import "fmt"

// Split returns the sub-address portions of a multiaddr.
func Split(m Multiaddr) []Multiaddr {
	b := []byte(m)
	var ret []Multiaddr
	for len(b) > 0 {
		code, n, err := ReadVarintCode(b)
		if err != nil {
			panic(fmt.Errorf("invalid multiaddr %s", m.String()))
		}

		p := ProtocolWithCode(code)
		if p.Code == 0 {
			panic(fmt.Errorf("invalid multiaddr %s", m.String()))
		}

		size, err := sizeForAddr(p, b[n:])
		if err != nil {
			panic(fmt.Errorf("invalid multiaddr %s", m.String()))
		}

		length := n + size
		ret = append(ret, b[:length])
		b = b[length:]
	}

	return ret
}

// Join returns a combination of addresses.
func Join(ms ...Multiaddr) Multiaddr {
	length := 0
	for _, m := range ms {
		length += len(m)
	}

	offset := 0
	out := make([]byte, length)
	for _, m := range ms {
		offset += copy(out[offset:], m)
	}
	return out
}

// Cast re-casts a byte slice as a multiaddr. will panic if it fails to parse.
func Cast(b []byte) Multiaddr {
	_, err := bytesToString(b)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return b
}

// StringCast like Cast, but parses a string. Will also panic if it fails to parse.
func StringCast(s string) Multiaddr {
	m, err := NewMultiaddr(s)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return m
}
