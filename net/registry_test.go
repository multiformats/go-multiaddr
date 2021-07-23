package manet

import (
	"net"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

func TestRegisterFrom(t *testing.T) {
	cm := NewCodecMap()
	cm.RegisterFromNetAddr(
		func(a net.Addr) (ma.Multiaddr, error) { return nil, nil },
		"test", "iptest", "blahtest",
	)
	if _, ok := cm.addrParsers["test"]; !ok {
		t.Fatal("myproto not properly registered")
	}
	if _, ok := cm.addrParsers["iptest"]; !ok {
		t.Fatal("myproto not properly registered")
	}
	if _, ok := cm.addrParsers["blahtest"]; !ok {
		t.Fatal("myproto not properly registered")
	}
}

func TestRegisterTo(t *testing.T) {
	cm := NewCodecMap()
	cm.RegisterToNetAddr(
		func(a ma.Multiaddr) (net.Addr, error) { return nil, nil },
		"test", "iptest", "blahtest",
	)
	if _, ok := cm.maddrParsers["test"]; !ok {
		t.Fatal("myproto not properly registered")
	}
	if _, ok := cm.maddrParsers["iptest"]; !ok {
		t.Fatal("myproto not properly registered")
	}
	if _, ok := cm.maddrParsers["blahtest"]; !ok {
		t.Fatal("myproto not properly registered")
	}
}
