package multiaddr

import "testing"

func expectVarint(t *testing.T, x, expected int) {
	size := VarintSize(x)
	if size != expected {
		t.Fatalf("expected varintsize of %d to be %d, got %d", x, expected, size)
	}
}

func TestVarintSize(t *testing.T) {
	expectVarint(t, (1<<7)-1, 1)
	expectVarint(t, 0, 1)
	expectVarint(t, 1<<7, 2)
}
