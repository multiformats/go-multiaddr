package multiaddr

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
)

func newMultiaddr(t *testing.T, a string) Multiaddr {
	m, err := NewMultiaddr(a)
	if err != nil {
		t.Error(err)
	}
	return m
}

func TestReturnsNilOnEmpty(t *testing.T) {
	a := StringCast("/ip4/1.2.3.4")
	a, _ = SplitLast(a)
	require.Nil(t, a)
	a, _ = SplitLast(a)
	require.Nil(t, a)

	a, c := SplitLast(nil)
	require.Zero(t, len(a.Protocols()))
	require.Nil(t, a)
	require.Nil(t, c)

	// Test that empty multiaddr from various operations returns nil
	a = StringCast("/ip4/1.2.3.4/tcp/1234")
	_, a = SplitFirst(a)
	_, a = SplitFirst(a)
	require.Nil(t, a)
	_, a = SplitFirst(a)
	require.Nil(t, a)

	c, a = SplitFirst(nil)
	require.Nil(t, a)
	require.Nil(t, c)

	a = StringCast("/ip4/1.2.3.4/tcp/1234")
	a = a.Decapsulate(a)
	require.Nil(t, a)

	a = StringCast("/ip4/1.2.3.4/tcp/1234")
	a = a.Decapsulate(StringCast("/tcp/1234"))
	a = a.Decapsulate(StringCast("/ip4/1.2.3.4"))
	require.Nil(t, a)

	// Test that SplitFunc returns nil when we split at beginning and end
	a = StringCast("/ip4/1.2.3.4/tcp/1234")
	pre, _ := SplitFunc(a, func(c Component) bool {
		return c.Protocol().Code == P_IP4
	})
	require.Nil(t, pre)

	a = StringCast("/ip4/1.2.3.4/tcp/1234")
	_, post := SplitFunc(a, func(c Component) bool {
		return false
	})
	require.Nil(t, post)

	_, err := NewMultiaddr("")
	require.Error(t, err)

	var nilMultiaddr Multiaddr
	a = nilMultiaddr.AppendComponent()
	require.Nil(t, a)

	a = Join()
	require.Nil(t, a)
}

func TestJoinWithComponents(t *testing.T) {
	var m Multiaddr
	c, err := NewComponent("ip4", "127.0.0.1")
	require.NoError(t, err)

	expected := "/ip4/127.0.0.1"
	require.Equal(t, expected, Join(m, c).String())

}

func TestConstructFails(t *testing.T) {
	cases := []string{
		"/ip4",
		"/ip4/::1",
		"/ip4/fdpsofodsajfdoisa",
		"/ip4/::/ipcidr/256",
		"/ip6/::/ipcidr/1026",
		"/ip6",
		"/ip6zone",
		"/ip6zone/",
		"/ip6zone//ip6/fe80::1",
		"/udp",
		"/tcp",
		"/sctp",
		"/udp/65536",
		"/tcp/65536",
		"/quic/65536",
		"/quic-v1/65536",
		"/onion/9imaq4ygg2iegci7:80",
		"/onion/aaimaq4ygg2iegci7:80",
		"/onion/timaq4ygg2iegci7:0",
		"/onion/timaq4ygg2iegci7:-1",
		"/onion/timaq4ygg2iegci7",
		"/onion/timaq4ygg2iegci@:666",
		"/onion3/9ww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd:80",
		"/onion3/vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd7:80",
		"/onion3/vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd:0",
		"/onion3/vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd:-1",
		"/onion3/vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd",
		"/onion3/vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyy@:666",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA7:80",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA:0",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA:0",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA:-1",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA@:666",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA7:80",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA:0",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA:0",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA:-1",
		"/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA@:666",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzu",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzu77",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzu:80",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq:-1",
		"/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzu@",
		"/udp/1234/sctp",
		"/udp/1234/udt/1234",
		"/udp/1234/utp/1234",
		"/ip4/127.0.0.1/udp/jfodsajfidosajfoidsa",
		"/ip4/127.0.0.1/udp",
		"/ip4/127.0.0.1/tcp/jfodsajfidosajfoidsa",
		"/ip4/127.0.0.1/tcp",
		"/ip4/127.0.0.1/quic/1234",
		"/ip4/127.0.0.1/quic-v1/1234",
		"/ip4/127.0.0.1/udp/1234/quic-v1/webtransport/certhash",
		"/ip4/127.0.0.1/udp/1234/quic-v1/webtransport/certhash/b2uaraocy6yrdblb4sfptaddgimjmmp", // 1 character missing from certhash
		"/ip4/127.0.0.1/ipfs",
		"/ip4/127.0.0.1/ipfs/tcp",
		"/ip4/127.0.0.1/p2p",
		"/ip4/127.0.0.1/p2p/tcp",
		"/unix",
		"/ip4/1.2.3.4/tcp/80/unix",
		"/ip4/1.2.3.4/tcp/-1",
		"/ip4/127.0.0.1/tcp/9090/http/p2p-webcrt-direct",
		fmt.Sprintf("/memory/%d1", uint64(1<<63)),
		"/",
		"",
		"/p2p/QmxoHT6iViN5xAjoz1VZ553cL31U9F94ht3QvWR1FrEbZY", // sha256 multihash with digest len > 32
	}

	for _, a := range cases {
		if _, err := NewMultiaddr(a); err == nil {
			t.Errorf("should have failed: %s - %s", a, err)
		}
	}
}

func TestEmptyMultiaddr(t *testing.T) {
	_, err := NewMultiaddrBytes([]byte{})
	if err == nil {
		t.Fatal("should have failed to parse empty multiaddr")
	}
}

