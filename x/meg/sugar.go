package meg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Pattern is essentially a curried MatchState.
// Given the slice of current MatchStates and a handle (int index) to the next
// MatchState, it returns a handle to the inserted MatchState.
type Pattern = func(states *[]MatchState, nextIdx int) int

type Matcher struct {
	states   []MatchState
	startIdx int
}

func (s Matcher) String() string {
	states := make([]string, len(s.states))
	for i, state := range s.states {
		states[i] = state.String() + "@" + strconv.Itoa(i)
	}
	return fmt.Sprintf("RootMatchState{states: [%s], startIdx: %d}", strings.Join(states, ", "), s.startIdx)
}

func PatternToMatcher(patterns ...Pattern) Matcher {
	// Preallocate a slice to hold the MatchStates.
	// Avoids small allocations for each pattern.
	// The number is chosen experimentally. It is subject to change.
	states := make([]MatchState, 0, len(patterns)*3)
	// Append the done state.
	states = append(states, MatchState{codeOrKind: done})
	nextIdx := len(states) - 1
	// Build the chain by composing patterns from right to left.
	for i := len(patterns) - 1; i >= 0; i-- {
		nextIdx = patterns[i](&states, nextIdx)
	}
	return Matcher{states: states, startIdx: nextIdx}
}

func Cat(left, right Pattern) Pattern {
	return func(states *[]MatchState, nextIdx int) int {
		// First run the right pattern, then feed the result into left.
		return left(states, right(states, nextIdx))
	}
}

func Or(p ...Pattern) Pattern {
	return func(states *[]MatchState, nextIdx int) int {
		if len(p) == 0 {
			return nextIdx
		}
		// Evaluate the last pattern and use its result as the initial accumulator.
		accum := p[len(p)-1](states, nextIdx)
		// Iterate backwards from the second-to-last pattern to the first.
		for i := len(p) - 2; i >= 0; i-- {
			leftIdx := p[i](states, nextIdx)
			newState := MatchState{
				next:       leftIdx,
				codeOrKind: storeSplitIdx(accum),
			}
			*states = append(*states, newState)
			accum = len(*states) - 1
		}
		return accum
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
	return f
}

func captureMany(vals *[]string) captureFunc {
	if vals == nil {
		return nil
	}
	f := func(s string) error {
		*vals = append(*vals, s)
		return nil
	}
	return f
}

func captureValWithF(code int, f captureFunc) Pattern {
	return func(states *[]MatchState, nextIdx int) int {
		newState := MatchState{
			capture:    f,
			codeOrKind: code,
			next:       nextIdx,
		}
		*states = append(*states, newState)
		return len(*states) - 1
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
	return func(states *[]MatchState, nextIdx int) int {
		// Create the match state.
		matchState := MatchState{
			codeOrKind: code,
			capture:    f,
		}
		*states = append(*states, matchState)
		matchIdx := len(*states) - 1

		// Create the split state that branches to the match state and to the next state.
		s := MatchState{
			next:       matchIdx,
			codeOrKind: storeSplitIdx(nextIdx),
		}
		*states = append(*states, s)
		splitIdx := len(*states) - 1

		// Close the loop: update the match state's next field.
		(*states)[matchIdx].next = splitIdx

		return splitIdx
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
	return func(states *[]MatchState, nextIdx int) int {
		// First attach the zero-or-more loop.
		zeroOrMoreIdx := captureZeroOrMoreWithF(code, f)(states, nextIdx)
		// Then put the capture state before the loop.
		return captureValWithF(code, f)(states, zeroOrMoreIdx)
	}
}

func Optional(s Pattern) Pattern {
	return func(states *[]MatchState, nextIdx int) int {
		newState := MatchState{
			next:       s(states, nextIdx),
			codeOrKind: storeSplitIdx(nextIdx),
		}
		*states = append(*states, newState)
		return len(*states) - 1
	}
}
