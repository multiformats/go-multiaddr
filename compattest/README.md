# Compatibility Test

This package is used to test the backwards compatibility of the `go-multiaddr`
package against a previous version.

To update the previous version, from the root of the repo:
```sh
git subtree pull --prefix=compattest/internal/prev https://github.com/multiformats/go-multiaddr.git <tag> --squash
```