var good = []string{
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

func TestConstructSucceeds(t *testing.T) {
	for _, a := range good {
		if _, err := NewMultiaddr(a); err != nil {
			t.Errorf("should have succeeded: %s -- %s", a, err)
		}
	}
}

func TestEqual(t *testing.T) {
	m1 := newMultiaddr(t, "/ip4/127.0.0.1/udp/1234")
	m2 := newMultiaddr(t, "/ip4/127.0.0.1/tcp/1234")
	m3 := newMultiaddr(t, "/ip4/127.0.0.1/tcp/1234")
	m4 := newMultiaddr(t, "/ip4/127.0.0.1/tcp/1234/")

	if m1.Equal(m2) {
		t.Error("should not be equal")
	}

	if m2.Equal(m1) {
		t.Error("should not be equal")
	}

	if !m2.Equal(m3) {
		t.Error("should be equal")
	}

	if !m3.Equal(m2) {
		t.Error("should be equal")
	}

	if !m1.Equal(m1) {
		t.Error("should be equal")
	}

	if !m2.Equal(m4) {
		t.Error("should be equal")
	}

	if !m4.Equal(m3) {
		t.Error("should be equal")
	}
}

// TestNilInterface makes sure funcs that accept a multiaddr interface don't
// panic if it's passed a nil interface.
func TestNilInterface(t *testing.T) {
	m1 := newMultiaddr(t, "/ip4/127.0.0.1/udp/1234")
	var m2 Multiaddr
	m1.Equal(m2)
	m1.Encapsulate(m2)
	m1.Decapsulate(m2)

	// Test components
	c, _ := SplitFirst(m1)
	c.Multiaddr().Equal(m2)
	c.Encapsulate(m2)
	c.Decapsulate(m2)

	// Util funcs
	_ = Split(m2)
	_, _ = SplitFirst(m2)
	_, _ = SplitLast(m2)
	ForEach(m2, func(c Component) bool { return true })
}

func TestStringToBytes(t *testing.T) {

	testString := func(s string, h string) {
		b1, err := hex.DecodeString(h)
		if err != nil {
			t.Error("failed to decode hex", h)
		}

		// t.Log("196", h, []byte(b1))

		b2, err := stringToBytes(s)
		if err != nil {
			t.Error("failed to convert", s, err)
		}

		if !bytes.Equal(b1, b2) {
			t.Error("failed to convert \n", s, "to\n", hex.EncodeToString(b1), "got\n", hex.EncodeToString(b2))
		}

		if _, err := NewMultiaddrBytes(b2); err != nil {
			t.Error(err, "len:", len(b2))
		}
	}

	testString("/ip4/127.0.0.1/udp/1234", "047f000001910204d2")
	testString("/ip4/127.0.0.1/tcp/4321", "047f0000010610e1")
	testString("/ip4/127.0.0.1/udp/1234/ip4/127.0.0.1/tcp/4321", "047f000001910204d2047f0000010610e1")
	testString("/onion/aaimaq4ygg2iegci:80", "bc030010c0439831b48218480050")
	testString("/onion3/vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd:1234", "bd03adadec040be047f9658668b11a504f3155001f231a37f54c4476c07fb4cc139ed7e30304d2")
	testString("/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA",
		"be0383038d3fc8c976a86ae4e78ba378e75ec41bc9ab1542a9cb422581987e118f5cb0c024f3639d6ad9b3aff613672f07bfbbbfc2f920ef910534ecaa6ff9c03e0fa4872a764d2fce6d4cfc5a5a9800cd95944cc9ef0241f753fe71494a175f334b35682459acadc4076428ab49b5a83a49d2ea2366b06461e4a559b0111fa750e0de0c138a94d1231ed5979572ff53922905636221994bdabc44bd0c17fef11622b16432db3f193400af53cc61aa9bfc0c4c8d874b41a6e18732f0b60f5662ef1a89c80589dd8366c90bb58bb85ead56356aba2a244950ca170abbd01094539014f84bdd383e4a10e00cee63dfc3e809506e2d9b54edbdca1bace6eaa119e68573d30533791fba830f5d80be5c051a77c09415e3b8fe3139400848be5244b8ae96bb0c4a24f819cba0488f34985eac741d3359180bd72cafa1559e4c19f54ea8cedbb6a5afde4319396eb92aab340c60a50cc2284580cb3ad09017e8d9abc60269b3d8d687680bd86ce834412273d4f2e3bf68dd3d6fe87e2426ac658cd5c77fd5c0aa000000")
	testString("/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq",
		"bf0320efbcd45d0c5dc79781ac6f20ea5055a036afb48d45a52e7d68ec7d4338919e69")

}

func TestBytesToString(t *testing.T) {
	testString := func(s1 string, h string) {
		t.Helper()
		b, err := hex.DecodeString(h)
		if err != nil {
			t.Error("failed to decode hex", h)
		}

		if _, err := NewMultiaddrBytes(b); err != nil {
			t.Error(err)
		}

		m, err := NewMultiaddrBytes(b)
		s2 := m.String()
		if err != nil {
			t.Log("236", s1, ":", string(h), ":", s2)
			t.Error("failed to convert", b, err)
		}

		if s1 != s2 {
			t.Error("failed to convert", b, "to", s1, "got", s2)
		}
	}

	testString("/ip4/127.0.0.1/udp/1234", "047f000001910204d2")
	testString("/ip4/127.0.0.1/tcp/4321", "047f0000010610e1")
	testString("/ip4/127.0.0.1/udp/1234/ip4/127.0.0.1/tcp/4321", "047f000001910204d2047f0000010610e1")
	testString("/onion/aaimaq4ygg2iegci:80", "bc030010c0439831b48218480050")
	testString("/onion3/vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd:1234", "bd03adadec040be047f9658668b11a504f3155001f231a37f54c4476c07fb4cc139ed7e30304d2")
	testString("/garlic64/jT~IyXaoauTni6N4517EG8mrFUKpy0IlgZh-EY9csMAk82Odatmzr~YTZy8Hv7u~wvkg75EFNOyqb~nAPg-khyp2TS~ObUz8WlqYAM2VlEzJ7wJB91P-cUlKF18zSzVoJFmsrcQHZCirSbWoOknS6iNmsGRh5KVZsBEfp1Dg3gwTipTRIx7Vl5Vy~1OSKQVjYiGZS9q8RL0MF~7xFiKxZDLbPxk0AK9TzGGqm~wMTI2HS0Gm4Ycy8LYPVmLvGonIBYndg2bJC7WLuF6tVjVquiokSVDKFwq70BCUU5AU-EvdOD5KEOAM7mPfw-gJUG4tm1TtvcobrObqoRnmhXPTBTN5H7qDD12AvlwFGnfAlBXjuP4xOUAISL5SRLiulrsMSiT4GcugSI80mF6sdB0zWRgL1yyvoVWeTBn1TqjO27alr95DGTluuSqrNAxgpQzCKEWAyzrQkBfo2avGAmmz2NaHaAvYbOg0QSJz1PLjv2jdPW~ofiQmrGWM1cd~1cCqAAAA",
		"be0383038d3fc8c976a86ae4e78ba378e75ec41bc9ab1542a9cb422581987e118f5cb0c024f3639d6ad9b3aff613672f07bfbbbfc2f920ef910534ecaa6ff9c03e0fa4872a764d2fce6d4cfc5a5a9800cd95944cc9ef0241f753fe71494a175f334b35682459acadc4076428ab49b5a83a49d2ea2366b06461e4a559b0111fa750e0de0c138a94d1231ed5979572ff53922905636221994bdabc44bd0c17fef11622b16432db3f193400af53cc61aa9bfc0c4c8d874b41a6e18732f0b60f5662ef1a89c80589dd8366c90bb58bb85ead56356aba2a244950ca170abbd01094539014f84bdd383e4a10e00cee63dfc3e809506e2d9b54edbdca1bace6eaa119e68573d30533791fba830f5d80be5c051a77c09415e3b8fe3139400848be5244b8ae96bb0c4a24f819cba0488f34985eac741d3359180bd72cafa1559e4c19f54ea8cedbb6a5afde4319396eb92aab340c60a50cc2284580cb3ad09017e8d9abc60269b3d8d687680bd86ce834412273d4f2e3bf68dd3d6fe87e2426ac658cd5c77fd5c0aa000000")
	testString("/garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuq",
		"bf0320efbcd45d0c5dc79781ac6f20ea5055a036afb48d45a52e7d68ec7d4338919e69")
}

func TestBytesSplitAndJoin(t *testing.T) {

	testString := func(s string, res []string) {
		m, err := NewMultiaddr(s)
		if err != nil {
			t.Fatal("failed to convert", s, err)
		}

		split := Split(m)
		if len(split) != len(res) {
			t.Error("not enough split components", split)
			return
		}

		for i, a := range split {
			if a.String() != res[i] {
				t.Errorf("split component failed: %s != %s", &a, res[i])
			}
		}

		joined := append(Multiaddr{}, split...)
		if !m.Equal(joined) {
			t.Errorf("joined components failed: %s != %s", m, joined)
		}

		for i, a := range split {
			if a.String() != res[i] {
				t.Errorf("split component failed: %s != %s", &a, res[i])
			}
		}
	}

	testString("/ip4/1.2.3.4/udp/1234", []string{"/ip4/1.2.3.4", "/udp/1234"})
	testString("/ip4/1.2.3.4/tcp/1/ip4/2.3.4.5/udp/2",
		[]string{"/ip4/1.2.3.4", "/tcp/1", "/ip4/2.3.4.5", "/udp/2"})
	testString("/ip4/1.2.3.4/utp/ip4/2.3.4.5/udp/2/udt",
		[]string{"/ip4/1.2.3.4", "/utp", "/ip4/2.3.4.5", "/udp/2", "/udt"})
}

func TestProtocols(t *testing.T) {
	m, err := NewMultiaddr("/ip4/127.0.0.1/udp/1234")
	if err != nil {
		t.Error("failed to construct", "/ip4/127.0.0.1/udp/1234")
	}

	ps := m.Protocols()
	if ps[0].Code != ProtocolWithName("ip4").Code {
		t.Error(ps[0], ProtocolWithName("ip4"))
		t.Error("failed to get ip4 protocol")
	}

	if ps[1].Code != ProtocolWithName("udp").Code {
		t.Error(ps[1], ProtocolWithName("udp"))
		t.Error("failed to get udp protocol")
	}

}

func TestProtocolsWithString(t *testing.T) {
	pwn := ProtocolWithName
	good := map[string][]Protocol{
		"/ip4":                    {pwn("ip4")},
		"/ip4/tcp":                {pwn("ip4"), pwn("tcp")},
		"ip4/tcp/udp/ip6":         {pwn("ip4"), pwn("tcp"), pwn("udp"), pwn("ip6")},
		"////////ip4/tcp":         {pwn("ip4"), pwn("tcp")},
		"ip4/udp/////////":        {pwn("ip4"), pwn("udp")},
		"////////ip4/tcp////////": {pwn("ip4"), pwn("tcp")},
	}

	for s, ps1 := range good {
		ps2, err := ProtocolsWithString(s)
		if err != nil {
			t.Errorf("ProtocolsWithString(%s) should have succeeded", s)
		}

		for i, ps1p := range ps1 {
			ps2p := ps2[i]
			if ps1p.Code != ps2p.Code {
				t.Errorf("mismatch: %s != %s, %s", ps1p.Name, ps2p.Name, s)
			}
		}
	}

	bad := []string{
		"dsijafd",                           // bogus proto
		"/ip4/tcp/fidosafoidsa",             // bogus proto
		"////////ip4/tcp/21432141/////////", // bogus proto
		"////////ip4///////tcp/////////",    // empty protos in between
	}

	for _, s := range bad {
		if _, err := ProtocolsWithString(s); err == nil {
			t.Errorf("ProtocolsWithString(%s) should have failed", s)
		}
	}

}

func TestEncapsulate(t *testing.T) {
	m, err := NewMultiaddr("/ip4/127.0.0.1/udp/1234")
	if err != nil {
		t.Error(err)
	}

	m2, err := NewMultiaddr("/udp/5678")
	if err != nil {
		t.Error(err)
	}

	b := m.Encapsulate(m2)
	if s := b.String(); s != "/ip4/127.0.0.1/udp/1234/udp/5678" {
		t.Error("encapsulate /ip4/127.0.0.1/udp/1234/udp/5678 failed.", s)
	}

	m3, _ := NewMultiaddr("/udp/5678")
	c := b.Decapsulate(m3)
	if s := c.String(); s != "/ip4/127.0.0.1/udp/1234" {
		t.Error("decapsulate /udp failed.", "/ip4/127.0.0.1/udp/1234", s)
	}

	m4, _ := NewMultiaddr("/ip4/127.0.0.1")
	d := c.Decapsulate(m4)
	if d != nil {
		t.Error("decapsulate /ip4 failed: ", d)
	}

	t.Run("Encapsulating with components", func(t *testing.T) {
		left, last := SplitLast(m)
		joined := left.Encapsulate(last)
		require.True(t, joined.Equal(m))

		first, rest := SplitFirst(m)
		joined = first.Encapsulate(rest)
		require.True(t, joined.Equal(m))
		// Component type
		joined = (*first).Encapsulate(rest)
		require.True(t, joined.Equal(m))
	})
}

func TestDecapsulateComment(t *testing.T) {
	// shows the behavior from the interface comment
	m := StringCast("/ip4/1.2.3.4/tcp/80")
	rest := m.Decapsulate(StringCast("/tcp/80"))
	if rest.String() != "/ip4/1.2.3.4" {
		t.Fatalf("Documented behavior is not correct. Expected %v saw %v", "/ip4/1.2.3.4/", rest.String())
	}

	m = StringCast("/ip4/1.2.3.4/tcp/80")
	rest = m.Decapsulate(StringCast("/udp/80"))
	if !rest.Equal(m) {
		t.Fatalf("Documented behavior is not correct. Expected %v saw %v", "/ip4/1.2.3.4/tcp/80", rest.String())
	}

	m = StringCast("/ip4/1.2.3.4/tcp/80")
	rest = m.Decapsulate(StringCast("/ip4/1.2.3.4"))
	require.Nil(t, rest, "expected a nil multiaddr if we decapsulate everything")
}

func TestDecapsulate(t *testing.T) {
	t.Run("right is nil", func(t *testing.T) {
		left := StringCast("/ip4/1.2.3.4/tcp/1")
		var right Multiaddr
		left.Decapsulate(right)
	})

	testcases := []struct {
		left, right, expected string
	}{
		{"/ip4/1.2.3.4/tcp/1234", "/ip4/1.2.3.4", ""},
		{"/ip4/1.2.3.4", "/ip4/1.2.3.4/tcp/1234", "/ip4/1.2.3.4"},
		{"/ip4/1.2.3.5/tcp/1234", "/ip4/5.3.2.1", "/ip4/1.2.3.5/tcp/1234"},
		{"/ip4/1.2.3.5/udp/1234/quic-v1", "/udp/1234", "/ip4/1.2.3.5"},
		{"/ip4/1.2.3.6/udp/1234/quic-v1", "/udp/1234/quic-v1", "/ip4/1.2.3.6"},
		{"/ip4/1.2.3.7/tcp/1234", "/ws", "/ip4/1.2.3.7/tcp/1234"},
		{"/dnsaddr/wss.com/tcp/4001", "/ws", "/dnsaddr/wss.com/tcp/4001"},
		{"/dnsaddr/wss.com/tcp/4001/ws", "/wss", "/dnsaddr/wss.com/tcp/4001/ws"},
		{"/dnsaddr/wss.com/ws", "/wss", "/dnsaddr/wss.com/ws"},
		{"/dnsaddr/wss.com/ws", "/dnsaddr/wss.com", ""},
		{"/dnsaddr/wss.com/tcp/4001/wss", "/wss", "/dnsaddr/wss.com/tcp/4001"},
	}

	for _, tc := range testcases {
		t.Run(tc.left, func(t *testing.T) {
			left := StringCast(tc.left)
			right := StringCast(tc.right)
			actualMa := left.Decapsulate(right)

			if tc.expected == "" {
				require.Nil(t, actualMa, "expected nil")
				return
			}

			actual := actualMa.String()
			expected := StringCast(tc.expected).String()
			require.Equal(t, expected, actual)
		})
	}

	for _, tc := range testcases {
		t.Run("Decapsulating with components"+tc.left, func(t *testing.T) {
			left, last := SplitLast(StringCast(tc.left))
			butLast := left.Decapsulate(last)
			require.Equal(t, butLast.String(), left.String())
			// Round trip
			require.Equal(t, tc.left, butLast.Encapsulate(last).String())
		})
	}
}

func assertValueForProto(t *testing.T, a Multiaddr, p int, exp string) {
	t.Logf("checking for %s in %s", ProtocolWithCode(p).Name, a)
	fv, err := a.ValueForProtocol(p)
	if err != nil {
		t.Fatal(err)
	}

	if fv != exp {
		t.Fatalf("expected %q for %d in %s, but got %q instead", exp, p, a, fv)
	}
}

func TestAppendComponent(t *testing.T) {
	var m Multiaddr
	res := m.AppendComponent(nil)
	require.Equal(t, m, res)

	c, err := NewComponent("ip4", "127.0.0.1")
	require.NoError(t, err)
	res = m.AppendComponent(c)
	require.Equal(t, "/ip4/127.0.0.1", res.String())
}

func TestGetValue(t *testing.T) {
	a := newMultiaddr(t, "/ip4/127.0.0.1/utp/tcp/5555/udp/1234/tls/utp/ipfs/QmbHVEEepCi7rn7VL7Exxpd2Ci9NNB6ifvqwhsrbRMgQFP")
	assertValueForProto(t, a, P_IP4, "127.0.0.1")
	assertValueForProto(t, a, P_UTP, "")
	assertValueForProto(t, a, P_TLS, "")
	assertValueForProto(t, a, P_TCP, "5555")
	assertValueForProto(t, a, P_UDP, "1234")
	assertValueForProto(t, a, P_IPFS, "QmbHVEEepCi7rn7VL7Exxpd2Ci9NNB6ifvqwhsrbRMgQFP")
	assertValueForProto(t, a, P_P2P, "QmbHVEEepCi7rn7VL7Exxpd2Ci9NNB6ifvqwhsrbRMgQFP")

	_, err := a.ValueForProtocol(P_IP6)
	switch err {
	case ErrProtocolNotFound:
		break
	case nil:
		t.Fatal("expected value lookup to fail")
	default:
		t.Fatalf("expected ErrProtocolNotFound but got: %s", err)
	}

	a = newMultiaddr(t, "/ip4/0.0.0.0") // only one addr
	assertValueForProto(t, a, P_IP4, "0.0.0.0")

	a = newMultiaddr(t, "/ip4/0.0.0.0/ip4/0.0.0.0/ip4/0.0.0.0") // same sub-addr
	assertValueForProto(t, a, P_IP4, "0.0.0.0")

	a = newMultiaddr(t, "/ip4/0.0.0.0/udp/12345/utp") // ending in a no-value one.
	assertValueForProto(t, a, P_IP4, "0.0.0.0")
	assertValueForProto(t, a, P_UDP, "12345")
	assertValueForProto(t, a, P_UTP, "")

	a = newMultiaddr(t, "/ip4/0.0.0.0/unix/a/b/c/d") // ending in a path one.
	assertValueForProto(t, a, P_IP4, "0.0.0.0")
	assertValueForProto(t, a, P_UNIX, "/a/b/c/d")
}

func FuzzNewMultiaddrBytes(f *testing.F) {
	for _, v := range good {
		ma, err := NewMultiaddr(v)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(ma.Bytes())
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		// just checking that it doesn't panic
		ma, err := NewMultiaddrBytes(b)
		if err == nil {
			// for any valid multiaddrs, make sure these calls don't panic
			ma.Protocols()
			roundTripBytes(t, ma)
			roundTripString(t, ma)
		}
	})
}

