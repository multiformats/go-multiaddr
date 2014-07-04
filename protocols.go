package multiaddr

type Protocol struct {
	Code int
	Size int
	Name string
}

// replicating table here to:
// 1. avoid parsing the csv
// 2. ensuring errors in the csv don't screw up code.
// 3. changing a number has to happen in two places.

var Protocols = []*Protocol{
	&Protocol{4, 32, "ip4"},
	&Protocol{6, 16, "tcp"},
	&Protocol{17, 16, "udp"},
	&Protocol{33, 16, "dccp"},
	&Protocol{41, 128, "ip6"},
	// these require varint:
	&Protocol{132, 16, "sctp"},
	// {480, 0, "http"},
	// {443, 0, "https"},
}

func ProtocolWithName(s string) *Protocol {
	for _, p := range Protocols {
		if p.Name == s {
			return p
		}
	}
	return nil
}

func ProtocolWithCode(c int) *Protocol {
	for _, p := range Protocols {
		if p.Code == c {
			return p
		}
	}
	return nil
}
