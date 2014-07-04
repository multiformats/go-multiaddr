package multiaddr

import (
  "encoding/binary"
  "net"
  "strconv"
)

func AddressStringToBytes(p *Protocol, s string) []byte {
  switch p.Code {

  // ipv4,6
  case 4, 41:
    return net.ParseIP(s).To4()

  // tcp udp dccp sctp
  case 6, 17, 33, 132:
    b := make([]byte, 2)
    i, err := strconv.Atoi(s)
    if err == nil {
      binary.BigEndian.PutUint16(b, uint16(i))
    }
    return b
  }

  return []byte{}
}

func AddressBytesToString(p *Protocol, b []byte) string {
  switch p.Code {

  // ipv4,6
  case 4, 41:
    return net.IP(b).String()

  // tcp udp dccp sctp
  case 6, 17, 33, 132:
    i := binary.BigEndian.Uint16(b)
    return strconv.Itoa(int(i))
  }

  return ""
}