func FuzzNewMultiaddrString(f *testing.F) {
	for _, v := range good {
		if _, err := NewMultiaddr(v); err != nil {
			// Validate maddrs
			f.Fatal(err)
		}
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, s string) {
		// just checking that it doesn't panic
		ma, err := NewMultiaddr(s)
		if err == nil {
			// for any valid multiaddrs, make sure these calls don't panic
			ma.Protocols()
			roundTripBytes(t, ma)
			roundTripString(t, ma)
		}
	})
}

func roundTripBytes(t *testing.T, orig Multiaddr) {
	m2, err := NewMultiaddrBytes(orig.Bytes())
	if err != nil {
		t.Fatalf("failed to parse maddr back from ma.Bytes, %v: %v", orig, err)
	}
	if !m2.Equal(orig) {
		t.Fatalf("unequal maddr after roundTripBytes %v %v", orig, m2)
	}
}

func roundTripString(t *testing.T, orig Multiaddr) {
	m2, err := NewMultiaddr(orig.String())
	if err != nil {
		t.Fatalf("failed to parse maddr back from ma.String, %v: %v", orig, err)
	}
	if !m2.Equal(orig) {
		t.Fatalf("unequal maddr after roundTripString %v %v\n% 02x\n% 02x\n", orig, m2, orig.Bytes(), m2.Bytes())
	}
}

