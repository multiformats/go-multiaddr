# go-multiaddr

> [multiaddr](https://github.com/multiformats/multiaddr) implementation in go

Multiaddr is a standard way to represent addresses that:

- Support any standard network protocols.
- Self-describe (include protocols).
- Have a binary packed format.
- Have a nice string representation.
- Encapsulate well.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
  - [Example](#example)
    - [Simple](#simple)
    - [Protocols](#protocols)
    - [En/decapsulate](#endecapsulate)
    - [Tunneling](#tunneling)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

```sh
go get github.com/TRON-US/go-multiaddr
```

## Usage

### Example

#### Simple

```go
import ma "github.com/TRON-US/go-multiaddr"

// construct from a string (err signals parse failure)
m1, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/1234")

// construct from bytes (err signals parse failure)
m2, err := ma.NewMultiaddrBytes(m1.Bytes())

// true
strings.Equal(m1.String(), "/ip4/127.0.0.1/udp/1234")
strings.Equal(m1.String(), m2.String())
bytes.Equal(m1.Bytes(), m2.Bytes())
m1.Equal(m2)
m2.Equal(m1)
```

#### Protocols

```go
// get the multiaddr protocol description objects
m1.Protocols()
// []Protocol{
//   Protocol{ Code: 4, Name: 'ip4', Size: 32},
//   Protocol{ Code: 17, Name: 'udp', Size: 16},
// }
```

#### En/decapsulate

```go
import ma "github.com/TRON-US/go-multiaddr"

m, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/1234")
// <Multiaddr /ip4/127.0.0.1/udp/1234>

sctpMA, err := ma.NewMultiaddr("/sctp/5678")

m.Encapsulate(sctpMA)
// <Multiaddr /ip4/127.0.0.1/udp/1234/sctp/5678>

udpMA, err := ma.NewMultiaddr("/udp/1234")

m.Decapsulate(udpMA) // up to + inc last occurrence of subaddr
// <Multiaddr /ip4/127.0.0.1>
```

#### Tunneling

Multiaddr allows expressing tunnels very nicely.

```js
printer, _ := ma.NewMultiaddr("/ip4/192.168.0.13/tcp/80")
proxy, _ := ma.NewMultiaddr("/ip4/10.20.30.40/tcp/443")
printerOverProxy := proxy.Encapsulate(printer)
// /ip4/10.20.30.40/tcp/443/ip4/192.168.0.13/tcp/80

proxyAgain := printerOverProxy.Decapsulate(printer)
// /ip4/10.20.30.40/tcp/443
```

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/TRON-US/go-multiaddr/issues).

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

[MIT](LICENSE) © TRON-US. 
