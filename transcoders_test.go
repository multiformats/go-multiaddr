package multiaddr

import "testing"

// Regression test for https://github.com/multiformats/go-multiaddr/issues/288.
// onionBtS and onion3BtS used to index their input at fixed offsets without
// a length check, panicking with `slice bounds out of range` when called with
// short input (e.g. directly via the exported TranscoderOnion / TranscoderOnion3).
func TestOnionTranscodersRejectShortInput(t *testing.T) {
	cases := []struct {
		name string
		tr   Transcoder
		b    []byte
	}{
		{"onion-short", TranscoderOnion, []byte{0x00, 0x01}},
		{"onion-empty", TranscoderOnion, nil},
		{"onion3-short", TranscoderOnion3, []byte{0x00, 0x01}},
		{"onion3-empty", TranscoderOnion3, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.tr.BytesToString(tc.b); err == nil {
				t.Fatalf("expected error for short input, got nil")
			}
		})
	}
}
