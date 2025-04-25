package pipedream

import "fmt"

var ErrTerminatedEarly = fmt.Errorf("pipeline terminated early on an unhandled branch")

// Shared interface for nodes.
// If you want you can define your own that fits this pattern.
type Node interface {
	Execute() error
}

// A Pipeline is made up of a sequence of nodes that are executed in order.
type Pipeline struct {
	Nodes []Node
}
