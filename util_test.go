package multiaddr

import (
	"strings"
	"testing"
)

func TestSplitFirstLast(t *testing.T) {
	ipStr := "/ip4/0.0.0.0"
	tcpStr := "/tcp/123"
	quicStr := "/quic"
	ipfsStr := "/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7"

	for _, x := range [][]string{
		{ipStr, tcpStr, quicStr, ipfsStr},
		{ipStr, tcpStr, ipfsStr},
		{ipStr, tcpStr},
		{ipStr},
	} {
		addr := StringCast(strings.Join(x, ""))
		head, tail := SplitFirst(addr)
		rest, last := SplitLast(addr)
		if len(x) == 0 {
			if !head.Empty() {
				t.Error("expected head to be empty")
			}
			if tail != nil {
				t.Error("expected tail to be nil")
			}
			if rest != nil {
				t.Error("expected rest to be nil")
			}
			if !last.Empty() {
				t.Error("expected last to be empty")
			}
			continue
		}
		if !head.AsMultiaddr().Equal(StringCast(x[0])) {
			t.Errorf("expected %s to be %s", head, x[0])
		}
		if !last.AsMultiaddr().Equal(StringCast(x[len(x)-1])) {
			t.Errorf("expected %s to be %s", head, x[len(x)-1])
		}
		if len(x) == 1 {
			if tail != nil {
				t.Error("expected tail to be nil")
			}
			if rest != nil {
				t.Error("expected rest to be nil")
			}
			continue
		}
		tailExp := strings.Join(x[1:], "")
		if !tail.Equal(StringCast(tailExp)) {
			t.Errorf("expected %s to be %s", tail, tailExp)
		}
		restExp := strings.Join(x[:len(x)-1], "")
		if !rest.Equal(StringCast(restExp)) {
			t.Errorf("expected %s to be %s", rest, restExp)
		}
	}

	c, err := NewComponent("ip4", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	ci, m := SplitFirst(c.AsMultiaddr())
	if !ci.Equal(c) || m != nil {
		t.Error("split first on component failed")
	}
	m, ci = SplitLast(c.AsMultiaddr())
	if !ci.Equal(c) || m != nil {
		t.Error("split last on component failed")
	}
	cis := Split(c.AsMultiaddr())
	if len(cis) != 1 || !cis[0].Equal(c) {
		t.Error("split on component failed")
	}
	m1, m2 := SplitFunc(c.AsMultiaddr(), func(c Component) bool {
		return true
	})
	if m1 != nil || !m2.Equal(c.AsMultiaddr()) {
		t.Error("split func(true) on component failed")
	}
	m1, m2 = SplitFunc(c.AsMultiaddr(), func(c Component) bool {
		return false
	})
	if !m1.Equal(c.AsMultiaddr()) || m2 != nil {
		t.Error("split func(false) on component failed")
	}

	i := 0
	ForEach(c.AsMultiaddr(), func(ci Component) bool {
		if i != 0 {
			t.Error("expected exactly one component")
		}
		i++
		if !ci.Equal(c) {
			t.Error("foreach on component failed")
		}
		return true
	})
}

func TestSplitFunc(t *testing.T) {
	ipStr := "/ip4/0.0.0.0"
	tcpStr := "/tcp/123"
	quicStr := "/quic"
	ipfsStr := "/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7"

	for _, x := range [][]string{
		{ipStr, tcpStr, quicStr, ipfsStr},
		{ipStr, tcpStr, ipfsStr},
		{ipStr, tcpStr},
		{ipStr},
	} {
		addr := StringCast(strings.Join(x, ""))
		for i, cs := range x {
			target := StringCast(cs)
			a, b := SplitFunc(addr, func(c Component) bool {
				return c.AsMultiaddr().Equal(target)
			})
			if i == 0 {
				if !a.Empty() {
					t.Error("expected nil addr")
				}
			} else {
				if !a.Equal(StringCast(strings.Join(x[:i], ""))) {
					t.Error("split failed")
				}
				if !b.Equal(StringCast(strings.Join(x[i:], ""))) {
					t.Error("split failed")
				}
			}
		}
		a, b := SplitFunc(addr, func(_ Component) bool { return false })
		if !a.Equal(addr) || !b.Empty() {
			t.Error("should not have split")
		}
	}
}
