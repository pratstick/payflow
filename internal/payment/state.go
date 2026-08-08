package payment

import "fmt"

var validTransitions = map[Status]map[Status]struct{}{
	StatusCreated: {
		StatusProcessing: {},
		StatusFailed:     {},
	},
	StatusProcessing: {
		StatusSuccess: {},
		StatusFailed:  {},
	},
	StatusSuccess: {
		StatusRefunded: {},
	},
	StatusFailed:   {},
	StatusRefunded: {},
}

func ValidateTransition(from, to Status) error {
	if from == to {
		return nil
	}
	next, ok := validTransitions[from]
	if !ok {
		return fmt.Errorf("unknown current status %q", from)
	}
	if _, ok := next[to]; !ok {
		return ErrInvalidStateTransition
	}
	return nil
}
