package manet

import (
	"fmt"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

func TestIsPublicAddr(t *testing.T) {
	tests := []struct {
		addr      ma.Multiaddr
		isPublic  bool
		isPrivate bool
	}{
		{
			addr:      ma.StringCast("/ip4/192.168.1.1/tcp/80"),
			isPublic:  false,
			isPrivate: true,
		},
		{
			addr:      ma.StringCast("/ip4/1.1.1.1/tcp/80"),
			isPublic:  true,
			isPrivate: false,
		},
		{
			addr:      ma.StringCast("/tcp/80/ip4/1.1.1.1"),
			isPublic:  false,
			isPrivate: false,
		},
		{
			addr:      ma.StringCast("/dns/node.libp2p.io/udp/1/quic-v1"),
			isPublic:  true,
			isPrivate: false,
		},
		{
			addr:      ma.StringCast("/dnsaddr/node.libp2p.io/udp/1/quic-v1"),
			isPublic:  true,
			isPrivate: false,
		},
		{
			addr:      ma.StringCast("/dns/node.libp2p.local/udp/1/quic-v1"),
			isPublic:  false,
			isPrivate: false, // You can configure .local domains in local networks to return public addrs
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			isPublic := IsPublicAddr(tt.addr)
			isPrivate := IsPrivateAddr(tt.addr)
			if isPublic != tt.isPublic {
				t.Errorf("IsPublicAddr check failed for %s: expected %t, got %t", tt.addr, tt.isPublic, isPublic)
			}
			if isPrivate != tt.isPrivate {
				t.Errorf("IsPrivateAddr check failed for %s: expected %t, got %t", tt.addr, tt.isPrivate, isPrivate)
			}
		})
	}
}
