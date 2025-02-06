// Package matest provides utilities for testing with multiaddrs.
package matest

import (
	"slices"

	"github.com/multiformats/go-multiaddr"
)

type TestingT interface {
	Errorf(format string, args ...interface{})
}

type tHelper interface {
	Helper()
}

func AssertEqualMultiaddr(t TestingT, expected, actual multiaddr.Multiaddr) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !expected.Equal(actual) {
		t.Errorf("expected %v, got %v", expected, actual)
		return false
	}
	return true
}

func AssertEqualMultiaddrs(t TestingT, expected, actual []multiaddr.Multiaddr) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if len(expected) != len(actual) {
		t.Errorf("expected %v, got %v", expected, actual)
		return false
	}
	for i, e := range expected {
		if !e.Equal(actual[i]) {
			t.Errorf("expected %v, got %v", expected, actual)
			return false
		}
	}
	return true
}

// AssertMultiaddrsMatch is the same as AssertEqualMultiaddrs, but it ignores the order of the elements.
func AssertMultiaddrsMatch(t TestingT, expected, actual []multiaddr.Multiaddr) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	e := slices.Clone(expected)
	a := slices.Clone(actual)
	slices.SortFunc(e, func(a, b multiaddr.Multiaddr) int { return a.Compare(b) })
	slices.SortFunc(a, func(a, b multiaddr.Multiaddr) int { return a.Compare(b) })
	return AssertEqualMultiaddrs(t, e, a)
}

func AssertMultiaddrsContain(t TestingT, haystack []multiaddr.Multiaddr, needle multiaddr.Multiaddr) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	for _, h := range haystack {
		if h.Equal(needle) {
			return true
		}
	}
	t.Errorf("expected %v to contain %v", haystack, needle)
	return false
}

type MultiaddrMatcher struct {
	multiaddr.Multiaddr
}

// Implements the Matcher interface for gomock.Matcher
// Let's us use this struct in gomock tests. Example:
// Expect(mock.Method(gomock.Any(), multiaddrMatcher).Return(nil)
func (m MultiaddrMatcher) Matches(x interface{}) bool {
	if m2, ok := x.(multiaddr.Multiaddr); ok {
		return m.Equal(m2)
	}
	return false
}
