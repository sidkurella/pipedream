package pipedream

import "fmt"

var ErrNoContextKeyProvided = fmt.Errorf("no context key provided")
var ErrValueNotFoundInContext = fmt.Errorf("value not found in context")

// ValueBuilder builds a concrete value.
// For more complex cases you may need to write your own.
type ValueBuilder interface {
	Build(pctx PipelineContext) (any, error)
}

// LiteralValue ignores the context and always returns the literal value given.
type LiteralValue[T any] struct {
	Value T
}

func (l LiteralValue[T]) Build(pctx PipelineContext) (any, error) {
	return l.Value, nil
}

// DynamicValue uses a ValueGetter to extract a value from the context.
type DynamicValue struct {
	ContextKey string      // Key for the value to use from the context.
	Getter     ValueGetter // ValueGetter to use to process the context value.
	Key        any         // The key or keys needed by the Getter
}

// Build implements the ValueBuilder interface.
// It retrieves a value from the pipeline context using ContextKey.
// If a Getter is provided, it processes the retrieved value using the Getter and Key.
func (dv DynamicValue) Build(pctx PipelineContext) (any, error) {
	if dv.ContextKey == "" {
		return nil, ErrNoContextKeyProvided
	}

	// Retrieve the raw value from the context
	rawValue, found := pctx.GetValue(dv.ContextKey)
	if !found {
		return nil, ErrValueNotFoundInContext
	}

	// If a getter is specified, use it to process the value
	if dv.Getter != nil {
		return dv.Getter.GetValue(rawValue, dv.Key)
	}

	// Otherwise, return the raw value directly
	return rawValue, nil
}
