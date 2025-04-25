package nodes

import "github.com/sidkurella/pipedream"

// Executes one of two pipelines depending on the output of the condition
type BranchNode struct {
	// Condition that this branch should evaluate
	Condition pipedream.Condition

	// Pipeline to execute if the condition is true.
	// If nil and the condition is true, the pipeline execution will return ErrTerminatedEarly.
	TruePipeline *pipedream.Pipeline

	// Pipeline to execute if the condition is false.
	// If nil and the condition is false, the pipeline execution will return ErrTerminatedEarly.
	FalsePipeline *pipedream.Pipeline
}
