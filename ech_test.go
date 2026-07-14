package multiaddr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestECHProtocol(t *testing.T) {
	echConfigList := []byte{0x00, 0x04, 0xff, 0xff, 0x00, 0x00}
	const encodedECHConfigList = "uAAT__wAA"

	protocol := ProtocolWithName("ech")
	require.Equal(t, P_ECH, protocol.Code)
	require.Equal(t, []byte{0xf9, 0x4c}, protocol.VCode)
	require.Equal(t, LengthPrefixedVarSize, protocol.Size)
	require.NotNil(t, protocol.Transcoder)
	require.Equal(t, protocol, ProtocolWithCode(P_ECH))

	component, err := NewComponent("ech", encodedECHConfigList)
	require.NoError(t, err)
	require.Equal(t, "/ech/"+encodedECHConfigList, component.String())
	require.Equal(t, echConfigList, component.RawValue())
	require.Equal(t, append([]byte{0xf9, 0x4c, 0x06}, echConfigList...), component.Bytes())

	addr, err := NewMultiaddr("/ip4/192.0.2.1/tcp/443/tls/ech/" + encodedECHConfigList)
	require.NoError(t, err)
	require.Equal(t, "/ip4/192.0.2.1/tcp/443/tls/ech/"+encodedECHConfigList, addr.String())

	value, err := addr.ValueForProtocol(P_ECH)
	require.NoError(t, err)
	require.Equal(t, encodedECHConfigList, value)

	fromBytes, err := NewMultiaddrBytes(addr.Bytes())
	require.NoError(t, err)
	require.True(t, addr.Equal(fromBytes))
}

func TestECHAcceptsAnyMultibaseEncoding(t *testing.T) {
	addr, err := NewMultiaddr("/quic-v1/ech/baacp77yaaa")
	require.NoError(t, err)
	require.Equal(t, "/quic-v1/ech/uAAT__wAA", addr.String())
}

func TestECHRejectsInvalidValues(t *testing.T) {
	for _, value := range []string{"not-multibase", "u"} {
		_, err := NewComponent("ech", value)
		require.Error(t, err)
	}
}