func TestBinaryRepresentation(t *testing.T) {
	expected := []byte{0x4, 0x7f, 0x0, 0x0, 0x1, 0x91, 0x2, 0x4, 0xd2}
	ma, err := NewMultiaddr("/ip4/127.0.0.1/udp/1234")
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(ma.Bytes(), expected) {
		t.Errorf("expected %x, got %x", expected, ma.Bytes())
	}
}

func TestRoundTrip(t *testing.T) {
	for _, s := range []string{
		"/unix/a/b/c/d",
		"/ip6/::ffff:127.0.0.1/tcp/111",
		"/ip4/127.0.0.1/tcp/123",
		"/ip4/127.0.0.1/tcp/123/tls",
		"/ip4/127.0.0.1/udp/123",
		"/ip4/127.0.0.1/udp/123/ip6/::",
		"/ip4/127.0.0.1/udp/1234/quic-v1/webtransport/certhash/uEiDDq4_xNyDorZBH3TlGazyJdOWSwvo4PUo5YHFMrvDE8g",
		"/p2p/QmbHVEEepCi7rn7VL7Exxpd2Ci9NNB6ifvqwhsrbRMgQFP",
		"/p2p/QmbHVEEepCi7rn7VL7Exxpd2Ci9NNB6ifvqwhsrbRMgQFP/unix/a/b/c",
		"/http-path/tmp%2Fbar",
	} {
		ma, err := NewMultiaddr(s)
		if err != nil {
			t.Errorf("error when parsing %q: %s", s, err)
			continue
		}
		if ma.String() != s {
			t.Errorf("failed to round trip %q", s)
		}
	}
}

