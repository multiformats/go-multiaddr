// package meg implements Regular Expressions for multiaddr Components. It's short for "Megular Expressions"
package meg

// The developer is assumed to be familiar with the Thompson NFA approach to
// regex before making changes to this file. Refer to
// https://swtch.com/~rsc/regexp/regexp1.html for an introduction.

import (
	"fmt"
)

type stateKind uint8

const (
	matchCode stateKind = iota
	split
	done
)

// MatchState is the Thompson NFA for a regular expression.
type Matcher = *MatchState

type MatchState struct {
	capture   captureFunc
	next      *MatchState
	nextSplit *MatchState

	kind       stateKind
	generation int
	code       int
}

type captureFunc *func(string) error

type capture struct {
	f    captureFunc
	v    string
	prev *capture
}

type statesAndCaptures struct {
	states   []*MatchState
	captures []*capture
}

func (s *MatchState) String() string {
	return fmt.Sprintf("state{kind: %d, generation: %d, code: %d}", s.kind, s.generation, s.code)
}

type Matchable interface {
	Code() int
	Value() string // Used when capturing the value
}

// Match returns whether the given Components match the Pattern defined in MatchState.
// Errors are used to communicate capture errors.
// If the error is non-nil the returned bool will be false.
func Match[S ~[]T, T Matchable](s *MatchState, components S) (bool, error) {
	listGeneration := s.generation + 1               // Start at the last generation + 1
	defer func() { s.generation = listGeneration }() // In case we reuse this state, store our highest generation number

	currentStates := statesAndCaptures{
		states:   make([]*MatchState, 0, 16),
		captures: make([]*capture, 0, 16),
	}
	nextStates := statesAndCaptures{
		states:   make([]*MatchState, 0, 16),
		captures: make([]*capture, 0, 16),
	}

	currentStates = appendState(currentStates, s, nil, listGeneration)

	for _, c := range components {
		if len(currentStates.states) == 0 {
			return false, nil
		}
		for i, s := range currentStates.states {
			if s.kind == matchCode && s.code == c.Code() {
				cm := currentStates.captures[i]
				if s.capture != nil {
					next := &capture{
						f: s.capture,
						v: c.Value(),
					}
					if cm == nil {
						cm = next
					} else {
						next.prev = cm
						cm = next
					}
					currentStates.captures[i] = cm
				}
				nextStates = appendState(nextStates, s.next, cm, listGeneration)
			}
		}
		currentStates, nextStates = nextStates, currentStates
		nextStates.states = nextStates.states[:0]
		nextStates.captures = nextStates.captures[:0]
		listGeneration++
	}

	for i, s := range currentStates.states {
		if s.kind == done {
			// We found a complete path. Run the captures now
			c := currentStates.captures[i]
			for c != nil {
				if err := (*c.f)(c.v); err != nil {
					return false, err
				}
				c = c.prev
			}
			return true, nil
		}
	}
	return false, nil
}

func appendState(arr statesAndCaptures, s *MatchState, c *capture, listGeneration int) statesAndCaptures {
	if s == nil || s.generation == listGeneration {
		return arr
	}
	s.generation = listGeneration
	if s.kind == split {
		arr = appendState(arr, s.next, c, listGeneration)
		arr = appendState(arr, s.nextSplit, c, listGeneration)
	} else {
		arr.states = append(arr.states, s)
		arr.captures = append(arr.captures, c)
	}
	return arr
}
