## Breaking changes in the large refactor of go-multiaddr v0.15

- There is no `Multiaddr` interface type.
- Multiaddr is now a concrete type. Not an interface.
- Empty Multiaddrs/ should be checked with `.Empty()`, not `== nil`
- Components do not implement `Multiaddr` as there is no `Multiaddr` to implement.
