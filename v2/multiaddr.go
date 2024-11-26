package multiaddr

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

type ProtocolName string
type ProtocolCode uint64

type Component struct {
	value          []byte
	protocolCode   ProtocolCode
	isVariableSize bool
}

func (c Component) Code() int {
	return int(c.protocolCode)
}

func (c Component) Value() []byte {
	out := make([]byte, len(c.value))
	copy(out, c.value)
	return out
}

// MultiaddrBytes represents the binary form of a multiaddr
type MultiaddrBytes struct {
	Bytes []byte
}

type Multiaddr []Component

func (m Multiaddr) PopLast() (Multiaddr, Component) {
	if len(m) == 0 {
		return nil, Component{}
	}
	return m[:len(m)-1], m[len(m)-1]
}

func (m Multiaddr) ToBinary() ([]byte, error) {
	var out []byte

	for _, c := range m {
		if c.protocolCode == 0 {
			return nil, errors.New("invalid multiaddr: component has no protocol code")
		}
		out = binary.AppendUvarint(out, uint64(c.protocolCode))
		if c.isVariableSize {
			out = binary.AppendUvarint(out, uint64(len(c.value)))
		}
		out = append(out, c.value...)
	}

	return out, nil
}

// MultiaddrTranscoder coverts between Strings/Bytes to Multiaddrs and back.
// A default MultiaddrTranscoder is provided with all configured protocols.
// Users are encouraged to build their own MultiaddrTranscoder and support only
// the protocols they are working with.
type MultiaddrTranscoder struct {
	protocols       []Protocol
	protocolsByName map[ProtocolName]Protocol
	protocolsByCode map[ProtocolCode]Protocol
}

func (t *MultiaddrTranscoder) AddProtocol(p Protocol) error {
	if t.protocolsByName == nil {
		t.protocolsByName = make(map[ProtocolName]Protocol)
	}
	if t.protocolsByCode == nil {
		t.protocolsByCode = make(map[ProtocolCode]Protocol)
	}

	t.protocols = append(t.protocols, p)
	if _, ok := t.protocolsByName[p.Name]; ok {
		return fmt.Errorf("protocol by the name %q already exists", p.Name)
	}

	if _, ok := t.protocolsByCode[p.Code]; ok {
		return fmt.Errorf("protocol code %d already taken by %q", p.Code, p.Code)
	}

	if p.Size != 0 && p.Transcoder == nil {
		return fmt.Errorf("protocols with arguments must define transcoders")
	}
	if p.Path && p.Size >= 0 {
		return fmt.Errorf("path protocols must have variable-length sizes")
	}

	Protocols = append(Protocols, p)
	t.protocolsByName[p.Name] = p
	t.protocolsByCode[p.Code] = p
	return nil
}

func (t *MultiaddrTranscoder) AliasProtocolName(from ProtocolName, to ProtocolName) error {
	if t.protocolsByName == nil {
		t.protocolsByName = make(map[ProtocolName]Protocol)
	}
	if t.protocolsByCode == nil {
		t.protocolsByCode = make(map[ProtocolCode]Protocol)
	}

	p, ok := t.protocolsByName[to]
	if !ok {
		return fmt.Errorf("protocol %q is missing", to)
	}

	t.protocolsByName[from] = p
	return nil
}

func (t *MultiaddrTranscoder) ProtocolWithName(s ProtocolName) (Protocol, bool) {
	v, ok := t.protocolsByName[s]
	return v, ok
}

func (t *MultiaddrTranscoder) ProtocolWithCode(c ProtocolCode) (Protocol, bool) {
	v, ok := t.protocolsByCode[c]
	return v, ok
}

func (t *MultiaddrTranscoder) Decode(mab MultiaddrBytes) (Multiaddr, error) {
	return t.FromBinary(mab.Bytes)
}

