package manet

import (
	"fmt"
	"net"
	"sync"

	ma "github.com/jbenet/go-multiaddr"
)

type AddrParser func(a net.Addr) (ma.Multiaddr, error)
type MaddrParser func(ma ma.Multiaddr) (net.Addr, error)

var maddrParsers map[string]MaddrParser
var addrParsers map[string]AddrParser
var addrParsersLock sync.Mutex

func RegisterAddressType(netname, maname string, ap AddrParser, mp MaddrParser) {
	addrParsersLock.Lock()
	defer addrParsersLock.Unlock()
	addrParsers[netname] = ap
	maddrParsers[maname] = mp
}

func init() {
	addrParsers = make(map[string]AddrParser)
	maddrParsers = make(map[string]MaddrParser)

	funcs := map[string]AddrParser{
		"tcp": parseTcpNetAddr,
		"udp": parseUdpNetAddr,
		"utp": parseUtpNetAddr,
	}

	for k, v := range funcs {
		RegisterAddressType(k, k, v, parseBasicNetMaddr)
		RegisterAddressType(k+"4", k, v, parseBasicNetMaddr)
		RegisterAddressType(k+"6", k, v, parseBasicNetMaddr)
	}

	for _, i := range []string{"ip", "ip4", "ip6"} {
		RegisterAddressType(i, i, parseIpNetAddr, parseBasicNetMaddr)
	}

	addrParsers["ip+net"] = parseIpPlusNetAddr
}

func getAddrParser(net string) (AddrParser, error) {
	addrParsersLock.Lock()
	defer addrParsersLock.Unlock()

	parser, ok := addrParsers[net]
	if !ok {
		return nil, fmt.Errorf("unknown network %v", net)
	}
	return parser, nil
}

func getMaddrParser(name string) (MaddrParser, error) {
	addrParsersLock.Lock()
	defer addrParsersLock.Unlock()
	p, ok := maddrParsers[name]
	if !ok {
		return nil, fmt.Errorf("network not supported: %s", name)
	}

	return p, nil
}
