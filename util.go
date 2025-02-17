package multiaddr

import (
	"fmt"
)

// Split returns the sub-address portions of a multiaddr.
func Split(m Multiaddr) []Component {
	return m
}

func JoinComponents(cs ...Component) Multiaddr {
	if len(cs) == 0 {
		return nil
	}
	out := make([]Component, 0, len(cs))
	for _, c := range cs {
		if !c.Empty() {
			out = append(out, c)
		}
	}
	return out
}

// Join returns a combination of addresses.
// Note: This copies all the components from the input Multiaddrs. Depending on
// your use case, you may prefer to use `append(leftMA, rightMA...)` instead.
func Join(ms ...Multiaddr) Multiaddr {
	size := 0
	for _, m := range ms {
		size += len(m)
	}
	if size == 0 {
		return nil
	}

	out := make([]Component, 0, size)
	for _, m := range ms {
		for _, c := range m {
			if !c.Empty() {
				out = append(out, c)
			}
		}
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
		return Component{}, nil
	}
	if len(m) == 1 {
		return m[0], nil
	}
	return m[0], m[1:]
}

// SplitLast returns the rest of the multiaddr and the last component.
func SplitLast(m Multiaddr) (Multiaddr, Component) {
	if m.Empty() {
		return nil, Component{}
	}
	if len(m) == 1 {
		// We want to explicitly return a nil slice if the prefix is now empty.
		return nil, m[0]
	}
	return m[:len(m)-1], m[len(m)-1]
}

// SplitFunc splits the multiaddr when the callback first returns true. The
// component on which the callback first returns will be included in the
// *second* multiaddr.
func SplitFunc(m Multiaddr, cb func(Component) bool) (Multiaddr, Multiaddr) {
	if m.Empty() {
		return nil, nil
	}

	idx := len(m)
	for i, c := range m {
		if cb(c) {
			idx = i
			break
		}
	}
	pre, post := m[:idx], m[idx:]
	if pre.Empty() {
		pre = nil
	}
	if post.Empty() {
		post = nil
	}
	return pre, post
}

// ForEach walks over the multiaddr, component by component.
//
// This function iterates over components.
// Return true to continue iteration, false to stop.
//
// Prefer to use a standard for range loop instead
// e.g. `for _, c := range m { ... }`
func ForEach(m Multiaddr, cb func(c Component) bool) {
	for _, c := range m {
		if !cb(c) {
			return
		}
	}
}
