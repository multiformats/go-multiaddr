package multiaddr

import (
  "fmt"
)

type Multiaddr struct {
  Bytes []byte
}

func NewString(s string) (*Multiaddr, error) {
  b, err := StringToBytes(s)
  if err != nil {
    return nil, err
  }
  return &Multiaddr{Bytes: b}, nil
}

func (m *Multiaddr) String() (string, error) {
  return BytesToString(m.Bytes)
}

func (m *Multiaddr) Protocols() (ret []*Protocol, err error) {

  // panic handler, in case we try accessing bytes incorrectly.
  defer func() {
    if e := recover(); e != nil {
      ret = nil
      err = e.(error)
    }
  }()

  ps := []*Protocol{}
  b := m.Bytes[:]
  for ; len(b) > 0 ; {
    p := ProtocolWithCode(int(b[0]))
    if p == nil {
      return nil, fmt.Errorf("no protocol with code %d", b[0])
    }
    ps = append(ps, p)
    b = b[1 + (p.Size / 8):]
  }
  return ps, nil
}