func TestIPFSvP2P(t *testing.T) {
	var (
		p2pAddr  = "/p2p/QmbHVEEepCi7rn7VL7Exxpd2Ci9NNB6ifvqwhsrbRMgQFP"
		ipfsAddr = "/ipfs/QmbHVEEepCi7rn7VL7Exxpd2Ci9NNB6ifvqwhsrbRMgQFP"
	)

	for _, s := range []string{p2pAddr, ipfsAddr} {
		ma, err := NewMultiaddr(s)
		if err != nil {
			t.Errorf("error when parsing %q: %s", s, err)
		}
		if ma.String() != p2pAddr {
			t.Errorf("expected %q, got %q", p2pAddr, ma.String())
		}
	}
}

func TestInvalidP2PAddrBytes(t *testing.T) {
	badAddr := "a503221221c05877cbae039d70a5e600ea02c6f9f2942439285c9e344e26f8d280c850fad6"
	bts, err := hex.DecodeString(badAddr)
	if err != nil {
		t.Fatal(err)
	}
	ma, err := NewMultiaddrBytes(bts)
	if err == nil {
		t.Error("should have failed")
		// Check for panic
		_ = ma.String()
	}
}

func TestInvalidP2PAddrString(t *testing.T) {
	hashedData, err := mh.Sum([]byte("test"), mh.SHA2_256, -1)
	if err != nil {
		t.Fatal(err)
	}

	// using MD5 since it's not a valid data codec
	unknownCodecCID := cid.NewCidV1(mh.MD5, hashedData).String()

	badStringAddrs := []string{
		"/p2p/k2k4r8oqamigqdo6o7hsbfwd45y70oyynp98usk7zmyfrzpqxh1pohl-", // invalid multibase encoding
		"/p2p/?unknownmultibase", // invalid multibase encoding
		"/p2p/k2jmtxwoe2phm1hbqp0e7nufqf6umvuu2e9qd7ana7h411a0haqj6i2z", // non-libp2p-key codec
		"/p2p/" + unknownCodecCID, // impossible codec
	}
	for _, a := range badStringAddrs {
		ma, err := NewMultiaddr(a)
		if err == nil {
			t.Error("should have failed")
			// Check for panic
			_ = ma.String()
		}
	}
}

