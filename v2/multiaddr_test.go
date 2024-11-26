package multiaddr

import (
	"bytes"
	"fmt"
	"testing"
)

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// TestAssertForm is an example of how a user of a multiaddr should validate
// that the multiaddr is of some form.
func TestAssertForm(t *testing.T) {
	expectedStr := "/ip4/127.0.0.1/tcp/1234"
	ma, err := FromString(expectedStr)
	assertNoErr(t, err)

	ma, c := ma.PopLast()
	if c.Code() != P_TCP {
		t.Fatalf("expected %q, got %q", P_TCP, c.Code())
	}

	ma, c = ma.PopLast()
	switch c.Code() {
	case P_IP4, P_IP6:
	default:
		t.Fatalf("expected IP component, got %q", c.Code())
	}

	if len(ma) != 0 {
		t.Fatalf("expected empty Multiaddr, got %v", ma)
	}
}

func TestRoundTrip(t *testing.T) {
	expectedStr := "/ip4/127.0.0.1/tcp/1234"
	ma, err := FromString(expectedStr)
	assertNoErr(t, err)
	str, err := ma.ToString()
	assertNoErr(t, err)

	if str != expectedStr {
		t.Fatalf("expected %q, got %q", expectedStr, str)
	}
}

func FuzzRoundTrip(f *testing.F) {
	corpus := []string{
		"/ip4/1.2.3.4",
		"/ip4/0.0.0.0",
		"/ip4/192.0.2.0/ipcidr/24",
		"/ip6/::1",
		"/ip6/2601:9:4f81:9700:803e:ca65:66e8:c21",
		"/ip6/2601:9:4f81:9700:803e:ca65:66e8:c21/udp/1234/quic",
		"/ip6/2601:9:4f81:9700:803e:ca65:66e8:c21/udp/1234/quic-v1",
		"/ip6/2001:db8::/ipcidr/32",
		"/ip6zone/x/ip6/fe80::1",
		"/ip6zone/x%y/ip6/fe80::1",
		"/ip6zone/x%y/ip6/::",
		"/ip6zone/x/ip6/fe80::1/udp/1234/quic",
		"/ip6zone/x/ip6/fe80::1/udp/1234/quic-v1",
		"/onion/timaq4ygg2iegci7:1234",
		"/onion/timaq4ygg2iegci7:80/http",
		"/onion3/vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd:1234",
		"/onion3/vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd:80/http",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA/http",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA/udp/8080",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA/tcp/8080",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuqzwas",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuqzwassw",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq/http",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq/tcp/8080",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq/udp/8080",
		"/udp/0",
		"/tcp/0",
		"/sctp/0",
		"/udp/1234",
		"/tcp/1234",
		"/sctp/1234",
		"/udp/65535",
		"/tcp/65535",
		"/ipfs/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC",
		"/ipfs/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7",
		"/p2p/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC",
		"/p2p/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7",
		"/p2p/bafzbeigvf25ytwc3akrijfecaotc74udrhcxzh2cx3we5qqnw5vgrei4bm",
		"/p2p/12D3KooWCryG7Mon9orvQxcS1rYZjotPgpwoJNHHKcLLfE4Hf5mV",
		"/p2p/k51qzi5uqu5dhb6l8spkdx7yxafegfkee5by8h7lmjh2ehc2sgg34z7c15vzqs",
		"/p2p/bafzaajaiaejcalj543iwv2d7pkjt7ykvefrkfu7qjfi6sduakhso4lay6abn2d5u",
		"/udp/1234/sctp/1234",
		"/udp/1234/udt",
		"/udp/1234/utp",
		"/tcp/1234/http",
		"/tcp/1234/tls/http",
		"/tcp/1234/https",
		"/ipfs/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC/tcp/1234",
		"/ipfs/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7/tcp/1234",
		"/p2p/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC/tcp/1234",
		"/p2p/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7/tcp/1234",
		"/ip4/127.0.0.1/udp/1234",
		"/ip4/127.0.0.1/udp/0",
		"/ip4/127.0.0.1/tcp/1234",
		"/ip4/127.0.0.1/tcp/1234/",
		"/ip4/127.0.0.1/udp/1234/quic",
		"/ip4/127.0.0.1/udp/1234/quic-v1",
		"/ip4/127.0.0.1/udp/1234/quic-v1/webtransport",
		"/ip4/127.0.0.1/udp/1234/quic-v1/webtransport/certhash/b2uaraocy6yrdblb4sfptaddgimjmmpy",
		"/ip4/127.0.0.1/udp/1234/quic-v1/webtransport/certhash/b2uaraocy6yrdblb4sfptaddgimjmmpy/certhash/zQmbWTwYGcmdyK9CYfNBcfs9nhZs17a6FQ4Y8oea278xx41",
		"/ip4/127.0.0.1/ipfs/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC",
		"/ip4/127.0.0.1/ipfs/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC/tcp/1234",
		"/ip4/127.0.0.1/ipfs/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7",
		"/ip4/127.0.0.1/ipfs/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7/tcp/1234",
		"/ip4/127.0.0.1/p2p/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC",
		"/ip4/127.0.0.1/p2p/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC/tcp/1234",
		"/ip4/127.0.0.1/p2p/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7",
		"/ip4/127.0.0.1/p2p/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7/tcp/1234",
		"/unix/a/b/c/d/e",
		"/unix/stdio",
		"/ip4/1.2.3.4/tcp/80/unix/a/b/c/d/e/f",
		"/ip4/127.0.0.1/ipfs/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC/tcp/1234/unix/stdio",
		"/ip4/127.0.0.1/ipfs/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7/tcp/1234/unix/stdio",
		"/ip4/127.0.0.1/p2p/QmcgpsyWgH8Y8ajJz1Cu72KnS5uo2Aa2LpzU7kinSupNKC/tcp/1234/unix/stdio",
		"/ip4/127.0.0.1/p2p/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl7/tcp/1234/unix/stdio",
		"/ip4/127.0.0.1/tcp/9090/http/p2p-webrtc-direct",
		"/ip4/127.0.0.1/tcp/127/ws",
		"/ip4/127.0.0.1/tcp/127/ws",
		"/ip4/127.0.0.1/tcp/127/tls",
		"/ip4/127.0.0.1/tcp/127/tls/ws",
		"/ip4/127.0.0.1/tcp/127/noise",
		"/ip4/127.0.0.1/tcp/127/wss",
		"/ip4/127.0.0.1/tcp/127/wss",
		"/ip4/127.0.0.1/tcp/127/webrtc-direct",
		"/ip4/127.0.0.1/tcp/127/webrtc",
		"/http-path/tmp%2Fbar",
		"/http-path/tmp%2Fbar%2Fbaz",
		"/http-path/foo",
		"/ip4/127.0.0.1/tcp/0/p2p/12D3KooWCryG7Mon9orvQxcS1rYZjotPgpwoJNHHKcLLfE4Hf5mV/http-path/foo",
		"/ip4/127.0.0.1/tcp/443/tls/sni/example.com/http/http-path/foo",
		"/memory/4",
	}
	for _, d := range corpus {
		f.Add(d)
	}

	f.Fuzz(func(t *testing.T, data string) {
		ma, err := FromString(data)
		if err != nil {
			return
		}

		b, err := ma.ToBinary()
		if err != nil {
			panic("failed to roundtrip to bytes")
		}

		s, err := ma.ToString()
		if err != nil {
			panic("failed to roundtrip to string: " + err.Error() + " " + data)
		}

		ma2, err := FromString(s)
		if err != nil {
			panic("failed to roundtrip from string: " + err.Error() + " " + s + data)
		}
		b2, err := ma2.ToBinary()
		if err != nil {
			panic("failed to roundtrip to bytes second time")
		}

		if !bytes.Equal(b, b2) {
			panic(fmt.Sprintf("expected %q, got %q", data, s))
		}
	})
}
