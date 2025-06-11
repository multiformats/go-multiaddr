//go:build !tinygo

package manet

import (
	"fmt"
	"net"

	ma "github.com/multiformats/go-multiaddr"
)

func parseBasicNetMaddr(maddr ma.Multiaddr) (net.Addr, error) {
	network, host, err := DialArgs(maddr)
	if err != nil {
		return nil, err
	}

	switch network {
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, host)
	case "udp", "udp4", "udp6":
		return net.ResolveUDPAddr(network, host)
	case "ip", "ip4", "ip6":
		return net.ResolveIPAddr(network, host)
	case "unix":
		return net.ResolveUnixAddr(network, host)
	}

	return nil, fmt.Errorf("network not supported: %s", network)
}

func wrap(nconn net.Conn, laddr, raddr ma.Multiaddr) Conn {
	endpts := maEndpoints{
		laddr: laddr,
		raddr: raddr,
	}
	// This sucks. However, it's the only way to reliably expose the
	// underlying methods. This way, users that need access to, e.g.,
	// CloseRead and CloseWrite, can do so via type assertions.
	switch nconn := nconn.(type) {
	case *net.TCPConn:
		return &struct {
			*net.TCPConn
			maEndpoints
		}{nconn, endpts}
	case *net.UDPConn:
		return &struct {
			*net.UDPConn
			maEndpoints
		}{nconn, endpts}
	case *net.IPConn:
		return &struct {
			*net.IPConn
			maEndpoints
		}{nconn, endpts}
	case *net.UnixConn:
		return &struct {
			*net.UnixConn
			maEndpoints
		}{nconn, endpts}
	case halfOpen:
		return &struct {
			halfOpen
			maEndpoints
		}{nconn, endpts}
	default:
		return &struct {
			net.Conn
			maEndpoints
		}{nconn, endpts}
	}
}

func listenPacket(network, address string) (net.PacketConn, error) {
	return net.ListenPacket(network, address)
}
