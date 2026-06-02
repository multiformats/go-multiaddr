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

// portBtS used to call binary.BigEndian.Uint16 without checking len(b) >= 2.
// Direct callers of TranscoderPort (outside the codec.go length checks) could
// panic with `index out of range`. The internal guard and the new validator
// must turn both into normal errors.
func TestTranscoderPortRejectsShortInput(t *testing.T) {
	t.Run("BytesToString-short", func(t *testing.T) {
		if _, err := TranscoderPort.BytesToString(nil); err == nil {
			t.Fatalf("expected error on nil input")
		}
		if _, err := TranscoderPort.BytesToString([]byte{0x00}); err == nil {
			t.Fatalf("expected error on 1-byte input")
		}
	})
	t.Run("Validate-rejects-wrong-length", func(t *testing.T) {
		if err := TranscoderPort.ValidateBytes(nil); err == nil {
			t.Fatalf("expected validate error on nil input")
		}
		if err := TranscoderPort.ValidateBytes([]byte{0x00, 0x01, 0x02}); err == nil {
			t.Fatalf("expected validate error on 3-byte input")
		}
	})
	t.Run("Validate-accepts-2-byte", func(t *testing.T) {
		if err := TranscoderPort.ValidateBytes([]byte{0x00, 0x50}); err != nil {
			t.Fatalf("expected no error on 2-byte input, got %v", err)
		}
	})
}
