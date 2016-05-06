package manet

import (
	"net"
	"testing"

	ma "github.com/jbenet/go-multiaddr"
)

func TestRegisterSpec(t *testing.T) {
	myproto := &AddressSpec{
		Key:              "test",
		NetNames:         []string{"test", "iptest", "blahtest"},
		ConvertMultiaddr: func(a ma.Multiaddr) (net.Addr, error) { return nil, nil },
		ParseNetAddr:     func(a net.Addr) (ma.Multiaddr, error) { return nil, nil },
	}

	RegisterAddressType(myproto)

	_, ok := addrParsers["test"]
	if !ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = addrParsers["iptest"]
	if !ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = addrParsers["blahtest"]
	if !ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = maddrParsers["test"]
	if !ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = maddrParsers["iptest"]
	if ok {
		t.Fatal("myproto not properly registered")
	}

	_, ok = maddrParsers["blahtest"]
	if ok {
		t.Fatal("myproto not properly registered")
	}

}