func TestZone(t *testing.T) {
	ip6String := "/ip6zone/eth0/ip6/::1"
	ip6Bytes := []byte{
		0x2a, 4,
		'e', 't', 'h', '0',
		0x29,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 1,
	}

	ma, err := NewMultiaddr(ip6String)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(ma.Bytes(), ip6Bytes) {
		t.Errorf("expected %x, got %x", ip6Bytes, ma.Bytes())
	}

	ma2, err2 := NewMultiaddrBytes(ip6Bytes)
	if err2 != nil {
		t.Error(err)
	}
	if ma2.String() != ip6String {
		t.Errorf("expected %s, got %s", ip6String, ma2.String())
	}
}

func TestBinaryMarshaler(t *testing.T) {
	addr := newMultiaddr(t, "/ip4/0.0.0.0/tcp/4001/tls")
	b, err := addr.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var addr2 Multiaddr
	if err = addr2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if !addr.Equal(addr2) {
		t.Error("expected equal addresses in circular marshaling test")
	}
}

func TestTextMarshaler(t *testing.T) {
	addr := newMultiaddr(t, "/ip4/0.0.0.0/tcp/4001/tls")
	b, err := addr.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	var addr2 Multiaddr
	if err = addr2.UnmarshalText(b); err != nil {
		t.Fatal(err)
	}
	if !addr.Equal(addr2) {
		t.Error("expected equal addresses in circular marshaling test")
	}
}

func TestJSONMarshaler(t *testing.T) {
	addr := newMultiaddr(t, "/ip4/0.0.0.0/tcp/4001/tls")
	b, err := addr.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var addr2 Multiaddr
	if err = addr2.UnmarshalJSON(b); err != nil {
		t.Fatal(err)
	}
	if !addr.Equal(addr2) {
		t.Error("expected equal addresses in circular marshaling test")
	}
}

func TestComponentBinaryMarshaler(t *testing.T) {
	comp, err := NewComponent("ip4", "0.0.0.0")
	if err != nil {
		t.Fatal(err)
	}
	b, err := comp.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var comp2 Component
	if err = comp2.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if !comp.Equal(&comp2) {
		t.Error("expected equal components in circular marshaling test")
	}
}

