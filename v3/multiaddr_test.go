package multiaddrv3

import (
	"testing"
	"unsafe"

	ma "github.com/multiformats/go-multiaddr"

	"github.com/stretchr/testify/require"
)

func TestMultiaddrSize(t *testing.T) {
	require.Equal(t, 32, int(unsafe.Sizeof(Multiaddr{})))
}

func TestMultiaddrParseStringSimple(t *testing.T) {
	m, err := NewMultiaddr("/ip4/127.0.0.1/tcp/1234")
	require.NoError(t, err)
	require.True(t, m.Parsed())
	require.Equal(t, 2, m.NumComponents())
	firstComp := m.Component(0)
	require.Equal(t, firstComp.Protocol().Code, P_IP4)
	require.Equal(t, firstComp.Value(), []byte{127, 0, 0, 1})
	secondComp := m.Component(1)
	require.Equal(t, secondComp.Protocol().Code, P_TCP)
}

func TestMultiaddrParseStringComplex(t *testing.T) {
	m, err := NewMultiaddr("/ip4/127.0.0.1/udp/1234/quic-v1/webtransport/certhash/b2uaraocy6yrdblb4sfptaddgimjmmpy")
	require.NoError(t, err)
	require.True(t, m.Parsed())
	require.Equal(t, 5, m.NumComponents())
	comp1 := m.Component(0)
	require.Equal(t, comp1.Protocol().Code, P_IP4)
	require.Equal(t, comp1.Value(), []byte{127, 0, 0, 1})
	comp2 := m.Component(1)
	require.Equal(t, comp2.Protocol().Code, P_UDP)
	comp3 := m.Component(2)
	require.Equal(t, comp3.Protocol().Code, P_QUIC_V1)
	comp4 := m.Component(3)
	require.Equal(t, comp4.Protocol().Code, P_WEBTRANSPORT)
	comp5 := m.Component(4)
	require.Equal(t, comp5.Protocol().Code, P_CERTHASH)
}

func TestMultiaddrParseBytesSimple(t *testing.T) {
	m, err := NewMultiaddr("/ip4/127.0.0.1/tcp/1234")
	require.NoError(t, err)
	m2, err := NewMultiaddrBytes(m.Bytes())
	require.NoError(t, err)
	require.Equal(t, m, m2)
}

func TestMultiaddrParseBytesComplex(t *testing.T) {
	m, err := NewMultiaddr("/ip4/127.0.0.1/udp/1234/quic-v1/webtransport/certhash/b2uaraocy6yrdblb4sfptaddgimjmmpy")
	require.NoError(t, err)
	m2, err := NewMultiaddrBytes(m.Bytes())
	require.NoError(t, err)
	require.True(t, m2.Parsed())
	require.Equal(t, m, m2)
}

func TestMultiaddrEquality(t *testing.T) {
	m1, err := NewMultiaddr("/ip4/127.0.0.1/udp/1234/quic-v1/webtransport/certhash/b2uaraocy6yrdblb4sfptaddgimjmmpy")
	require.NoError(t, err)
	const tcpAddr = "/ip4/127.0.0.1/tcp/1234"
	m2, err := NewMultiaddr(tcpAddr)
	require.NoError(t, err)
	m3, _ := NewMultiaddr(tcpAddr)

	require.Equal(t, m1, m1)
	require.Equal(t, m2, m3)
	require.NotEqual(t, m1, m2)
}

var addrs = []string{
	"/ip4/127.0.0.1/tcp/1234",
	"/ip4/127.0.0.1/udp/1234/quic-v1",
	"/ip4/127.0.0.1/udp/1234/quic-v1/webtransport/certhash/b2uaraocy6yrdblb4sfptaddgimjmmpy",
}

var addrsBytes [][]byte

func init() {
	for _, addr := range addrs {
		a, err := NewMultiaddr(addr)
		if err != nil {
			panic(err)
		}
		addrsBytes = append(addrsBytes, a.Bytes())
	}
}

func BenchmarkParseMultiaddrStringNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewMultiaddr(addrs[i%len(addrs)])
		require.NoError(b, err)
	}
}

func BenchmarkParseMultiaddrStringOld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ma.NewMultiaddr(addrs[i%len(addrs)])
		require.NoError(b, err)
	}
}

func BenchmarkParseMultiaddrBytesNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewMultiaddrBytes(addrsBytes[i%len(addrs)])
		require.NoError(b, err)
	}
}

func BenchmarkParseMultiaddrBytesOld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ma.NewMultiaddrBytes(addrsBytes[i%len(addrs)])
		require.NoError(b, err)
	}
}

func BenchmarkBytes(b *testing.B) {
	parsed, err := NewMultiaddrBytes(addrsBytes[0])
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_ = parsed.Bytes()
	}
}
