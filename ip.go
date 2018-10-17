package manet

import (
	"net"

	ma "github.com/multiformats/go-multiaddr"
)

// Loopback Addresses
var (
	// IP4Loopback is the ip4 loopback multiaddr
	IP4Loopback = ma.StringCast("/ip4/127.0.0.1")

	// IP6Loopback is the ip6 loopback multiaddr
	IP6Loopback = ma.StringCast("/ip6/::1")

	// IP4MappedIP6Loopback is the IPv4 Mapped IPv6 loopback address.
	IP4MappedIP6Loopback = ma.StringCast("/ip6/::ffff:127.0.0.1")
)

// Unspecified Addresses (used for )
var (
	IP4Unspecified = ma.StringCast("/ip4/0.0.0.0")
	IP6Unspecified = ma.StringCast("/ip6/::")
)

// IsThinWaist returns whether a Multiaddr starts with "Thin Waist" Protocols.
// This means: /{IP4, IP6}[/{TCP, UDP}]
func IsThinWaist(m ma.Multiaddr) bool {
	p := m.Protocols()

	// nothing? not even a waist.
	if len(p) == 0 {
		return false
	}

	if p[0].Code != ma.P_IP4 && p[0].Code != ma.P_IP6 {
		return false
	}

	// only IP? still counts.
	if len(p) == 1 {
		return true
	}

	switch p[1].Code {
	case ma.P_TCP, ma.P_UDP, ma.P_IP4, ma.P_IP6:
		return true
	default:
		return false
	}
}

// IsIPLoopback returns whether a Multiaddr is a "Loopback" IP address
// This means either /ip4/127.*.*.*, /ip6/::1, or /ip6/::ffff:127.*.*.*.*
func IsIPLoopback(m ma.Multiaddr) bool {
	c, rest := ma.SplitFirst(m)
	if rest != nil {
		// Not *just* an IPv4 addr
		return false
	}
	switch c.Protocol().Code {
	case ma.P_IP4, ma.P_IP6:
		return net.IP(c.RawValue()).IsLoopback()
	}
	return false
}

// IsIP6LinkLocal returns if a an IPv6 link-local multiaddress (with zero or
// more leading zones). These addresses are non routable.
func IsIP6LinkLocal(m ma.Multiaddr) bool {
	matched := false
	ma.ForEach(m, func(c ma.Component) bool {
		// Too much.
		if matched {
			matched = false
			return false
		}

		switch c.Protocol().Code {
		case ma.P_IP6ZONE:
			return true
		case ma.P_IP6:
			ip := net.IP(c.RawValue())
			matched = ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast()
			return true
		default:
			return false
		}
	})
	return matched
}

// IsIPUnspecified returns whether a Multiaddr is am Unspecified IP address
// This means either /ip4/0.0.0.0 or /ip6/::
func IsIPUnspecified(m ma.Multiaddr) bool {
	return IP4Unspecified.Equal(m) || IP6Unspecified.Equal(m)
}
