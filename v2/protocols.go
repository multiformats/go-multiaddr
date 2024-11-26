package multiaddr

import (
	"encoding/binary"
	"math"
)

// You **MUST** register your multicodecs with
// https://github.com/multiformats/multicodec before adding them here.
const (
	P_IP4               = 4
	P_TCP               = 6
	P_DNS               = 53 // 4 or 6
	P_DNS4              = 54
	P_DNS6              = 55
	P_DNSADDR           = 56
	P_UDP               = 273
	P_DCCP              = 33
	P_IP6               = 41
	P_IP6ZONE           = 42
	P_IPCIDR            = 43
	P_QUIC              = 460
	P_QUIC_V1           = 461
	P_WEBTRANSPORT      = 465
	P_CERTHASH          = 466
	P_SCTP              = 132
	P_CIRCUIT           = 290
	P_UDT               = 301
	P_UTP               = 302
	P_UNIX              = 400
	P_P2P               = 421
	P_IPFS              = P_P2P // alias for backwards compatibility
	P_HTTP              = 480
	P_HTTP_PATH         = 481
	P_HTTPS             = 443 // deprecated alias for /tls/http
	P_ONION             = 444 // also for backwards compatibility
	P_ONION3            = 445
	P_GARLIC64          = 446
	P_GARLIC32          = 447
	P_P2P_WEBRTC_DIRECT = 276 // Deprecated. use webrtc-direct instead
	P_TLS               = 448
	P_SNI               = 449
	P_NOISE             = 454
	P_WS                = 477
	P_WSS               = 478 // deprecated alias for /tls/ws
	P_PLAINTEXTV2       = 7367777
	P_WEBRTC_DIRECT     = 280
	P_WEBRTC            = 281
	P_MEMORY            = 777
)

func codeToVarint(num int) []byte {
	if num < 0 || num > math.MaxInt32 {
		panic("invalid code")
	}
	return binary.AppendUvarint(nil, uint64(num))
}

