package multiaddr

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
