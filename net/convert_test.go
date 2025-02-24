package manet

import (
	"net"
	"runtime"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

type GenFunc func() (ma.Multiaddr, error)

func testConvert(t *testing.T, s string, gen GenFunc) {
	m, err := gen()
	if err != nil {
		t.Fatal("failed to generate.")
	}

	if s2 := m.String(); err != nil || s2 != s {
		t.Fatal("failed to convert: " + s + " != " + s2)
	}
}

func testToNetAddr(t *testing.T, maddr, ntwk, addr string) {
	m, err := ma.NewMultiaddr(maddr)
	if err != nil {
		t.Fatal("failed to generate.")
	}

	naddr, err := ToNetAddr(m)
	if addr == "" { // should fail
		if err == nil {
			t.Fatalf("failed to error: %s", m)
		}
		return
	}

	// shouldn't fail
	if err != nil {
		t.Fatalf("failed to convert to net addr: %s", m)
	}

	if naddr.String() != addr {
		t.Fatalf("naddr.Address() == %s != %s", naddr, addr)
	}

	if naddr.Network() != ntwk {
		t.Fatalf("naddr.Network() == %s != %s", naddr.Network(), ntwk)
	}

	// should convert properly
	switch ntwk {
	case "tcp":
		taddr := naddr.(*net.TCPAddr)
		if ip, err := ToIP(m); err != nil || !taddr.IP.Equal(ip) {
			t.Fatalf("ToIP() and ToNetAddr diverged: %s != %s", taddr, ip)
		}
	case "udp":
		uaddr := naddr.(*net.UDPAddr)
		if ip, err := ToIP(m); err != nil || !uaddr.IP.Equal(ip) {
			t.Fatalf("ToIP() and ToNetAddr diverged: %s != %s", uaddr, ip)
		}
	case "ip":
		ipaddr := naddr.(*net.IPAddr)
		if ip, err := ToIP(m); err != nil || !ipaddr.IP.Equal(ip) {
			t.Fatalf("ToIP() and ToNetAddr diverged: %s != %s", ipaddr, ip)
		}
	}
}

func TestFromIP4(t *testing.T) {
	testConvert(t, "/ip4/10.20.30.40", func() (ma.Multiaddr, error) {
		return FromNetAddr(&net.IPAddr{IP: net.ParseIP("10.20.30.40")})
	})
}

func TestFromUnix(t *testing.T) {
	path := "/C:/foo/bar"
	if runtime.GOOS == "windows" {
		path = `C:\foo\bar`
	}
	testConvert(t, "/unix/C:/foo/bar", func() (ma.Multiaddr, error) {
		return FromNetAddr(&net.UnixAddr{Name: path, Net: "unix"})
	})
}

func TestToUnix(t *testing.T) {
	path := "/C:/foo/bar"
	if runtime.GOOS == "windows" {
		path = `C:\foo\bar`
	}
	testToNetAddr(t, "/unix/C:/foo/bar", "unix", path)
}

func TestFromIP6(t *testing.T) {
	testConvert(t, "/ip6/2001:4860:0:2001::68", func() (ma.Multiaddr, error) {
		return FromNetAddr(&net.IPAddr{IP: net.ParseIP("2001:4860:0:2001::68")})
	})
}

func TestFromTCP(t *testing.T) {
	testConvert(t, "/ip4/10.20.30.40/tcp/1234", func() (ma.Multiaddr, error) {
		return FromNetAddr(&net.TCPAddr{
			IP:   net.ParseIP("10.20.30.40"),
			Port: 1234,
		})
	})
}

func TestFromUDP(t *testing.T) {
	testConvert(t, "/ip4/10.20.30.40/udp/1234", func() (ma.Multiaddr, error) {
		return FromNetAddr(&net.UDPAddr{
			IP:   net.ParseIP("10.20.30.40"),
			Port: 1234,
		})
	})
}

func TestThinWaist(t *testing.T) {
	addrs := map[string]bool{
		"/ip4/127.0.0.1/udp/1234":                true,
		"/ip4/127.0.0.1/tcp/1234":                true,
		"/ip4/127.0.0.1/udp/1234/tcp/1234":       true,
		"/ip4/127.0.0.1/tcp/12345/ip4/1.2.3.4":   true,
		"/ip6/::1/tcp/80":                        true,
		"/ip6/::1/udp/80":                        true,
		"/ip6/::1":                               true,
		"/ip6zone/hello/ip6/fe80::1/tcp/80":      true,
		"/ip6zone/hello/ip6/fe80::1":             true,
		"/tcp/1234/ip4/1.2.3.4":                  false,
		"/tcp/1234":                              false,
		"/tcp/1234/udp/1234":                     false,
		"/ip4/1.2.3.4/ip4/2.3.4.5":               true,
		"/ip6/fe80::1/ip4/2.3.4.5":               true,
		"/ip6zone/hello/ip6/fe80::1/ip4/2.3.4.5": true,

		// Invalid ip6zone usage:
		"/ip6zone/hello":             false,
		"/ip6zone/hello/ip4/1.1.1.1": false,
	}

	for a, res := range addrs {
		m, err := ma.NewMultiaddr(a)
		if err != nil {
			t.Fatalf("failed to construct Multiaddr: %s", a)
		}

		if IsThinWaist(m) != res {
			t.Fatalf("IsThinWaist(%s) != %v", a, res)
		}
	}
}

func TestDialArgs(t *testing.T) {
	test := func(e_maddr, e_nw, e_host string) {
		m, err := ma.NewMultiaddr(e_maddr)
		if err != nil {
			t.Fatal("failed to construct", e_maddr)
		}

		nw, host, err := DialArgs(m)
		if err != nil {
			t.Fatal("failed to get dial args", e_maddr, m, err)
		}

		if nw != e_nw {
			t.Error("failed to get udp network Dial Arg", e_nw, nw)
		}

		if host != e_host {
			t.Error("failed to get host:port Dial Arg", e_host, host)
		}
	}

	test_error := func(e_maddr string) {
		m, err := ma.NewMultiaddr(e_maddr)
		if err != nil {
			t.Fatal("failed to construct", e_maddr)
		}

		_, _, err = DialArgs(m)
		if err == nil {
			t.Fatal("expected DialArgs to fail on", e_maddr)
		}
	}

	test("/ip4/127.0.0.1/udp/1234", "udp4", "127.0.0.1:1234")
	test("/ip4/127.0.0.1/tcp/4321", "tcp4", "127.0.0.1:4321")
	test("/ip6/::1/udp/1234", "udp6", "[::1]:1234")
	test("/ip6/::1/tcp/4321", "tcp6", "[::1]:4321")
	test("/ip6/::1", "ip6", "::1")                                  // Just an IP
	test("/ip4/1.2.3.4", "ip4", "1.2.3.4")                          // Just an IP
	test("/ip6zone/foo/ip6/::1/tcp/4321", "tcp6", "[::1%foo]:4321") // zone
	test("/ip6zone/foo/ip6/::1/udp/4321", "udp6", "[::1%foo]:4321") // zone
	test("/ip6zone/foo/ip6/::1", "ip6", "::1%foo")                  // no TCP
	test_error("/ip6zone/foo/ip4/127.0.0.1")                        // IP4 doesn't take zone
	test("/ip6zone/foo/ip6/::1/ip6zone/bar", "ip6", "::1%foo")      // IP over IP
	test_error("/ip6zone/foo/ip6zone/bar/ip6/::1")                  // Only one zone per IP6
	test("/dns/abc.com/tcp/1234", "tcp", "abc.com:1234")            // DNS4:port
	test("/dns4/abc.com/tcp/1234", "tcp4", "abc.com:1234")          // DNS4:port
	test("/dns4/abc.com", "ip4", "abc.com")                         // Just DNS4
	test("/dns6/abc.com/udp/1234", "udp6", "abc.com:1234")          // DNS6:port
	test("/dns6/abc.com", "ip6", "abc.com")                         // Just DNS6
}

func TestMultiaddrToIPNet(t *testing.T) {
	type testCase struct {
		name      string
		ma        string
		ips       []string
		contained []bool
	}

	testCases := []testCase{
		{
			name:      "basic",
			ma:        "/ip4/1.2.3.0/ipcidr/24",
			ips:       []string{"1.2.3.4", "1.2.3.9", "2.1.1.1"},
			contained: []bool{true, true, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ma := ma.StringCast(tc.ma)

			ipnet, err := MultiaddrToIPNet(ma)
			if err != nil {
				t.Fatalf("failed to parse multiaddr %v into ipnet", ma)
			}
			for i, ipString := range tc.ips {
				ip := net.ParseIP(ipString)
				if ip == nil {
					t.Fatalf("failed to parse IP %s", ipString)
				}
				if ipnet.Contains(ip) != tc.contained[i] {
					t.Fatalf("Contains check failed. Expected %v got %v", tc.contained[i], ipnet.Contains(ip))
				}
			}
		})
	}
}

func TestFailMultiaddrToIPNet(t *testing.T) {
	type testCase struct {
		name string
		ma   string
	}

	testCases := []testCase{
		{name: "missing ip addr", ma: "/ipcidr/24"},
		{name: "wrong mask", ma: "/ip4/1.2.3.0/ipcidr/128"},
		{name: "wrong mask", ma: "/ip6/::/ipcidr/255"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ma := ma.StringCast(tc.ma)

			_, err := MultiaddrToIPNet(ma)
			if err == nil {
				t.Fatalf("Expected error when parsing: %s", tc.ma)
			}
		})
	}
}
