package manet

import (
	"net"
	"strings"

	ma "github.com/multiformats/go-multiaddr"
)

// Private4 and Private6 are well-known private networks
var Private4, Private6 []*net.IPNet
var privateCIDR4 = []string{
	// localhost
	"127.0.0.0/8",
	// private networks
	"10.0.0.0/8",
	"100.64.0.0/10",
	"172.16.0.0/12",
	"192.168.0.0/16",
	// link local
	"169.254.0.0/16",
}
var privateCIDR6 = []string{
	// localhost
	"::1/128",
	// ULA reserved
	"fc00::/7",
	// link local
	"fe80::/10",
}

// Unroutable4 and Unroutable6 are well known unroutable address ranges
var Unroutable4, Unroutable6 []*net.IPNet
var unroutableCIDR4 = []string{
	"0.0.0.0/8",
	"192.0.0.0/26",
	"192.0.2.0/24",
	"192.88.99.0/24",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"255.255.255.255/32",
}
var unroutableCIDR6 = []string{
	"ff00::/8",
}

// specialUseDomains are reserved for various purposes and do not have a central authority
// for consistent resolution in different networks.
// see: https://en.wikipedia.org/wiki/Special-use_domain_name#Reserved_domain_names
// This list doesn't contain `.onion` addresses as they are consistently resolved everywhere.
var specialUseDomains = []string{
	"6tisch.arpa",
	"10.in-addr.arpa",
	"16.172.in-addr.arpa",
	"17.172.in-addr.arpa",
	"18.172.in-addr.arpa",
	"19.172.in-addr.arpa",
	"20.172.in-addr.arpa",
	"21.172.in-addr.arpa",
	"22.172.in-addr.arpa",
	"23.172.in-addr.arpa",
	"24.172.in-addr.arpa",
	"25.172.in-addr.arpa",
	"26.172.in-addr.arpa",
	"27.172.in-addr.arpa",
	"28.172.in-addr.arpa",
	"29.172.in-addr.arpa",
	"30.172.in-addr.arpa",
	"31.172.in-addr.arpa",
	"168.192.in-addr.arpa",
	"170.0.0.192.in-addr.arpa",
	"171.0.0.192.in-addr.arpa",
	"ipv4only.arpa",
	"254.169.in-addr.arpa",
	"8.e.f.ip6.arpa",
	"9.e.f.ip6.arpa",
	"a.e.f.ip6.arpa",
	"b.e.f.ip6.arpa",
	"home.arpa",
	"example",
	"example.com",
	"example.net",
	"example.org",
	"invalid",
	"intranet",
	"internal",
	"private",
	"corp",
	"home",
	"lan",
	"local",
	"localhost",
	"test",
}

func init() {
	Private4 = parseCIDR(privateCIDR4)
	Private6 = parseCIDR(privateCIDR6)
	Unroutable4 = parseCIDR(unroutableCIDR4)
	Unroutable6 = parseCIDR(unroutableCIDR6)
}

func parseCIDR(cidrs []string) []*net.IPNet {
	ipnets := make([]*net.IPNet, len(cidrs))
	for i, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		ipnets[i] = ipnet
	}
	return ipnets
}

// IsPublicAddr returns true if the IP part of the multiaddr is a publicly routable address
// or if it's a dns address without a special use domain e.g. .local.
func IsPublicAddr(a ma.Multiaddr) bool {
	isPublic := false
	ma.ForEach(a, func(c ma.Component) bool {
		switch c.Protocol().Code {
		case ma.P_IP6ZONE:
			return true
		case ma.P_IP4:
			ip := net.IP(c.RawValue())
			isPublic = !inAddrRange(ip, Private4) && !inAddrRange(ip, Unroutable4)
		case ma.P_IP6:
			ip := net.IP(c.RawValue())
			isPublic = !inAddrRange(ip, Private6) && !inAddrRange(ip, Unroutable6)
		case ma.P_DNS, ma.P_DNS4, ma.P_DNS6, ma.P_DNSADDR:
			dnsAddr := c.Value()
			isPublic = true
			for _, sd := range specialUseDomains {
				if strings.HasSuffix(dnsAddr, sd) {
					isPublic = false
					break
				}
			}
		}
		return false
	})
	return isPublic
}

// IsPrivateAddr returns true if the IP part of the mutiaddr is in a private network
func IsPrivateAddr(a ma.Multiaddr) bool {
	isPrivate := false
	ma.ForEach(a, func(c ma.Component) bool {
		switch c.Protocol().Code {
		case ma.P_IP6ZONE:
			return true
		case ma.P_IP4:
			isPrivate = inAddrRange(net.IP(c.RawValue()), Private4)
		case ma.P_IP6:
			isPrivate = inAddrRange(net.IP(c.RawValue()), Private6)
		}
		return false
	})
	return isPrivate
}

func inAddrRange(ip net.IP, ipnets []*net.IPNet) bool {
	for _, ipnet := range ipnets {
		if ipnet.Contains(ip) {
			return true
		}
	}

	return false
}
