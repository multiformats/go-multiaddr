package multiaddr

import (
	"net/netip"
	"testing"

	"github.com/multiformats/go-multiaddr/x/meg"
)

func TestMatchAndCaptureMultiaddr(t *testing.T) {
	m := StringCast("/ip4/1.2.3.4/udp/8231/quic-v1/webtransport/certhash/b2uaraocy6yrdblb4sfptaddgimjmmpy/certhash/zQmbWTwYGcmdyK9CYfNBcfs9nhZs17a6FQ4Y8oea278xx41")

	var udpPort string
	var certhashes []string
	found, _ := m.Match(
		meg.Or(
			meg.Val(P_IP4),
			meg.Val(P_IP6),
		),
		meg.CaptureString(P_UDP, &udpPort),
		meg.Val(P_QUIC_V1),
		meg.Val(P_WEBTRANSPORT),
		meg.CaptureZeroOrMoreStrings(P_CERTHASH, &certhashes),
	)
	if !found {
		t.Fatal("failed to match")
	}
	if udpPort != "8231" {
		t.Fatal("unexpected value")
	}

	if len(certhashes) != 2 {
		t.Fatal("Didn't capture all certhashes")
	}

	{
		m, c := SplitLast(m)
		if c.Value() != certhashes[1] {
			t.Fatal("unexpected value. Expected", c.RawValue(), "but got", []byte(certhashes[1]))
		}
		_, c = SplitLast(m)
		if c.Value() != certhashes[0] {
			t.Fatal("unexpected value. Expected", c.RawValue(), "but got", []byte(certhashes[0]))
		}
	}
}

func TestCaptureAddrPort(t *testing.T) {
	m := StringCast("/ip4/1.2.3.4/udp/8231/quic-v1/webtransport")
	var addrPort netip.AddrPort
	var network string

	found, err := m.Match(
		CaptureAddrPort(&network, &addrPort),
		meg.ZeroOrMore(meg.Any),
	)
	if err != nil {
		t.Fatal("error", err)
	}
	if !found {
		t.Fatal("failed to match")
	}
	if !addrPort.IsValid() {
		t.Fatal("failed to capture addrPort")
	}
	if network != "udp" {
		t.Fatal("unexpected network", network)
	}
	if addrPort.String() != "1.2.3.4:8231" {
		t.Fatal("unexpected ipPort", addrPort)
	}
}
