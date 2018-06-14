package multiaddr

import (
	"bytes"
	"fmt"
	"strings"
)

func stringToMultiaddr(s string) (multiaddr, error) {

	// consume trailing slashes
	s = strings.TrimRight(s, "/")

	var b bytes.Buffer
	sp := strings.Split(s, "/")

	if sp[0] != "" {
		return "", fmt.Errorf("invalid multiaddr, must begin with /")
	}

	// consume first empty elem
	sp = sp[1:]

	for len(sp) > 0 {
		p := ProtocolWithName(sp[0])
		if p.Code == 0 {
			return "", fmt.Errorf("no protocol with name %s", sp[0])
		}
		b.Write(p.VCode)
		sp = sp[1:]

		if p.Size == 0 { // no length.
			continue
		}

		if len(sp) < 1 {
			return "", fmt.Errorf("protocol requires address, none given: %s", p.Name)
		}

		if p.Path {
			// it's a path protocol (terminal).
			// consume the rest of the address as the next component.
			sp = []string{"/" + strings.Join(sp, "/")}
		}

		if p.Transcoder == nil {
			return "", fmt.Errorf("no transcoder for %s protocol", p.Name)
		}
		a, err := p.Transcoder.StringToBytes(sp[0])
		if err != nil {
			return "", fmt.Errorf("failed to parse %s: %s %s", p.Name, sp[0], err)
		}
		b.Write(a)
		sp = sp[1:]
	}

	return multiaddr(b.Bytes()), nil
}

func validateBytes(b []byte) (err error) {
	for len(b) > 0 {
		code, n, err := ReadVarintCode(b)
		if err != nil {
			return err
		}

		b = b[n:]
		p := ProtocolWithCode(code)
		if p.Code == 0 {
			return fmt.Errorf("no protocol with code %d", code)
		}

		if p.Size == 0 {
			continue
		}

		size, err := sizeForAddr(p, b)
		if err != nil {
			return err
		}

		if len(b) < size || size < 0 {
			return fmt.Errorf("invalid value for size")
		}

		b = b[size:]
	}

	return nil
}

func bytesToString(b []byte) (ret string, err error) {
	s := ""

	for len(b) > 0 {
		code, n, err := ReadVarintCode(b)
		if err != nil {
			return "", err
		}

		b = b[n:]
		p := ProtocolWithCode(code)
		if p.Code == 0 {
			return "", fmt.Errorf("no protocol with code %d", code)
		}
		s += "/" + p.Name

		if p.Size == 0 {
			continue
		}

		size, err := sizeForAddr(p, b)
		if err != nil {
			return "", err
		}

		if len(b) < size || size < 0 {
			return "", fmt.Errorf("invalid value for size")
		}

		if p.Transcoder == nil {
			return "", fmt.Errorf("no transcoder for %s protocol", p.Name)
		}
		a, err := p.Transcoder.BytesToString(b[:size])
		if err != nil {
			return "", err
		}
		if len(a) > 0 {
			s += "/" + a
		}
		b = b[size:]
	}

	return s, nil
}

func sizeForAddr(p Protocol, b []byte) (int, error) {
	switch {
	case p.Size > 0:
		return (p.Size / 8), nil
	case p.Size == 0:
		return 0, nil
	case p.Path:
		size, n, err := ReadVarintCode(b)
		if err != nil {
			return 0, err
		}
		return size + n, nil
	default:
		size, n, err := ReadVarintCode(b)
		if err != nil {
			return 0, err
		}
		return size + n, nil
	}
}
