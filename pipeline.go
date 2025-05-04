package pipedream

import "fmt"

// Sentinel error to return to stop pipeline execution.
// Otherwise, the pipeline will proceed to the next node in the chain.
var ErrPipelineExecutionStop = fmt.Errorf("pipeline execution stop")

// Shared interface for nodes.
// If you want you can define your own that fits this pattern.
type Node interface {
	Execute(ectx ExecutionContext, pctx PipelineContext) error
}

// A Pipeline is made up of a sequence of nodes that are executed in order.
type Pipeline struct {
	Nodes []Node
}
