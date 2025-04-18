package pipedream

import "fmt"

var ErrTerminatedEarly = fmt.Errorf("pipeline terminated early on an unhandled branch")

type Node interface {
	isPipelineNode()
}

type Pipeline struct {
	Nodes []Node `json:"nodes"`
}

type QueryNode struct {
	// Data source to query from.
	DataSourceName string
	// Params to query by.
	Params ValueBuilder
	// Name to save this into the pipeline context.
	SaveToName string
}

func (q QueryNode) isPipelineNode() {}

type BranchNode struct {
	// Condition that this branch should evaluate
	Condition Condition

	// Pipeline to execute if the condition is true.
	// If nil and the condition is true, the pipeline execution will return ErrTerminatedEarly.
	TruePipeline *Pipeline

	// Pipeline to execute if the condition is false.
	// If nil and the condition is false, the pipeline execution will return ErrTerminatedEarly.
	FalsePipeline *Pipeline
}

// Builds a value and returns it as the end goal of the pipeline.
type ReturnNode struct {
	ValueBuilder ValueBuilder
}
