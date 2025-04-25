package nodes

import "github.com/sidkurella/pipedream"

// Builds a value and returns it as the end goal of the pipeline.
type ReturnNode struct {
	ValueBuilder pipedream.ValueBuilder
}
