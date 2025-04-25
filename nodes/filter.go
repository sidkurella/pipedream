package nodes

import "github.com/sidkurella/pipedream"

// Filters out data from a value and saves the filtered result to the pipeline context.
type FilterNode struct {
	// Condition to test for.
	Condition pipedream.Condition

	// Value to filter.
	Source pipedream.ValueBuilder

	// Exclude inverts the filter so only elements passing the condition are excluded instead.
	// Default is that elements passing the condition are kept.
	Exclude bool

	// Name to save the filtered result into the pipeline context.
	SaveToName string
}
