package nodes

import "github.com/sidkurella/pipedream"

// Queries a data source for data which is then saved to the pipeline context.
type QueryNode struct {
	// Data source to query from.
	DataSourceName string

	// Params to query by.
	Params pipedream.ValueBuilder

	// Name to save this into the pipeline context.
	SaveToName string
}