var (
	protoIP4 = Protocol{
		Name:       "ip4",
		Code:       P_IP4,
		VCode:      codeToVarint(P_IP4),
		Size:       32,
		Path:       false,
		Transcoder: TranscoderIP4,
	}
	protoTCP = Protocol{
		Name:       "tcp",
		Code:       P_TCP,
		VCode:      codeToVarint(P_TCP),
		Size:       16,
		Path:       false,
		Transcoder: TranscoderPort,
	}
	protoDNS = Protocol{
		Code:       P_DNS,
		Size:       LengthPrefixedVarSize,
		Name:       "dns",
		VCode:      codeToVarint(P_DNS),
		Transcoder: TranscoderDns,
	}
	protoDNS4 = Protocol{
		Code:       P_DNS4,
		Size:       LengthPrefixedVarSize,
		Name:       "dns4",
		VCode:      codeToVarint(P_DNS4),
		Transcoder: TranscoderDns,
	}
	protoDNS6 = Protocol{
		Code:       P_DNS6,
		Size:       LengthPrefixedVarSize,
		Name:       "dns6",
		VCode:      codeToVarint(P_DNS6),
		Transcoder: TranscoderDns,
	}
	protoDNSADDR = Protocol{
		Code:       P_DNSADDR,
		Size:       LengthPrefixedVarSize,
		Name:       "dnsaddr",
		VCode:      codeToVarint(P_DNSADDR),
		Transcoder: TranscoderDns,
	}
	protoUDP = Protocol{
		Name:       "udp",
		Code:       P_UDP,
		VCode:      codeToVarint(P_UDP),
		Size:       16,
		Path:       false,
		Transcoder: TranscoderPort,
	}
	protoDCCP = Protocol{
		Name:       "dccp",
		Code:       P_DCCP,
		VCode:      codeToVarint(P_DCCP),
		Size:       16,
		Path:       false,
		Transcoder: TranscoderPort,
	}
	protoIP6 = Protocol{
		Name:       "ip6",
		Code:       P_IP6,
		VCode:      codeToVarint(P_IP6),
		Size:       128,
		Transcoder: TranscoderIP6,
	}
	protoIPCIDR = Protocol{
		Name:       "ipcidr",
		Code:       P_IPCIDR,
		VCode:      codeToVarint(P_IPCIDR),
		Size:       8,
		Transcoder: TranscoderIPCIDR,
	}
	// these require varint
	protoIP6ZONE = Protocol{
		Name:       "ip6zone",
		Code:       P_IP6ZONE,
		VCode:      codeToVarint(P_IP6ZONE),
		Size:       LengthPrefixedVarSize,
		Path:       false,
		Transcoder: TranscoderIP6Zone,
	}
	protoSCTP = Protocol{
		Name:       "sctp",
		Code:       P_SCTP,
		VCode:      codeToVarint(P_SCTP),
		Size:       16,
		Transcoder: TranscoderPort,
	}

	protoCIRCUIT = Protocol{
		Code:  P_CIRCUIT,
		Size:  0,
		Name:  "p2p-circuit",
		VCode: codeToVarint(P_CIRCUIT),
	}

	protoONION2 = Protocol{
		Name:       "onion",
		Code:       P_ONION,
		VCode:      codeToVarint(P_ONION),
		Size:       96,
		Transcoder: TranscoderOnion,
	}
	protoONION3 = Protocol{
		Name:       "onion3",
		Code:       P_ONION3,
		VCode:      codeToVarint(P_ONION3),
		Size:       296,
		Transcoder: TranscoderOnion3,
	}
	protoGARLIC64 = Protocol{
		Name:       "garlic64",
		Code:       P_GARLIC64,
		VCode:      codeToVarint(P_GARLIC64),
		Size:       LengthPrefixedVarSize,
		Transcoder: TranscoderGarlic64,
	}
	protoGARLIC32 = Protocol{
		Name:       "garlic32",
		Code:       P_GARLIC32,
		VCode:      codeToVarint(P_GARLIC32),
		Size:       LengthPrefixedVarSize,
		Transcoder: TranscoderGarlic32,
	}
	protoUTP = Protocol{
		Name:  "utp",
		Code:  P_UTP,
		VCode: codeToVarint(P_UTP),
	}
	protoUDT = Protocol{
		Name:  "udt",
		Code:  P_UDT,
		VCode: codeToVarint(P_UDT),
	}
	protoQUIC = Protocol{
		Name:  "quic",
		Code:  P_QUIC,
		VCode: codeToVarint(P_QUIC),
	}
	protoQUICV1 = Protocol{
		Name:  "quic-v1",
		Code:  P_QUIC_V1,
		VCode: codeToVarint(P_QUIC_V1),
	}
	protoWEBTRANSPORT = Protocol{
		Name:  "webtransport",
		Code:  P_WEBTRANSPORT,
		VCode: codeToVarint(P_WEBTRANSPORT),
	}
	protoCERTHASH = Protocol{
		Name:       "certhash",
		Code:       P_CERTHASH,
		VCode:      codeToVarint(P_CERTHASH),
		Size:       LengthPrefixedVarSize,
		Transcoder: TranscoderCertHash,
	}
	protoHTTP = Protocol{
		Name:  "http",
		Code:  P_HTTP,
		VCode: codeToVarint(P_HTTP),
	}
	protoHTTPPath = Protocol{
		Name:       "http-path",
		Code:       P_HTTP_PATH,
		VCode:      codeToVarint(P_HTTP_PATH),
		Size:       LengthPrefixedVarSize,
		Transcoder: TranscoderHTTPPath,
	}
	protoHTTPS = Protocol{
		Name:  "https",
		Code:  P_HTTPS,
		VCode: codeToVarint(P_HTTPS),
	}
	protoP2P = Protocol{
		Name:       "p2p",
		Code:       P_P2P,
		VCode:      codeToVarint(P_P2P),
		Size:       LengthPrefixedVarSize,
		Transcoder: TranscoderP2P,
	}
	protoUNIX = Protocol{
		Name:       "unix",
		Code:       P_UNIX,
		VCode:      codeToVarint(P_UNIX),
		Size:       LengthPrefixedVarSize,
		Path:       true,
		Transcoder: TranscoderUnix,
	}
	protoP2P_WEBRTC_DIRECT = Protocol{
		Name:  "p2p-webrtc-direct",
		Code:  P_P2P_WEBRTC_DIRECT,
		VCode: codeToVarint(P_P2P_WEBRTC_DIRECT),
	}
	protoTLS = Protocol{
		Name:  "tls",
		Code:  P_TLS,
		VCode: codeToVarint(P_TLS),
	}
	protoSNI = Protocol{
		Name:       "sni",
		Size:       LengthPrefixedVarSize,
		Code:       P_SNI,
		VCode:      codeToVarint(P_SNI),
		Transcoder: TranscoderDns,
	}
	protoNOISE = Protocol{
		Name:  "noise",
		Code:  P_NOISE,
		VCode: codeToVarint(P_NOISE),
	}
	protoPlaintextV2 = Protocol{
		Name:  "plaintextv2",
		Code:  P_PLAINTEXTV2,
		VCode: codeToVarint(P_PLAINTEXTV2),
	}
	protoWS = Protocol{
		Name:  "ws",
		Code:  P_WS,
		VCode: codeToVarint(P_WS),
	}
	protoWSS = Protocol{
		Name:  "wss",
		Code:  P_WSS,
		VCode: codeToVarint(P_WSS),
	}
	protoWebRTCDirect = Protocol{
		Name:  "webrtc-direct",
		Code:  P_WEBRTC_DIRECT,
		VCode: codeToVarint(P_WEBRTC_DIRECT),
	}
	protoWebRTC = Protocol{
		Name:  "webrtc",
		Code:  P_WEBRTC,
		VCode: codeToVarint(P_WEBRTC),
	}

	protoMemory = Protocol{
		Name:       "memory",
		Code:       P_MEMORY,
		VCode:      codeToVarint(P_MEMORY),
		Size:       64,
		Transcoder: TranscoderMemory,
	}
)

func init() {
	for _, p := range []Protocol{
		protoIP4,
		protoTCP,
		protoDNS,
		protoDNS4,
		protoDNS6,
		protoDNSADDR,
		protoUDP,
		protoDCCP,
		protoIP6,
		protoIP6ZONE,
		protoIPCIDR,
		protoSCTP,
		protoCIRCUIT,
		protoONION2,
		protoONION3,
		protoGARLIC64,
		// protoGARLIC32, // Disabling for now because /garlic32/566niximlxdzpanmn4qouucvua3k7neniwss47li5r6ugoertzuqzwassw does not round trip
		protoUTP,
		protoUDT,
		protoQUIC,
		protoQUICV1,
		protoWEBTRANSPORT,
		protoCERTHASH,
		protoHTTP,
		protoHTTPPath,
		protoHTTPS,
		protoP2P,
		// protoUNIX, // disabling while debugging /unix/stdio
		protoP2P_WEBRTC_DIRECT,
		protoTLS,
		protoSNI,
		protoNOISE,
		protoWS,
		protoWSS,
		protoPlaintextV2,
		protoWebRTCDirect,
		protoWebRTC,
		protoMemory,
	} {
		if err := AddProtocol(p); err != nil {
			panic(err)
		}
	}

	// explicitly set both of these
	protocolsByName["p2p"] = protoP2P
	protocolsByName["ipfs"] = protoP2P
}
