package multiaddr_test

import (
	"testing"

	"github.com/multiformats/go-multiaddr"
)

func IsWebTransportMultiaddrLoop(addr multiaddr.Multiaddr) (bool, int) {
	var ip4Addr string
	var ip6Addr string
	var dnsName string
	var udpPort string
	var sni string

	components := []multiaddr.Component{}

	multiaddr.ForEach(addr, func(c multiaddr.Component) bool {
		components = append(components, c)
		return true
	})

	// Expected pattern:
	// 0: one of: P_IP4, P_IP6, P_DNS4, P_DNS6, P_DNS
	// 1: P_UDP
	// 2: P_QUIC_V1
	// 3: optional P_SNI (if present)
	// Next: P_WEBTRANSPORT
	// Trailing: zero or more P_CERTHASH

	// Check minimum length (at least without SNI: 4 components)
	if len(components) < 4 {
		return false, 0
	}

	idx := 0

	// Component 0: Must be one of IP or DNS protocols.
	switch components[idx].Protocol().Code {
	case multiaddr.P_IP4:
		ip4Addr = components[idx].String()
	case multiaddr.P_IP6:
		ip6Addr = components[idx].String()
	case multiaddr.P_DNS4, multiaddr.P_DNS6, multiaddr.P_DNS:
		dnsName = components[idx].String()
	default:
		return false, 0
	}
	idx++

	// Component 1: Must be UDP.
	if idx >= len(components) || components[idx].Protocol().Code != multiaddr.P_UDP {
		return false, 0
	}
	udpPort = components[idx].String()
	idx++

	// Component 2: Must be QUIC_V1.
	if idx >= len(components) || components[idx].Protocol().Code != multiaddr.P_QUIC_V1 {
		return false, 0
	}
	idx++

	// Optional component: SNI.
	if idx < len(components) && components[idx].Protocol().Code == multiaddr.P_SNI {
		sni = components[idx].String()
		idx++
	}

	// Next component: Must be WEBTRANSPORT.
	if idx >= len(components) || components[idx].Protocol().Code != multiaddr.P_WEBTRANSPORT {
		return false, 0
	}
	idx++

	// All remaining components must be CERTHASH.
	certHashCount := 0
	for ; idx < len(components); idx++ {
		if components[idx].Protocol().Code != multiaddr.P_CERTHASH {
			return false, 0
		}
		_ = components[idx].String()
		certHashCount++
	}

	_ = ip4Addr
	_ = ip6Addr
	_ = dnsName
	_ = udpPort
	_ = sni

	return true, certHashCount
}

func BenchmarkIsWebTransportMultiaddrLoop(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/quic-v1/sni/example.com/webtransport")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWT, count := IsWebTransportMultiaddrLoop(addr)
		if !isWT || count != 0 {
			b.Fatal("unexpected result")
		}
	}
}
