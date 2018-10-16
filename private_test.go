package manet

import (
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

func TestIsPublicAddr(t *testing.T) {
	a, err := ma.NewMultiaddr("/ip4/192.168.1.1/tcp/80")
	if err != nil {
		t.Fatal(err)
	}

	if IsPublicAddr(a) {
		t.Fatal("192.168.1.1 is not a public address!")
	}

	a, err = ma.NewMultiaddr("/ip4/1.1.1.1/tcp/80")
	if err != nil {
		t.Fatal(err)
	}

	if !IsPublicAddr(a) {
		t.Fatal("1.1.1.1 is a public address!")
	}
}
