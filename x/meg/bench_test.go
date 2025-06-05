package meg_test

import (
	"testing"

	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr/x/meg"
)

type preallocatedCapture struct {
	certHashes []string
	matcher    meg.Matcher
}

func preallocateCapture() *preallocatedCapture {
	p := &preallocatedCapture{}
	p.matcher = meg.PatternToMatcher(
		meg.Or(
			meg.Val(multiaddr.P_IP4),
			meg.Val(multiaddr.P_IP6),
			meg.Val(multiaddr.P_DNS),
		),
		meg.Val(multiaddr.P_UDP),
		meg.Val(multiaddr.P_WEBRTC_DIRECT),
		meg.CaptureZeroOrMoreStrings(multiaddr.P_CERTHASH, &p.certHashes),
	)
	return p
}

var webrtcMatchPrealloc *preallocatedCapture

type componentList []multiaddr.Component

func (m componentList) Get(i int) meg.Matchable {
	return &m[i]
}

func (m componentList) Len() int {
	return len(m)
}

func (p *preallocatedCapture) IsWebRTCDirectMultiaddr(addr multiaddr.Multiaddr) (bool, int) {
	found, _ := meg.Match(p.matcher, componentList(addr))
	return found, len(p.certHashes)
}

// IsWebRTCDirectMultiaddr returns whether addr is a /webrtc-direct multiaddr with the count of certhashes
// in addr
func IsWebRTCDirectMultiaddr(addr multiaddr.Multiaddr) (bool, int) {
	if webrtcMatchPrealloc == nil {
		webrtcMatchPrealloc = preallocateCapture()
	}
	return webrtcMatchPrealloc.IsWebRTCDirectMultiaddr(addr)
}

// IsWebRTCDirectMultiaddrLoop returns whether addr is a /webrtc-direct multiaddr with the count of certhashes
// in addr
func IsWebRTCDirectMultiaddrLoop(addr multiaddr.Multiaddr) (bool, int) {
	protos := [...]int{multiaddr.P_IP4, multiaddr.P_IP6, multiaddr.P_DNS, multiaddr.P_UDP, multiaddr.P_WEBRTC_DIRECT}
	matchProtos := [...][]int{protos[:3], {protos[3]}, {protos[4]}}
	certHashCount := 0
	for i, c := range addr {
		if i >= len(matchProtos) {
			if c.Code() == multiaddr.P_CERTHASH {
				certHashCount++
			} else {
				return false, 0
			}
		} else {
			found := false
			for _, proto := range matchProtos[i] {
				if c.Code() == proto {
					found = true
					break
				}
			}
			if !found {
				return false, 0
			}
		}
	}
	return true, certHashCount
}

var wtPrealloc *preallocatedCapture

func isWebTransportMultiaddrPrealloc() *preallocatedCapture {
	if wtPrealloc != nil {
		return wtPrealloc
	}

	p := &preallocatedCapture{}
	var dnsName string
	var ip4Addr string
	var ip6Addr string
	var udpPort string
	var sni string
	p.matcher = meg.PatternToMatcher(
		meg.Or(
			meg.CaptureString(multiaddr.P_IP4, &ip4Addr),
			meg.CaptureString(multiaddr.P_IP6, &ip6Addr),
			meg.CaptureString(multiaddr.P_DNS4, &dnsName),
			meg.CaptureString(multiaddr.P_DNS6, &dnsName),
			meg.CaptureString(multiaddr.P_DNS, &dnsName),
		),
		meg.CaptureString(multiaddr.P_UDP, &udpPort),
		meg.Val(multiaddr.P_QUIC_V1),
		meg.Optional(
			meg.CaptureString(multiaddr.P_SNI, &sni),
		),
		meg.Val(multiaddr.P_WEBTRANSPORT),
		meg.CaptureZeroOrMoreStrings(multiaddr.P_CERTHASH, &p.certHashes),
	)
	wtPrealloc = p
	return p
}

func IsWebTransportMultiaddrPrealloc(m multiaddr.Multiaddr) (bool, int) {
	p := isWebTransportMultiaddrPrealloc()
	found, _ := meg.Match(p.matcher, componentList(m))
	return found, len(p.certHashes)
}

func IsWebTransportMultiaddr(m multiaddr.Multiaddr) (bool, int) {
	var dnsName string
	var ip4Addr string
	var ip6Addr string
	var udpPort string
	var sni string
	var certHashesStr []string
	matched, _ := m.Match(
		meg.Or(
			meg.CaptureString(multiaddr.P_IP4, &ip4Addr),
			meg.CaptureString(multiaddr.P_IP6, &ip6Addr),
			meg.CaptureString(multiaddr.P_DNS4, &dnsName),
			meg.CaptureString(multiaddr.P_DNS6, &dnsName),
			meg.CaptureString(multiaddr.P_DNS, &dnsName),
		),
		meg.CaptureString(multiaddr.P_UDP, &udpPort),
		meg.Val(multiaddr.P_QUIC_V1),
		meg.Optional(
			meg.CaptureString(multiaddr.P_SNI, &sni),
		),
		meg.Val(multiaddr.P_WEBTRANSPORT),
		meg.CaptureZeroOrMoreStrings(multiaddr.P_CERTHASH, &certHashesStr),
	)
	if !matched {
		return false, 0
	}
	return true, len(certHashesStr)
}

