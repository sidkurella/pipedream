package pipedream

import (
	"cmp"
	"fmt"
	"reflect"
)

var ErrInvalidCondition = fmt.Errorf("condition has no valid operand or builders")
var ErrComparisonFailed = fmt.Errorf("failed to compare values")
var ErrIncompatibleTypes = fmt.Errorf("cannot compare incompatible types")
var ErrOperationNotSupported = fmt.Errorf("comparison operation not supported for these types")
var ErrNilCondition = fmt.Errorf("nil condition provided")
var ErrNilCustomCompareFunc = fmt.Errorf("custom compare function is nil")
var ErrTypeAssertionFailed = fmt.Errorf("failed to assert value to expected type for custom comparison")

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

// Returns a string representation of the operand.
func (op ConditionOperand) String() string {
	switch op {
	case ConditionEqual:
		return "=="
	case ConditionNotEqual:
		return "!="
	case ConditionGreaterThan:
		return ">"
	case ConditionLessThan:
		return "<"
	case ConditionGreaterThanOrEqual:
		return ">="
	case ConditionLessThanOrEqual:
		return "<="
	case ConditionOperandInvalid:
		return "Invalid"
	default:
		return fmt.Sprintf("Unknown(%d)", op)
	}
}

// Condition represents a condition that can be true or false based on context.
type Condition interface {
	// Evaluate checks the condition against the given PipelineContext.
	Evaluate(pctx PipelineContext) (bool, error)
}

// ValueCondition compares a left-hand-side (LHS) value to a right-hand-side (RHS) value
// using one of a list of operands. The LHS and RHS values are obtained using ValueBuilders.
type ValueCondition struct {
	LHS ValueBuilder
	RHS ValueBuilder

	Operand ConditionOperand
}

// Evaluate implements the Condition interface for ValueCondition.
// It builds the LHS and RHS values using the provided PipelineContext and compares them.
func (c *ValueCondition) Evaluate(pctx PipelineContext) (bool, error) {
	if c.LHS == nil || c.RHS == nil || c.Operand == ConditionOperandInvalid {
		return false, ErrInvalidCondition
	}

	lhsVal, err := c.LHS.Build(pctx)
	if err != nil {
		return false, fmt.Errorf("evaluating LHS: %w", err)
	}

	rhsVal, err := c.RHS.Build(pctx)
	if err != nil {
		return false, fmt.Errorf("evaluating RHS: %w", err)
	}

	// Now compare lhsVal and rhsVal based on c.Operand
	return compareValues(lhsVal, rhsVal, c.Operand)
}

// AndCondition only evaluates to true if all the conditions inside also evaluate to true.
// Short-circuits as soon as one condition is false.
type AndCondition struct {
	Conditions []Condition
}

