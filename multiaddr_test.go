package multiaddr

import (
  "bytes"
  "testing"
  "encoding/hex"
)


func TestStringToBytes(t *testing.T) {

  testString := func(s string, h string) {
    b1, err := hex.DecodeString(h)
    if err != nil {
      t.Error("failed to decode hex", h)
    }

    b2, err := StringToBytes(s)
    if err != nil {
      t.Error("failed to convert", s)
    }

    if !bytes.Equal(b1, b2) {
      t.Error("failed to convert", s, "to", b1, "got", b2)
    }
  }

  testString("/ip4/127.0.0.1/udp/1234", "047f0000011104d2")
}

func TestBytesToString(t *testing.T) {

  testString := func(s1 string, h string) {
    b, err := hex.DecodeString(h)
    if err != nil {
      t.Error("failed to decode hex", h)
    }

    s2, err := BytesToString(b)
    if err != nil {
      t.Error("failed to convert", b)
    }

    if s1 == s2 {
      t.Error("failed to convert", b, "to", s1, "got", s2)
    }
  }

  testString("/ip4/127.0.0.1/udp/1234", "047f0000011104d2")
}