func IsWebTransportMultiaddrCaptureBytes(m multiaddr.Multiaddr) (bool, int) {
	var dnsName []byte
	var ip4Addr []byte
	var ip6Addr []byte
	var udpPort []byte
	var sni []byte
	var certHashes [][]byte
	matched, _ := m.Match(
		meg.Or(
			meg.CaptureBytes(multiaddr.P_IP4, &ip4Addr),
			meg.CaptureBytes(multiaddr.P_IP6, &ip6Addr),
			meg.CaptureBytes(multiaddr.P_DNS4, &dnsName),
			meg.CaptureBytes(multiaddr.P_DNS6, &dnsName),
			meg.CaptureBytes(multiaddr.P_DNS, &dnsName),
		),
		meg.CaptureBytes(multiaddr.P_UDP, &udpPort),
		meg.Val(multiaddr.P_QUIC_V1),
		meg.Optional(
			meg.CaptureBytes(multiaddr.P_SNI, &sni),
		),
		meg.Val(multiaddr.P_WEBTRANSPORT),
		meg.CaptureZeroOrMoreBytes(multiaddr.P_CERTHASH, &certHashes),
	)
	if !matched {
		return false, 0
	}
	return true, len(certHashes)
}

func IsWebTransportMultiaddrNoCapture(m multiaddr.Multiaddr) (bool, int) {
	matched, _ := m.Match(
		meg.Or(
			meg.Val(multiaddr.P_IP4),
			meg.Val(multiaddr.P_IP6),
			meg.Val(multiaddr.P_DNS4),
			meg.Val(multiaddr.P_DNS6),
			meg.Val(multiaddr.P_DNS),
		),
		meg.Val(multiaddr.P_UDP),
		meg.Val(multiaddr.P_QUIC_V1),
		meg.Optional(
			meg.Val(multiaddr.P_SNI),
		),
		meg.Val(multiaddr.P_WEBTRANSPORT),
		meg.ZeroOrMore(multiaddr.P_CERTHASH),
	)
	if !matched {
		return false, 0
	}
	return true, 0
}

func IsWebTransportMultiaddrLoop(m multiaddr.Multiaddr) (bool, int) {
	var ip4Addr string
	var ip6Addr string
	var dnsName string
	var udpPort string
	var sni string

	// Expected pattern:
	// 0: one of: P_IP4, P_IP6, P_DNS4, P_DNS6, P_DNS
	// 1: P_UDP
	// 2: P_QUIC_V1
	// 3: optional P_SNI (if present)
	// Next: P_WEBTRANSPORT
	// Trailing: zero or more P_CERTHASH

	// Check minimum length (at least without SNI: 4 components)
	if len(m) < 4 {
		return false, 0
	}

	idx := 0

	// Component 0: Must be one of IP or DNS protocols.
	switch m[idx].Code() {
	case multiaddr.P_IP4:
		ip4Addr = m[idx].String()
	case multiaddr.P_IP6:
		ip6Addr = m[idx].String()
	case multiaddr.P_DNS4, multiaddr.P_DNS6, multiaddr.P_DNS:
		dnsName = m[idx].String()
	default:
		return false, 0
	}
	idx++

	// Component 1: Must be UDP.
	if idx >= len(m) || m[idx].Code() != multiaddr.P_UDP {
		return false, 0
	}
	udpPort = m[idx].String()
	idx++

	// Component 2: Must be QUIC_V1.
	if idx >= len(m) || m[idx].Code() != multiaddr.P_QUIC_V1 {
		return false, 0
	}
	idx++

	// Optional component: SNI.
	if idx < len(m) && m[idx].Code() == multiaddr.P_SNI {
		sni = m[idx].String()
		idx++
	}

	// Next component: Must be WEBTRANSPORT.
	if idx >= len(m) || m[idx].Code() != multiaddr.P_WEBTRANSPORT {
		return false, 0
	}
	idx++

	// All remaining components must be CERTHASH.
	certHashCount := 0
	for ; idx < len(m); idx++ {
		if m[idx].Code() != multiaddr.P_CERTHASH {
			return false, 0
		}
		_ = m[idx].String()
		certHashCount++
	}

	_ = ip4Addr
	_ = ip6Addr
	_ = dnsName
	_ = udpPort
	_ = sni

	return true, certHashCount
}

