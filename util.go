package multiaddr

import (
	"fmt"
)

// Split returns the sub-address portions of a multiaddr.
func Split(m Multiaddr) []Component {
	return m
}

func JoinComponents(cs ...Component) Multiaddr {
	return cs
}

// Join returns a combination of addresses.
func Join(ms ...Multiaddr) Multiaddr {
	size := 0
	for _, m := range ms {
		size += len(m)
	}

	out := make([]Component, 0, size)
	for _, m := range ms {
		out = append(out, m...)
	}
	return out
}

// Cast re-casts a byte slice as a multiaddr. will panic if it fails to parse.
func Cast(b []byte) Multiaddr {
	m, err := NewMultiaddrBytes(b)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return m
}

// StringCast like Cast, but parses a string. Will also panic if it fails to parse.
func StringCast(s string) Multiaddr {
	m, err := NewMultiaddr(s)
	if err != nil {
		panic(fmt.Errorf("multiaddr failed to parse: %s", err))
	}
	return m
}

// SplitFirst returns the first component and the rest of the multiaddr.
func SplitFirst(m Multiaddr) (Component, Multiaddr) {
	if m.Empty() {
		return Component{}, Multiaddr{}
	}
	return m[0], m[1:]
}

// SplitLast returns the rest of the multiaddr and the last component.
func SplitLast(m Multiaddr) (Multiaddr, Component) {
	if m.Empty() {
		return Multiaddr{}, Component{}
	}
	return m[:len(m)-1], m[len(m)-1]
}

// SplitFunc splits the multiaddr when the callback first returns true. The
// component on which the callback first returns will be included in the
// *second* multiaddr.
func SplitFunc(m Multiaddr, cb func(Component) bool) (Multiaddr, Multiaddr) {
	if m.Empty() {
		return Multiaddr{}, Multiaddr{}
	}

	idx := len(m)
	for i, c := range m {
		if cb(c) {
			idx = i
			break
		}
	}
	return m[:idx], m[idx:]
}

// ForEach walks over the multiaddr, component by component.
//
// Deprecated: use a simple `for _, c := range m` instead.
//
// This function iterates over components.
// Return true to continue iteration, false to stop.
func ForEach(m Multiaddr, cb func(c Component) bool) {
	for _, c := range m {
		if !cb(c) {
			return
		}
	}
}
