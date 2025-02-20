module github.com/multiformats/go-multiaddr/compattest

go 1.23.4

replace github.com/multiformats/go-multiaddr => ../

replace github.com/multiformats/go-multiaddr/compattest/internal/prev => ./internal/prev

require (
	github.com/multiformats/go-multiaddr v0.0.0-00010101000000-000000000000
	github.com/multiformats/go-multiaddr/compattest/internal/prev v0.0.0-00010101000000-000000000000
)

require (
	github.com/ipfs/go-cid v0.0.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multihash v0.2.3 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/exp v0.0.0-20230725012225-302865e7556b // indirect
	golang.org/x/sys v0.28.0 // indirect
	lukechampine.com/blake3 v1.2.1 // indirect
)
