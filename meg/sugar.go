package meg

import (
	"errors"
)

type Pattern = func(next *MatchState) *MatchState

func PatternToMatchState(states ...Pattern) *MatchState {
	nextState := &MatchState{kind: done}
	for i := len(states) - 1; i >= 0; i-- {
		nextState = states[i](nextState)
	}
	return nextState
}

func Cat(left, right Pattern) Pattern {
	return func(next *MatchState) *MatchState {
		return left(right(next))
	}
}

func Or(p ...Pattern) Pattern {
	return func(next *MatchState) *MatchState {
		if len(p) == 0 {
			return next
		}
		if len(p) == 1 {
			return p[0](next)
		}

		return &MatchState{
			kind:      split,
			next:      p[0](next),
			nextSplit: Or(p[1:]...)(next),
		}
	}
}

var errAlreadyCapture = errors.New("already captured")

func captureOneValueOrErr(val *string) captureFunc {
	if val == nil {
		return nil
	}
	var set bool
	f := func(s string) error {
		if set {
			*val = ""
			return errAlreadyCapture
		}
		*val = s
		return nil
	}
	return &f
}

func captureMany(vals *[]string) captureFunc {
	if vals == nil {
		return nil
	}
	f := func(s string) error {
		*vals = append(*vals, s)
		return nil
	}
	return &f
}

func captureValWithF(code int, f captureFunc) Pattern {
	return func(next *MatchState) *MatchState {
		return &MatchState{
			kind:    matchCode,
			capture: f,
			code:    code,
			next:    next,
		}
	}
}

func Val(code int) Pattern {
	return CaptureVal(code, nil)
}

func CaptureVal(code int, val *string) Pattern {
	return captureValWithF(code, captureOneValueOrErr(val))
}

func ZeroOrMore(code int) Pattern {
	return CaptureZeroOrMore(code, nil)
}

func captureZeroOrMoreWithF(code int, f captureFunc) Pattern {
	return func(next *MatchState) *MatchState {
		match := &MatchState{
			code:    code,
			capture: f,
		}
		s := &MatchState{
			kind:      split,
			next:      match,
			nextSplit: next,
		}
		match.next = s // Loop back to the split.
		return s
	}
}

func CaptureZeroOrMore(code int, vals *[]string) Pattern {
	return captureZeroOrMoreWithF(code, captureMany(vals))
}

func OneOrMore(code int) Pattern {
	return CaptureOneOrMore(code, nil)
}
func CaptureOneOrMore(code int, vals *[]string) Pattern {
	f := captureMany(vals)
	return func(next *MatchState) *MatchState {
		return captureValWithF(code, f)(captureZeroOrMoreWithF(code, f)(next))
	}
}

func Optional(s Pattern) Pattern {
	return func(next *MatchState) *MatchState {
		return &MatchState{
			kind:      split,
			next:      s(next),
			nextSplit: next,
		}
	}
}
