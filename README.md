# go-multiaddr-net

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://ipn.io)
[![](https://img.shields.io/badge/project-multiformats-blue.svg?style=flat-square)](http://github.com/multiformats/multiformats)
[![](https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23ipfs)

> multiaddr net tools

This package provides [Multiaddr](http://github.com/multiformats/go-multiaddr) specific versions of common functions in [stdlib](https://github.com/golang/go/tree/master/src)'s `net` package. This means wrappers of standard net symbols like `net.Dial` and `net.Listen`, as well
as conversion to and from `net.Addr`.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

```sh
go get github.com/multiformats/go-multiaddr-net
```

## Usage

See the docs:

- `multiaddr/net`: https://godoc.org/github.com/multiformats/go-multiaddr-net
- `multiaddr`: https://godoc.org/github.com/multiformats/go-multiaddr

## Maintainers

Captain: [@whyrusleeping](https://github.com/whyrusleeping).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/multiformats/go-multiaddr-net/issues).

Check out our [contributing document](https://github.com/multiformats/multiformats/blob/master/contributing.md) for more information on how we work, and about contributing in general. Please be aware that all interactions related to multiformats are subject to the IPFS [Code of Conduct](https://github.com/ipfs/community/blob/master/code-of-conduct.md).

Small note: If editing the Readme, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

[MIT](LICENSE) © Juan Batiz-Benet
