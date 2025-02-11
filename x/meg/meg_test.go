package meg

import (
	"regexp"
	"slices"
	"testing"
	"testing/quick"
)

type codeAndValue struct {
	code int
	val  string // Uses the string type to ensure immutability.
}

// Code implements Matchable.
func (c codeAndValue) Code() int {
	return c.code
}

// Value implements Matchable.
func (c codeAndValue) Value() string {
	return c.val
}

var _ Matchable = codeAndValue{}

func TestSimple(t *testing.T) {
	type testCase struct {
		pattern        Matcher
		skipQuickCheck bool
		shouldMatch    [][]int
		shouldNotMatch [][]int
	}
	testCases :=
		[]testCase{
			{
				pattern:     PatternToMatchState(Val(0), Val(1)),
				shouldMatch: [][]int{{0, 1}},
				shouldNotMatch: [][]int{
					{0},
					{0, 0},
					{0, 1, 0},
				},
			},
			{
				pattern: PatternToMatchState(Optional(Val(1))),
				shouldMatch: [][]int{
					{1},
					{},
				},
				shouldNotMatch: [][]int{{0}},
			},
			{
				pattern: PatternToMatchState(Val(0), Val(1), Optional(Val(2))),
				shouldMatch: [][]int{
					{0, 1, 2},
					{0, 1},
				},
				shouldNotMatch: [][]int{
					{0},
					{0, 0},
					{0, 1, 0},
					{0, 1, 2, 0},
				}}, {
				pattern:        PatternToMatchState(Val(0), Val(1), OneOrMore(2)),
				skipQuickCheck: true,
				shouldMatch: [][]int{
					{0, 1, 2, 2, 2, 2},
					{0, 1, 2},
				},
				shouldNotMatch: [][]int{
					{0},
					{0, 0},
					{0, 1},
					{0, 1, 0},
					{0, 1, 1, 0},
					{0, 1, 2, 0},
				}},
		}

	for i, tc := range testCases {
		for _, m := range tc.shouldMatch {
			if matches, err := Match(tc.pattern, codesToCodeAndValue(m)); !matches {
				t.Fatalf("failed to match %v with %v. idx=%d. err=%v", m, tc.pattern, i, err)
			}
		}
		for _, m := range tc.shouldNotMatch {
			if matches, _ := Match(tc.pattern, codesToCodeAndValue(m)); matches {
				t.Fatalf("failed to not match %v with %v. idx=%d", m, tc.pattern, i)
			}
		}
		if tc.skipQuickCheck {
			continue
		}
		if err := quick.Check(func(notMatch []int) bool {
			for _, shouldMatch := range tc.shouldMatch {
				if slices.Equal(notMatch, shouldMatch) {
					// The random `notMatch` is actually something that shouldMatch. Skip it.
					return true
				}
			}
			matches, _ := Match(tc.pattern, codesToCodeAndValue(notMatch))
			return !matches
		}, &quick.Config{}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCapture(t *testing.T) {
	type setupStateAndAssert func() (Matcher, func())
	type testCase struct {
		setup setupStateAndAssert
		parts []codeAndValue
	}

	testCases :=
		[]testCase{
			{
				setup: func() (Matcher, func()) {
					var code0str string
					return PatternToMatchState(CaptureVal(0, &code0str), Val(1)), func() {
						if code0str != "hello" {
							panic("unexpected value")
						}
					}
				},
				parts: []codeAndValue{{0, "hello"}, {1, "world"}},
			},
			{
				setup: func() (Matcher, func()) {
					var code0strs []string
					return PatternToMatchState(CaptureOneOrMore(0, &code0strs), Val(1)), func() {
						if code0strs[0] != "hello" {
							panic("unexpected value")
						}
						if code0strs[1] != "world" {
							panic("unexpected value")
						}
					}
				},
				parts: []codeAndValue{{0, "hello"}, {0, "world"}, {1, ""}},
			},
		}

	_ = testCases
	for _, tc := range testCases {
		state, assert := tc.setup()
		if matches, _ := Match(state, tc.parts); !matches {
			t.Fatalf("failed to match %v with %v", tc.parts, state)
		}
		assert()
	}
}

func codesToCodeAndValue(codes []int) []codeAndValue {
	out := make([]codeAndValue, len(codes))
	for i, c := range codes {
		out[i] = codeAndValue{code: c}
	}
	return out
}

func bytesToCodeAndValue(codes []byte) []codeAndValue {
	out := make([]codeAndValue, len(codes))
	for i, c := range codes {
		out[i] = codeAndValue{code: int(c)}
	}
	return out
}

// FuzzMatchesRegexpBehavior fuzz tests the expression matcher by comparing it to the behavior of the regexp package.
func FuzzMatchesRegexpBehavior(f *testing.F) {
	bytesToRegexpAndPattern := func(exp []byte) (string, []Pattern) {
		if len(exp) < 3 {
			panic("regexp too short")
		}
		pattern := make([]Pattern, 0, len(exp)-2)
		for i, b := range exp {
			b = b % 32
			if i == 0 {
				exp[i] = '^'
				continue
			} else if i == len(exp)-1 {
				exp[i] = '$'
				continue
			}
			switch {
			case b < 26:
				exp[i] = b + 'a'
				pattern = append(pattern, Val(int(exp[i])))
			case i > 1 && b == 26:
				exp[i] = '?'
				pattern = pattern[:len(pattern)-1]
				pattern = append(pattern, Optional(Val(int(exp[i-1]))))
			case i > 1 && b == 27:
				exp[i] = '*'
				pattern = pattern[:len(pattern)-1]
				pattern = append(pattern, ZeroOrMore(int(exp[i-1])))
			case i > 1 && b == 28:
				exp[i] = '+'
				pattern = pattern[:len(pattern)-1]
				pattern = append(pattern, OneOrMore(int(exp[i-1])))
			default:
				exp[i] = 'a'
				pattern = append(pattern, Val(int(exp[i])))
			}
		}

		return string(exp), pattern
	}

	simplifyB := func(buf []byte) []byte {
		for i, b := range buf {
			buf[i] = (b % 26) + 'a'
		}
		return buf
	}

	f.Fuzz(func(t *testing.T, expRules []byte, corpus []byte) {
		if len(expRules) < 3 || len(expRules) > 1024 || len(corpus) > 1024 {
			return
		}
		corpus = simplifyB(corpus)
		regexpPattern, pattern := bytesToRegexpAndPattern(expRules)
		matched, err := regexp.Match(string(regexpPattern), corpus)
		if err != nil {
			// Malformed regex. Ignore
			return
		}
		p := PatternToMatchState(pattern...)
		otherMatched, _ := Match(p, bytesToCodeAndValue(corpus))
		if otherMatched != matched {
			t.Log("regexp", string(regexpPattern))
			t.Log("corpus", string(corpus))
			m2, err2 := regexp.Match(string(regexpPattern), corpus)
			t.Logf("regexp matched %v. %v. %v, %v. \n%v - \n%v", matched, err, m2, err2, regexpPattern, corpus)
			t.Logf("pattern %+v", pattern)
			t.Fatalf("mismatched results: %v %v %v", otherMatched, matched, p)
		}
	})

}
