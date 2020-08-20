package codependencies

import (
	// Packages imported into this package.

	// go-multiaddr-net < 0.2.0 conflict with this package.
	_ "github.com/multiformats/go-multiaddr-net"
	// go-maddr-filter < 0.1.0 conflicts with this package.
	_ "github.com/libp2p/go-maddr-filter"
)
