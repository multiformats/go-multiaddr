package multiaddrv3

type ProtocolCode uint32

// Component is the component of a multiaddr, e.g. /udp/1234
// It consumes exactly 64 bytes.
// It is comparable using the == operator.
type Component struct {
	proto ProtocolCode
	val   string
}

func (c Component) Protocol() Protocol {
	return ProtocolWithCode(int(c.proto))
}

func (c Component) Value() []byte {
	return []byte(c.val)
}

func (c Component) String() string {
	s, _ := c.Protocol().Transcoder.BytesToString(c.Value())
	return s
}
