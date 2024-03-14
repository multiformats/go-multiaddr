package manet

import (
	"fmt"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

func TestIsWellKnownPrefixIPv4ConvertedIPv6Address(t *testing.T) {
	cases := []struct {
		addr          ma.Multiaddr
		want          bool
		failureReason string
	}{
		{
			addr:          ma.StringCast("/ip4/1.2.3.4/tcp/1234"),
			want:          false,
			failureReason: "ip4 addresses should return false",
		},
		{
			addr:          ma.StringCast("/ip6/1::4/tcp/1234"),
			want:          false,
			failureReason: "ip6 addresses doesn't have well-known prefix",
		},
		{
			addr:          ma.StringCast("/ip6/::1/tcp/1234"),
			want:          false,
			failureReason: "localhost addresses should return false",
		},
		{
			addr:          ma.StringCast("/ip6/64:ff9b::192.0.1.2/tcp/1234"),
			want:          true,
			failureReason: "ip6 address begins with well-known prefix",
		},
		{
			addr:          ma.StringCast("/ip6/64:ff9b::1:192.0.1.2/tcp/1234"),
			want:          false,
			failureReason: "64:ff9b::1 is not well-known prefix",
		},
		{
			addr:          ma.StringCast("/ip6/64:ff9b:1::1:192.0.1.2/tcp/1234"),
			want:          true,
			failureReason: "64:ff9b:1::1 is allowed for NAT64 translation",
		},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if IsNAT64IPv4ConvertedIPv6Addr(tc.addr) != tc.want {
				t.Fatalf("%s %s", tc.addr, tc.failureReason)
			}
		})
	}
}