func TestComponentTextMarshaler(t *testing.T) {
	comp, err := NewComponent("ip4", "0.0.0.0")
	if err != nil {
		t.Fatal(err)
	}
	b, err := comp.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	var comp2 Component
	if err = comp2.UnmarshalText(b); err != nil {
		t.Fatal(err)
	}
	if !comp.Equal(&comp2) {
		t.Error("expected equal components in circular marshaling test")
	}
}

func TestComponentJSONMarshaler(t *testing.T) {
	comp, err := NewComponent("ip4", "0.0.0.0")
	if err != nil {
		t.Fatal(err)
	}
	b, err := comp.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var comp2 Component
	if err = comp2.UnmarshalJSON(b); err != nil {
		t.Fatal(err)
	}
	if !comp.Equal(&comp2) {
		t.Error("expected equal components in circular marshaling test")
	}
}

func TestUseNil(t *testing.T) {
	f := func() Multiaddr {
		return nil
	}

	_ = f()

	var foo Multiaddr = nil
	_, right := SplitFirst(foo)
	right.Protocols()
	foo.Protocols()
	foo.Bytes()
	foo.Compare(nil)
	foo.Decapsulate(nil)
	foo.Encapsulate(nil)
	foo.Equal(nil)
	_, _ = foo.MarshalBinary()
	_, _ = foo.MarshalJSON()
	_, _ = foo.MarshalText()
	foo.Protocols()
	_ = foo.String()
	_ = foo.UnmarshalBinary(nil)
	_ = foo.UnmarshalJSON(nil)
	_ = foo.UnmarshalText(nil)
	_, _ = foo.ValueForProtocol(0)
}

func TestUseNilComponent(t *testing.T) {
	var foo *Component
	foo.Multiaddr()
	foo.Encapsulate(nil)
	foo.Decapsulate(nil)
	require.True(t, foo == nil)
	foo.Bytes()
	foo.MarshalBinary()
	foo.MarshalJSON()
	foo.MarshalText()
	foo.UnmarshalBinary(nil)
	foo.UnmarshalJSON(nil)
	foo.UnmarshalText(nil)
	foo.Equal(nil)
	foo.Compare(nil)
	foo.Protocols()
	foo.ValueForProtocol(0)
	foo.Protocol()
	foo.RawValue()
	foo.Value()
	_ = foo.String()

	var m Multiaddr = nil
	m.Encapsulate(foo)
}

func TestFilterAddrs(t *testing.T) {
	bad := []Multiaddr{
		newMultiaddr(t, "/ip6/fe80::1/tcp/1234"),
		newMultiaddr(t, "/ip6/fe80::100/tcp/1234"),
	}
	good := []Multiaddr{
		newMultiaddr(t, "/ip4/127.0.0.1/tcp/1234"),
		newMultiaddr(t, "/ip4/1.1.1.1/tcp/999"),
		newMultiaddr(t, "/ip4/1.2.3.4/udp/1234/utp"),
	}
	goodAndBad := append(good, bad...)

	filter := func(addr Multiaddr) bool {
		return addr.Protocols()[0].Code == P_IP4
	}

	require.Empty(t, FilterAddrs(bad, filter))
	require.ElementsMatch(t, FilterAddrs(good, filter), good)
	require.ElementsMatch(t, FilterAddrs(goodAndBad, filter), good)
}

func TestContains(t *testing.T) {
	a1 := newMultiaddr(t, "/ip4/127.0.0.1/tcp/1234")
	a2 := newMultiaddr(t, "/ip4/1.1.1.1/tcp/999")
	a3 := newMultiaddr(t, "/ip4/1.2.3.4/udp/443/quic")
	a4 := newMultiaddr(t, "/ip4/1.2.3.4/udp/443/quic-v1")
	addrs := []Multiaddr{a1, a2, a3, a4}

	require.True(t, Contains(addrs, a1))
	require.True(t, Contains(addrs, a2))
	require.True(t, Contains(addrs, a3))
	require.True(t, Contains(addrs, a4))
	require.False(t, Contains(addrs, newMultiaddr(t, "/ip4/4.3.2.1/udp/1234/utp")))
	require.False(t, Contains(nil, a1))
}

func TestUniqueAddrs(t *testing.T) {
	tcpAddr := StringCast("/ip4/127.0.0.1/tcp/1234")
	quicAddr := StringCast("/ip4/127.0.0.1/udp/1234/quic-v1")
	wsAddr := StringCast("/ip4/127.0.0.1/tcp/1234/ws")

	type testcase struct {
		in, out []Multiaddr
	}

	for i, tc := range []testcase{
		{in: nil, out: nil},
		{in: []Multiaddr{tcpAddr}, out: []Multiaddr{tcpAddr}},
		{in: []Multiaddr{tcpAddr, tcpAddr, tcpAddr}, out: []Multiaddr{tcpAddr}},
		{in: []Multiaddr{tcpAddr, quicAddr, tcpAddr}, out: []Multiaddr{tcpAddr, quicAddr}},
		{in: []Multiaddr{tcpAddr, quicAddr, wsAddr}, out: []Multiaddr{tcpAddr, quicAddr, wsAddr}},
	} {
		tc := tc
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			deduped := Unique(tc.in)
			for _, a := range tc.out {
				require.Contains(t, deduped, a)
			}
		})
	}
}

