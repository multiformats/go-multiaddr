package multiaddr

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/multiformats/go-varint"
)

func stringToBytes(s string) ([]byte, error) {
	// consume trailing slashes
	s = strings.TrimRight(s, "/")

	var b bytes.Buffer
	sp := strings.Split(s, "/")

	if sp[0] != "" {
		return nil, fmt.Errorf("failed to parse multiaddr %q: must begin with /", s)
	}

	// consume first empty elem
	sp = sp[1:]

	if len(sp) == 0 {
		return nil, fmt.Errorf("failed to parse multiaddr %q: empty multiaddr", s)
	}

	for len(sp) > 0 {
		name := sp[0]
		p := ProtocolWithName(name)
		if p.Code == 0 {
			return nil, fmt.Errorf("failed to parse multiaddr %q: unknown protocol %s", s, sp[0])
		}
		_, _ = b.Write(p.VCode)
		sp = sp[1:]

		if p.Size == 0 { // no length.
			continue
		}

		if len(sp) < 1 {
			return nil, fmt.Errorf("failed to parse multiaddr %q: unexpected end of multiaddr", s)
		}

		if p.Path {
			// it's a path protocol (terminal).
			// consume the rest of the address as the next component.
			sp = []string{"/" + strings.Join(sp, "/")}
		}

		a, err := p.Transcoder.StringToBytes(sp[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse multiaddr %q: invalid value %q for protocol %s: %s", s, sp[0], p.Name, err)
		}
		err = p.Transcoder.ValidateBytes(a)
		if err != nil {
			return nil, err
		}
		if p.Size < 0 { // varint size.
			_, _ = b.Write(varint.ToUvarint(uint64(len(a))))
		}
		b.Write(a)
		sp = sp[1:]
	}

	return b.Bytes(), nil
}

func readComponent(b []byte) (int, Component, error) {
	var offset int
	code, n, err := ReadVarintCode(b)
	if err != nil {
		return 0, Component{}, err
	}
	offset += n

	p := ProtocolWithCode(code)
	if p.Code == 0 {
		return 0, Component{}, fmt.Errorf("no protocol with code %d", code)
	}

	if p.Size == 0 {
		c, err := validateComponent(Component{
			bytes:    string(b[:offset]),
			offset:   offset,
			protocol: p,
		})

		return offset, c, err
	}

	n, size, err := sizeForAddr(p, b[offset:])
	if err != nil {
		return 0, Component{}, err
	}

	offset += n

	if len(b[offset:]) < size || size < 0 {
		return 0, Component{}, fmt.Errorf("invalid value for size %d", len(b[offset:]))
	}

	c, err := validateComponent(Component{
		bytes:    string(b[:offset+size]),
		protocol: p,
		offset:   offset,
	})

	return offset + size, c, err
}

func readMultiaddr(b []byte) (int, Multiaddr, error) {
	if len(b) == 0 {
		return 0, Multiaddr{}, fmt.Errorf("empty multiaddr")
	}

	var res Multiaddr
	bytesRead := 0
	for len(b) > 0 {
		n, c, err := readComponent(b)
		if err != nil {
			return 0, Multiaddr{}, err
		}
		b = b[n:]
		bytesRead += n
		res = append(res, c)
	}
	return bytesRead, res, nil
}

func bytesToString(b []byte) (ret string, err error) {
	if len(b) == 0 {
		return "", fmt.Errorf("empty multiaddr")
	}
	var buf strings.Builder

	for len(b) > 0 {
		n, c, err := readComponent(b)
		if err != nil {
			return "", err
		}
		b = b[n:]
		c.writeTo(&buf)
	}

	return buf.String(), nil
}

func sizeForAddr(p Protocol, b []byte) (skip, size int, err error) {
	switch {
	case p.Size > 0:
		return 0, (p.Size / 8), nil
	case p.Size == 0:
		return 0, 0, nil
	default:
		size, n, err := ReadVarintCode(b)
		if err != nil {
			return 0, 0, err
		}
		return n, size, nil
	}
}