// Evaluate implements the Condition interface for AndCondition.
// It passes the PipelineContext to each child condition.
func (c *AndCondition) Evaluate(pctx PipelineContext) (bool, error) {
	if len(c.Conditions) == 0 {
		// An empty AND condition is typically considered true.
		return true, nil
	}
	for _, cond := range c.Conditions {
		if cond == nil { // Add check for nil condition in the slice
			return false, ErrNilCondition
		}
		res, err := cond.Evaluate(pctx)
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

// OrCondition evaluates to true if any of the conditions inside evaluates to true.
// Short-circuits as soon as one condition is true.
type OrCondition struct {
	Conditions []Condition
}

// Evaluate implements the Condition interface for OrCondition.
// It passes the PipelineContext to each child condition.
func (c *OrCondition) Evaluate(pctx PipelineContext) (bool, error) {
	if len(c.Conditions) == 0 {
		// An empty OR condition is typically considered false.
		return false, nil
	}
	for _, cond := range c.Conditions {
		if cond == nil {
			return false, ErrNilCondition
		}
		res, err := cond.Evaluate(pctx) // Pass context
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

// Generic comparison function for ordered types
func compareOrdered[T cmp.Ordered](l, r T, op ConditionOperand) (bool, error) {
	switch op {
	case ConditionEqual:
		return l == r, nil
	case ConditionNotEqual:
		return l != r, nil
	case ConditionGreaterThan:
		return l > r, nil
	case ConditionLessThan:
		return l < r, nil
	case ConditionGreaterThanOrEqual:
		return l >= r, nil
	case ConditionLessThanOrEqual:
		return l <= r, nil
	default:
		return false, ErrInvalidCondition
	}
}

// compareValues performs the comparison between two values based on the operand.
// Handles basic types (numeric, string, bool), nils, and numeric type coercion.
func compareValues(lhs, rhs any, op ConditionOperand) (bool, error) {
	lhsV := reflect.ValueOf(lhs)
	rhsV := reflect.ValueOf(rhs)
	lhsK := lhsV.Kind()
	rhsK := rhsV.Kind()

	// 1. Handle nil comparisons explicitly.
	lhsIsNil := lhs == nil || (lhsK == reflect.Pointer && lhsV.IsNil())
	rhsIsNil := rhs == nil || (rhsK == reflect.Pointer && rhsV.IsNil())

	if lhsIsNil || rhsIsNil {
		switch op {
		case ConditionEqual:
			return lhsIsNil == rhsIsNil, nil
		case ConditionNotEqual:
			return lhsIsNil != rhsIsNil, nil
		default:
			return false, fmt.Errorf("%w: comparison with nil for operand %s", ErrOperationNotSupported, op.String())
		}
	}

	// 2. Handle identical types
	if lhsV.Type() == rhsV.Type() {
		switch {
		case isSignedInteger(lhsK):
			return compareOrdered(lhsV.Int(), rhsV.Int(), op)
		case isUnsignedInteger(lhsK):
			return compareOrdered(lhsV.Uint(), rhsV.Uint(), op)
		case isFloat(lhsK):
			return compareOrdered(lhsV.Float(), rhsV.Float(), op)
		case lhsK == reflect.String:
			return compareOrdered(lhsV.String(), rhsV.String(), op)
		case lhsK == reflect.Bool:
			// Bool is not cmp.Ordered, needs separate handling
			return compareBool(lhsV.Bool(), rhsV.Bool(), op)
		case lhsV.Type().Comparable():
			// Other comparable types (structs, arrays etc) - only support ==/!=
			if op == ConditionEqual {
				return lhsV.Interface() == rhsV.Interface(), nil
			}
			if op == ConditionNotEqual {
				return lhsV.Interface() != rhsV.Interface(), nil
			}
			return false, fmt.Errorf("%w: operand %s not supported for comparable type %s", ErrOperationNotSupported, op.String(), lhsV.Type())
		default: // Non-comparable identical types
			if op == ConditionEqual {
				return reflect.DeepEqual(lhsV.Interface(), rhsV.Interface()), nil
			}
			if op == ConditionNotEqual {
				return !reflect.DeepEqual(lhsV.Interface(), rhsV.Interface()), nil
			}
			return false, fmt.Errorf("%w: cannot compare non-comparable type %s with operand %s", ErrOperationNotSupported, lhsV.Type(), op.String())
		}
	}

	if isNumeric(lhsK) && isNumeric(rhsK) {
		lIsSigned := isSignedInteger(lhsK)
		rIsSigned := isSignedInteger(rhsK)
		lIsUnsigned := isUnsignedInteger(lhsK)
		rIsUnsigned := isUnsignedInteger(rhsK)

		// 3.1 Both signed integers -> Promote to int64
		if lIsSigned && rIsSigned {
			lInt, _ := convertToSignedInt64(lhsV)
			rInt, _ := convertToSignedInt64(rhsV)
			return compareOrdered(lInt, rInt, op)
		}

		// 3.2 Both unsigned integers -> Promote to uint64
		if lIsUnsigned && rIsUnsigned {
			lUint, _ := convertToUnsignedUint64(lhsV)
			rUint, _ := convertToUnsignedUint64(rhsV)
			return compareOrdered(lUint, rUint, op)
		}

		// 3.3 Mixed signed/unsigned integers -> Use special comparison logic
		if lIsSigned && rIsUnsigned {
			sVal, _ := convertToSignedInt64(lhsV)
			uVal, _ := convertToUnsignedUint64(rhsV)
			return compareMixedInt(sVal, uVal, op) // LHS=Signed, RHS=Unsigned
		}
		if lIsUnsigned && rIsSigned {
			uVal, _ := convertToUnsignedUint64(lhsV)
			sVal, _ := convertToSignedInt64(rhsV)
			return compareMixedIntReversed(uVal, sVal, op) // LHS=Unsigned, RHS=Signed
		}

		// 3.4 Other mixed numeric types (must involve at least one float) -> Convert both to float64
		// This now covers: signed_int vs float, unsigned_int vs float, float vs float (different types)
		lFlt, lOk := convertToFloat64(lhsV)
		rFlt, rOk := convertToFloat64(rhsV)
		if lOk && rOk {
			return compareOrdered(lFlt, rFlt, op)
		} else {
			// Should generally not happen if isNumeric check passed, but safeguard.
			return false, fmt.Errorf("%w: failed to convert numeric types %s and %s for comparison (float path)", ErrComparisonFailed, lhsV.Type(), rhsV.Type())
		}
	}

	// Default case: types are different and are not coercable to each other.
	return false, fmt.Errorf("%w: cannot compare types %s and %s with operand %s", ErrIncompatibleTypes, lhsV.Type(), rhsV.Type(), op.String())
}

// --- Helper functions for comparison ---

func isSignedInteger(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func isUnsignedInteger(k reflect.Kind) bool {
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	default:
		return false
	}
}

func isInteger(k reflect.Kind) bool {
	return isSignedInteger(k) || isUnsignedInteger(k)
}

func isFloat(k reflect.Kind) bool {
	switch k {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isNumeric(k reflect.Kind) bool {
	return isInteger(k) || isFloat(k)
}

func convertToSignedInt64(v reflect.Value) (int64, bool) {
	k := v.Kind()
	if isSignedInteger(k) {
		return v.Int(), true
	}
	// We don't attempt converting unsigned to signed here.
	return 0, false
}

func convertToUnsignedUint64(v reflect.Value) (uint64, bool) {
	k := v.Kind()
	if isUnsignedInteger(k) {
		return v.Uint(), true
	}
	// We don't attempt converting signed non-negative to unsigned here.
	return 0, false
}

func convertToFloat64(v reflect.Value) (float64, bool) {
	k := v.Kind()
	switch {
	case isSignedInteger(k):
		return float64(v.Int()), true
	case isUnsignedInteger(k):
		// Note: Potential precision loss converting uint64 > 2^53 to float64
		return float64(v.Uint()), true
	case isFloat(k):
		return v.Float(), true
	default:
		return 0, false // Not a numeric type handled here
	}
}

// --- Comparison helpers for specific types ---

// compareMixedInt compares signed (s, LHS) vs unsigned (u, RHS).
func compareMixedInt(s int64, u uint64, op ConditionOperand) (bool, error) {
	if s < 0 { // Signed is negative, Unsigned is non-negative
		switch op {
		case ConditionEqual:
			return false, nil
		case ConditionNotEqual:
			return true, nil
		case ConditionGreaterThan:
			return false, nil // s > u is false
		case ConditionLessThan:
			return true, nil // s < u is true
		case ConditionGreaterThanOrEqual:
			return false, nil // s >= u is false
		case ConditionLessThanOrEqual:
			return true, nil // s <= u is true
		default:
			return false, fmt.Errorf("%w: invalid operand %s", ErrInvalidCondition, op.String())
		}
	} else { // Signed is non-negative, safe to convert to uint64
		sUint := uint64(s)
		return compareOrdered(sUint, u, op) // Compare sUint (LHS) vs u (RHS)
	}
}

// compareMixedIntReversed compares unsigned (u, LHS) vs signed (s, RHS).
func compareMixedIntReversed(u uint64, s int64, op ConditionOperand) (bool, error) {
	if s < 0 { // Unsigned is non-negative, Signed is negative
		switch op {
		case ConditionEqual:
			return false, nil
		case ConditionNotEqual:
			return true, nil
		case ConditionGreaterThan:
			return true, nil // u > s is true
		case ConditionLessThan:
			return false, nil // u < s is false
		case ConditionGreaterThanOrEqual:
			return true, nil // u >= s is true
		case ConditionLessThanOrEqual:
			return false, nil // u <= s is false
		default:
			return false, fmt.Errorf("%w: invalid operand %s", ErrInvalidCondition, op.String())
		}
	} else { // Signed is non-negative, safe to convert to uint64
		sUint := uint64(s)
		return compareOrdered(u, sUint, op) // Compare u (LHS) vs sUint (RHS)
	}
}

func compareBool(l, r bool, op ConditionOperand) (bool, error) {
	switch op {
	case ConditionEqual:
		return l == r, nil
	case ConditionNotEqual:
		return l != r, nil
	default: // >, <, >=, <= make no sense for bools
		return false, fmt.Errorf("%w: operand %s not supported for bool comparison", ErrOperationNotSupported, op.String())
	}
}

// CustomCompareFunc defines the signature for a user-provided comparison function
// between two values of a specific type T.
type CustomCompareFunc[T any] func(lhs T, rhs T) (bool, error)

// CustomCompareCondition allows defining a condition using a custom comparison function
// applied to values built by LHS and RHS ValueBuilders. It expects both builders
// to produce values of type T.
type CustomCompareCondition[T any] struct {
	LHS         ValueBuilder
	RHS         ValueBuilder
	CompareFunc CustomCompareFunc[T]
}

// Evaluate implements the Condition interface for CustomCompareCondition.
// It builds LHS and RHS values, asserts they are of type T, and then executes
// the custom comparison function.
func (c *CustomCompareCondition[T]) Evaluate(pctx PipelineContext) (bool, error) {
	if c.LHS == nil || c.RHS == nil {
		return false, fmt.Errorf("%w: LHS or RHS builder is nil in CustomCompareCondition", ErrInvalidCondition)
	}
	if c.CompareFunc == nil {
		return false, ErrNilCustomCompareFunc
	}

	lhsValAny, err := c.LHS.Build(pctx)
	if err != nil {
		return false, fmt.Errorf("evaluating LHS for custom comparison: %w", err)
	}

	rhsValAny, err := c.RHS.Build(pctx)
	if err != nil {
		return false, fmt.Errorf("evaluating RHS for custom comparison: %w", err)
	}

	// Get the expected type name using reflection on a zero value of T
	// Note: This might not be perfect for interface types, but good for concrete types.
	var zeroT T
	expectedTypeName := reflect.TypeOf(zeroT).String() // Get type name for error message

	// Attempt type assertion for LHS
	lhsValTyped, ok := lhsValAny.(T)
	if !ok {
		actualTypeName := "nil"
		if lhsValAny != nil {
			actualTypeName = reflect.TypeOf(lhsValAny).String()
		}
		return false, fmt.Errorf("%w: expected LHS type '%s', got '%s'", ErrTypeAssertionFailed, expectedTypeName, actualTypeName)
	}

	// Attempt type assertion for RHS
	rhsValTyped, ok := rhsValAny.(T)
	if !ok {
		actualTypeName := "nil"
		if rhsValAny != nil {
			actualTypeName = reflect.TypeOf(rhsValAny).String()
		}
		return false, fmt.Errorf("%w: expected RHS type '%s', got '%s'", ErrTypeAssertionFailed, expectedTypeName, actualTypeName)
	}

	// Call the custom comparison function with the correctly typed values
	return c.CompareFunc(lhsValTyped, rhsValTyped)
}