func IsWebTransportMultiaddrLoopNoCapture(m multiaddr.Multiaddr) (bool, int) {
	// Expected pattern:
	// 0: one of: P_IP4, P_IP6, P_DNS4, P_DNS6, P_DNS
	// 1: P_UDP
	// 2: P_QUIC_V1
	// 3: optional P_SNI (if present)
	// Next: P_WEBTRANSPORT
	// Trailing: zero or more P_CERTHASH

	// Check minimum length (at least without SNI: 4 components)
	if len(m) < 4 {
		return false, 0
	}

	idx := 0

	// Component 0: Must be one of IP or DNS protocols.
	switch m[idx].Code() {
	case multiaddr.P_IP4:
	case multiaddr.P_IP6:
	case multiaddr.P_DNS4, multiaddr.P_DNS6, multiaddr.P_DNS:
	default:
		return false, 0
	}
	idx++

	// Component 1: Must be UDP.
	if idx >= len(m) || m[idx].Code() != multiaddr.P_UDP {
		return false, 0
	}
	idx++

	// Component 2: Must be QUIC_V1.
	if idx >= len(m) || m[idx].Code() != multiaddr.P_QUIC_V1 {
		return false, 0
	}
	idx++

	// Optional component: SNI.
	if idx < len(m) && m[idx].Code() == multiaddr.P_SNI {
		idx++
	}

	// Next component: Must be WEBTRANSPORT.
	if idx >= len(m) || m[idx].Code() != multiaddr.P_WEBTRANSPORT {
		return false, 0
	}
	idx++

	// All remaining components must be CERTHASH.
	for ; idx < len(m); idx++ {
		if m[idx].Code() != multiaddr.P_CERTHASH {
			return false, 0
		}
		_ = m[idx].String()
	}

	return true, 0
}

func BenchmarkIsWebTransportMultiaddrPrealloc(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/quic-v1/sni/example.com/webtransport")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWT, count := IsWebTransportMultiaddrPrealloc(addr)
		if !isWT || count != 0 {
			b.Fatal("unexpected result")
		}
	}
}

func BenchmarkIsWebTransportMultiaddrNoCapturePrealloc(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/quic-v1/sni/example.com/webtransport")

	wtPreallocNoCapture := meg.PatternToMatcher(
		meg.Or(
			meg.Val(multiaddr.P_IP4),
			meg.Val(multiaddr.P_IP6),
			meg.Val(multiaddr.P_DNS4),
			meg.Val(multiaddr.P_DNS6),
			meg.Val(multiaddr.P_DNS),
		),
		meg.Val(multiaddr.P_UDP),
		meg.Val(multiaddr.P_QUIC_V1),
		meg.Optional(
			meg.Val(multiaddr.P_SNI),
		),
		meg.Val(multiaddr.P_WEBTRANSPORT),
		meg.ZeroOrMore(multiaddr.P_CERTHASH),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWT, _ := meg.Match(wtPreallocNoCapture, componentList(addr))
		if !isWT {
			b.Fatal("unexpected result")
		}
	}
}

func BenchmarkIsWebTransportMultiaddrNoCapture(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/quic-v1/sni/example.com/webtransport")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWT, count := IsWebTransportMultiaddrNoCapture(addr)
		if !isWT || count != 0 {
			b.Fatal("unexpected result")
		}
	}
}

func BenchmarkIsWebTransportMultiaddrCaptureBytes(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/quic-v1/sni/example.com/webtransport")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWT, count := IsWebTransportMultiaddrCaptureBytes(addr)
		if !isWT || count != 0 {
			b.Fatal("unexpected result")
		}
	}
}

func BenchmarkIsWebTransportMultiaddr(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/quic-v1/sni/example.com/webtransport")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWT, count := IsWebTransportMultiaddr(addr)
		if !isWT || count != 0 {
			b.Fatal("unexpected result")
		}
	}
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

func BenchmarkIsWebTransportMultiaddrLoopNoCapture(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/quic-v1/sni/example.com/webtransport")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWT, count := IsWebTransportMultiaddrLoopNoCapture(addr)
		if !isWT || count != 0 {
			b.Fatal("unexpected result")
		}
	}
}

func BenchmarkIsWebRTCDirectMultiaddr(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/webrtc-direct/")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWebRTC, count := IsWebRTCDirectMultiaddr(addr)
		if !isWebRTC || count != 0 {
			b.Fatal("unexpected result")
		}
	}
}

func BenchmarkIsWebRTCDirectMultiaddrLoop(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/webrtc-direct/")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWebRTC, count := IsWebRTCDirectMultiaddrLoop(addr)
		if !isWebRTC || count != 0 {
			b.Fatal("unexpected result")
		}
	}
}
