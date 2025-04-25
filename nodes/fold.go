package nodes

import "github.com/sidkurella/pipedream"

// FoldNode aggregates elements in a list of values according to the specified aggregation method.
type FoldNode struct {
	// Value to filter.
	Source pipedream.ValueBuilder

	// RightToLeft flips the direction of the fold to start with the last element of source, and work backwards.
	RightToLeft bool

	// StartValue sets the starting value of the accumulator.
	StartValue pipedream.ValueBuilder

	// TODO: How to specify the actual aggregation method itself?

	// Name to save the aggregated result into the pipeline context.
	SaveToName string
}