func BenchmarkUniqueAddrs(b *testing.B) {
	b.ReportAllocs()
	var addrs []Multiaddr
	r := rand.New(rand.NewSource(1234))
	for i := 0; i < 100; i++ {
		tcpAddr := StringCast(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", r.Intn(math.MaxUint16)))
		quicAddr := StringCast(fmt.Sprintf("/ip4/127.0.0.1/udp/%d/quic-v1", r.Intn(math.MaxUint16)))
		wsAddr := StringCast(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/ws", r.Intn(math.MaxUint16)))
		addrs = append(addrs, tcpAddr, tcpAddr, quicAddr, quicAddr, wsAddr)
	}
	for _, sz := range []int{10, 20, 30, 50, 100} {
		b.Run(fmt.Sprintf("%d", sz), func(b *testing.B) {
			items := make([]Multiaddr, sz)
			for i := 0; i < b.N; i++ {
				copy(items, addrs[:sz])
				Unique(items)
			}
		})
	}
}

func TestDNS(t *testing.T) {
	b := []byte("7*000000000000000000000000000000000000000000")
	a, err := NewMultiaddrBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	aa := StringCast(a.String())
	if !a.Equal(aa) {
		t.Fatal("expected equality")
	}
}

func TestHTTPPath(t *testing.T) {
	t.Run("bad addr", func(t *testing.T) {
		badAddr := "/http-path/thisIsMissingAfullByte%f"
		_, err := NewMultiaddr(badAddr)
		require.Error(t, err)
	})

	t.Run("only reads the http-path part", func(t *testing.T) {
		addr := "/http-path/tmp%2Fbar/p2p-circuit" // The http-path only reference the part immediately after it. It does not include the rest of the multiaddr (like the /path component sometimes does)
		m, err := NewMultiaddr(addr)
		require.NoError(t, err)
		m.ValueForProtocol(P_HTTP_PATH)
		v, err := m.ValueForProtocol(P_HTTP_PATH)
		require.NoError(t, err)
		require.Equal(t, "tmp%2Fbar", v)
	})

	t.Run("round trip", func(t *testing.T) {
		cases := []string{
			"/http-path/tmp%2Fbar",
			"/http-path/tmp%2Fbar%2Fbaz",
			"/http-path/foo",
			"/ip4/127.0.0.1/tcp/0/p2p/12D3KooWCryG7Mon9orvQxcS1rYZjotPgpwoJNHHKcLLfE4Hf5mV/http-path/foo",
			"/ip4/127.0.0.1/tcp/443/tls/sni/example.com/http/http-path/foo",
		}
		for _, c := range cases {
			m, err := NewMultiaddr(c)
			require.NoError(t, err)
			require.Equal(t, c, m.String())
		}
	})

	t.Run("value for protocol", func(t *testing.T) {
		m := StringCast("/http-path/tmp%2Fbar")
		v, err := m.ValueForProtocol(P_HTTP_PATH)
		require.NoError(t, err)
		// This gives us the url escaped version
		require.Equal(t, "tmp%2Fbar", v)

		// If we want the raw unescaped version, we can use the component and read it
		_, component := SplitLast(m)
		require.Equal(t, "tmp/bar", string(component.RawValue()))
	})
}

func FuzzSplitRoundtrip(f *testing.F) {
	for _, v := range good {
		f.Add(v)
	}
	otherMultiaddr := StringCast("/udp/1337")

	f.Fuzz(func(t *testing.T, addrStr string) {
		addr, err := NewMultiaddr(addrStr)
		if err != nil {
			t.Skip() // Skip inputs that are not valid multiaddrs
		}

		// Test SplitFirst
		first, rest := SplitFirst(addr)
		joined := Join(first, rest)
		require.True(t, addr.Equal(joined), "SplitFirst and Join should round-trip")

		// Test SplitLast
		rest, last := SplitLast(addr)
		joined = Join(rest, last)
		require.True(t, addr.Equal(joined), "SplitLast and Join should round-trip")

		p := addr.Protocols()
		if len(p) == 0 {
			t.Skip()
		}

		tryPubMethods := func(a Multiaddr) {
			if a == nil {
				return
			}
			_ = a.Equal(otherMultiaddr)
			_ = a.Bytes()
			_ = a.String()
			_ = a.Protocols()
			_ = a.Encapsulate(otherMultiaddr)
			_ = a.Decapsulate(otherMultiaddr)
			_, _ = a.ValueForProtocol(P_TCP)
		}

		for _, proto := range p {
			splitFunc := func(c Component) bool {
				return c.Protocol().Code == proto.Code
			}
			beforeC, after := SplitFirst(addr)
			joined = Join(beforeC, after)
			require.True(t, addr.Equal(joined))
			tryPubMethods(after)

			before, afterC := SplitLast(addr)
			joined = Join(before, afterC)
			require.True(t, addr.Equal(joined))
			tryPubMethods(before)

			before, after = SplitFunc(addr, splitFunc)
			joined = Join(before, after)
			require.True(t, addr.Equal(joined))
			tryPubMethods(before)
			tryPubMethods(after)
		}
	})
}

func BenchmarkComponentValidation(b *testing.B) {
	comp, err := NewComponent("ip4", "127.0.0.1")
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := validateComponent(comp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func FuzzComponents(f *testing.F) {
	for _, v := range good {
		m := StringCast(v)
		for _, c := range m {
			f.Add(c.Bytes())
		}
	}
	f.Fuzz(func(t *testing.T, compBytes []byte) {
		n, c, err := readComponent(compBytes)
		if err != nil {
			t.Skip()
		}
		if c.protocol == nil {
			t.Fatal("component has nil protocol")
		}
		if c.protocol.Code == 0 {
			t.Fatal("component has nil protocol code")
		}
		if !bytes.Equal(c.Bytes(), compBytes[:n]) {
			t.Logf("component bytes: %v", c.Bytes())
			t.Logf("original bytes: %v", compBytes[:n])
			t.Fatal("component bytes are not equal to the original bytes")
		}
	})
}
