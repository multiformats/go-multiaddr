## Breaking changes in the large refactor of go-multiaddr v0.15

- There is no `Multiaddr` interface type.
- Multiaddr is now a concrete type. Not an interface.
- Empty Multiaddrs/Component should generally be checked with `.Empty()`, not `== nil`. This is similar to how slices should be checked with `len(s) == 0` rather than `s == nil`.
- Components do not implement `Multiaddr` as there is no `Multiaddr` to implement.
- `Multiaddr` can no longer be a key in a Map. If you want unique Multiaddrs, use `Multiaddr.String()` as the key, otherwise you can use the pointer value `*Multiaddr`.

## Callouts

- Multiaddr.Bytes() is a `O(n)` operation for n Components, as opposed to a `O(1)` operation.

## Migration tips for v0.15

- If your use case supports it, prefer `append` to append a Component to a
  Multiaddr rather than using `Encapsulate`. It's much faster as it does not do
  a defensive copy.