func (t *MultiaddrTranscoder) FromBinary(b []byte) (Multiaddr, error) {
	var out []Component
	for len(b) > 0 {
		code, n := binary.Uvarint(b)
		if n <= 0 {
			return nil, errors.New("invalid multiaddr: invalid protocol code")
		}
		b = b[n:]
		p, ok := t.protocolsByCode[ProtocolCode(code)]
		if !ok {
			return nil, fmt.Errorf("unsupported protocol code: %d", code)
		}

		isVariableSize := p.Size < 0

		var valSize int
		if isVariableSize {
			valSizeu64, n := binary.Uvarint(b)
			if n <= 0 {
				return nil, errors.New("invalid multiaddr: invalid protocol value size")
			}
			b = b[n:]
			valSize = int(valSizeu64)
		} else {
			valSize = p.Size
		}

		if int(valSize) > len(b) {
			return nil, errors.New("invalid multiaddr: invalid protocol value size (too large)")
		}
		value := make([]byte, valSize)
		copy(value, b[:valSize])

		c := Component{
			value:          value,
			protocolCode:   ProtocolCode(code),
			isVariableSize: isVariableSize,
		}
		out = append(out, c)
	}

	if len(out) == 0 {
		return nil, errors.New("no components found")
	}

	return out, nil
}

func (t *MultiaddrTranscoder) ToBinary(m Multiaddr) ([]byte, error) {
	return m.ToBinary()
}

func (t *MultiaddrTranscoder) FromString(s string) (Multiaddr, error) {
	if len(s) == 0 {
		return nil, errors.New("invalid multiaddr: empty string")
	}
	if s[0] != '/' {
		return nil, errors.New("invalid multiaddr: must start with '/'")
	}
	s = s[1:] // Start after the first slash

	var out []Component

	scanner := bufio.NewScanner(strings.NewReader(s))
	// Set the split function for the scanning operation.
	scanner.Split(splitTextOnSlash)

	for scanner.Scan() {
		protoName := scanner.Text()
		p, ok := t.protocolsByName[ProtocolName(protoName)]
		if !ok {
			return nil, fmt.Errorf("unsupported protocol: %s", protoName)
		}
		if p.Size == 0 {
			// No value. We exit this iteration early
			out = append(out, Component{protocolCode: ProtocolCode(p.Code)})
			continue
		}

		if !scanner.Scan() {
			return nil, errors.New("missing value for protocol")
		}
		valStr := scanner.Text()
		val, err := p.Transcoder.StringToBytes(valStr)
		if err != nil {
			return nil, err
		}

		out = append(out, Component{
			value:          val,
			protocolCode:   ProtocolCode(p.Code),
			isVariableSize: p.Size < 0,
		})
	}

	if len(out) == 0 {
		return nil, errors.New("no components found")
	}

	return out, nil
}

func (t *MultiaddrTranscoder) ToString(m Multiaddr) (string, error) {
	if len(t.protocolsByCode) == 0 {
		return "", errors.New("no supported protocols")
	}

	var out strings.Builder

	for _, c := range m {
		if c.protocolCode == 0 {
			return "", errors.New("invalid multiaddr: component has no protocol code")
		}

		p, ok := t.protocolsByCode[c.protocolCode]
		if !ok {
			return "", fmt.Errorf("unsupported protocol code: %d", c.protocolCode)
		}

		out.WriteRune('/')
		out.WriteString(string(p.Name))

		if len(c.value) != 0 {
			err := p.Transcoder.ValidateBytes(c.value)
			if err != nil {
				return "", err
			}

			s, err := p.Transcoder.BytesToString(c.value)
			if err != nil {
				return "", err
			}
			out.WriteRune('/')
			out.WriteString(s)
		}
	}

	return out.String(), nil
}

// var _ = Transcoder(MultiaddrTranscoder{})

// splitTextOnSlash splits the input data on the '/' character.
// It returns the advance count, the token, and an error if any.
// If atEOF is true and there is remaining data, it returns the final token.
//
// For use with bufio.Scanner.
func splitTextOnSlash(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0

	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '/' {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}
