package pipedream

import "fmt"

var ErrInvalidCondition = fmt.Errorf("condition has no custom func nor valid operand")

type ConditionOperand int

const (
	ConditionOperandInvalid     ConditionOperand = iota
	ConditionEqual                               // ==
	ConditionNotEqual                            // !=
	ConditionGreaterThan                         // >
	ConditionLessThan                            // <
	ConditionGreaterThanOrEqual                  // >=
	ConditionLessThanOrEqual                     // <=
)

// Condition represents a condition that can be true or false.
type Condition interface {
	Evaluate() (bool, error)
}

// ValueCondition compares a left-hand-side (LHS) value to a right-hand-side (RHS) value using one of a list of operands.
type ValueCondition struct {
	LHS ValueBuilder
	RHS ValueBuilder

	Operand ConditionOperand
}

func (c *ValueCondition) Evaluate() (bool, error) {
	// TODO
	return false, nil
}

// AndCondition only evaluates to true if all the conditions inside also evaluate to true.
// Short-circuits as soon as one condition is false.
type AndCondition struct {
	Conditions []Condition
}

// Evaluate implements the Condition interface for AndCondition.
func (c *AndCondition) Evaluate() (bool, error) {
	if len(c.Conditions) == 0 {
		// An empty AND condition is typically considered true.
		return true, nil
	}
	for _, cond := range c.Conditions {
		res, err := cond.Evaluate()
		if err != nil {
			// Return the first error encountered.
			return false, err
		}
		if !res {
			// Short-circuit: if any condition is false, the AND is false.
			return false, nil
		}
	}
	// All conditions evaluated to true without error.
	return true, nil
}

// OrCondiiton evaluates to true if any of the condiitons inside evaluates to true.
// Short-circuits as soon as one condition is true.
type OrCondition struct {
	Conditions []Condition
}

// Evaluate implements the Condition interface for OrCondition.
func (c *OrCondition) Evaluate() (bool, error) {
	if len(c.Conditions) == 0 {
		// An empty OR condition is typically considered false.
		return false, nil
	}
	for _, cond := range c.Conditions {
		res, err := cond.Evaluate()
		if err != nil {
			// Return the first error encountered.
			return false, err
		}
		if res {
			// Short-circuit: if any condition is true, the OR is true.
			return true, nil
		}
	}
	// All conditions evaluated to false without error.
	return false, nil
}
