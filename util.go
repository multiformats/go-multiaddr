package multiaddr

import "fmt"

// Split returns the sub-address portions of a multiaddr.
func Split(m Multiaddr) []Multiaddr {
	s := m.ByteString()
	b := m.Bytes()
	var addrs []Multiaddr
	for len(b) > 0 {
		code, n, err := ReadVarintCode(b)
		if err != nil {
			return nil
		}

		p := ProtocolWithCode(code)
		if p.Code == 0 {
			panic(fmt.Errorf("invalid multiaddr %s, no protocol with code %d", m.String(), b[0]))
		}

		size, err := sizeForAddr(p, b[n:])
		if err != nil {
			panic(fmt.Errorf("invalid multiaddr %s: %s", m.String(), err))
		}

		length := n + size
		addrs = append(addrs, multiaddr(s[:length]))
		b = b[length:]
		s = s[length:]
	}
	return addrs
}

// Join returns a combination of addresses.
func Join(ms ...Multiaddr) Multiaddr {
	var ret string
	for _, m := range ms {
		ret += m.ByteString()
	}

	return multiaddr(ret)
}

// Cast re-casts a byte slice as a multiaddr. will panic if it fails to parse.
func Cast(b []byte) Multiaddr {
	_, err := bytesToString(b)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return multiaddr(b)
}

// StringCast like Cast, but parses a string. Will also panic if it fails to parse.
func StringCast(s string) Multiaddr {
	m, err := NewMultiaddr(s)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return m
}
