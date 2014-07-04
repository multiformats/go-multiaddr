package multiaddr

type Multiaddr struct {
  Bytes []byte
}

func NewString(s string) *Multiaddr {
  m := &Multiaddr{}
  return m
}
