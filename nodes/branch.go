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

	// CloneContext indicates if the child pipeline should receive a copy of the parent pipeline context.
	// If false (default), it will reuse the context. This may be useful if there are other nodes to run
	// after this one no matter if the condition was true or false.
	// If true, the context will be cloned.
	CloneContext bool
}
