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

func captureOneBytesOrErr(val *[]byte) CaptureFunc {
	if val == nil {
		return nil
	}
	var set bool
	f := func(s Matchable) error {
		if set {
			*val = nil
			return errAlreadyCapture
		}
		*val = s.Bytes()
		return nil
	}
	return f
}

func captureOneStringValueOrErr(val *string) CaptureFunc {
	if val == nil {
		return nil
	}
	var set bool
	f := func(s Matchable) error {
		if set {
			*val = ""
			return errAlreadyCapture
		}
		*val = s.Value()
		return nil
	}
	return f
}

func captureManyBytes(vals *[][]byte) CaptureFunc {
	if vals == nil {
		return nil
	}
	f := func(s Matchable) error {
		*vals = append(*vals, s.Bytes())
		return nil
	}
	return f
}

func captureManyStrings(vals *[]string) CaptureFunc {
	if vals == nil {
		return nil
	}
	f := func(s Matchable) error {
		*vals = append(*vals, s.Value())
		return nil
	}
	return f
}

func CaptureWithF(code int, f CaptureFunc) Pattern {
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
	return CaptureStringVal(code, nil)
}

func CaptureStringVal(code int, val *string) Pattern {
	return CaptureWithF(code, captureOneStringValueOrErr(val))
}

func CaptureBytes(code int, val *[]byte) Pattern {
	return CaptureWithF(code, captureOneBytesOrErr(val))
}

func ZeroOrMore(code int) Pattern {
	return CaptureZeroOrMoreStringVals(code, nil)
}

func CaptureZeroOrMoreWithF(code int, f CaptureFunc) Pattern {
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

func CaptureZeroOrMoreBytes(code int, vals *[][]byte) Pattern {
	return CaptureZeroOrMoreWithF(code, captureManyBytes(vals))
}

func CaptureZeroOrMoreStringVals(code int, vals *[]string) Pattern {
	return CaptureZeroOrMoreWithF(code, captureManyStrings(vals))
}

func OneOrMore(code int) Pattern {
	return CaptureOneOrMoreStringVals(code, nil)
}

func CaptureOneOrMoreStringVals(code int, vals *[]string) Pattern {
	f := captureManyStrings(vals)
	return func(states *[]MatchState, nextIdx int) int {
		// First attach the zero-or-more loop.
		zeroOrMoreIdx := CaptureZeroOrMoreWithF(code, f)(states, nextIdx)
		// Then put the capture state before the loop.
		return CaptureWithF(code, f)(states, zeroOrMoreIdx)
	}
}

func CaptureOneOrMoreBytes(code int, vals *[][]byte) Pattern {
	f := captureManyBytes(vals)
	return func(states *[]MatchState, nextIdx int) int {
		// First attach the zero-or-more loop.
		zeroOrMoreIdx := CaptureZeroOrMoreWithF(code, f)(states, nextIdx)
		// Then put the capture state before the loop.
		return CaptureWithF(code, f)(states, zeroOrMoreIdx)
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
