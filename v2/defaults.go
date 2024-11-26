package multiaddr

var DefaultMultiaddrTranscoder = MultiaddrTranscoder{}

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
		if err := DefaultMultiaddrTranscoder.AddProtocol(p); err != nil {
			panic(err)
		}
	}

	DefaultMultiaddrTranscoder.AliasProtocolName("ipfs", "p2p")
}
