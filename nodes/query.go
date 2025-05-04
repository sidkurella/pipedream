package nodes

import (
	"fmt"

	"github.com/sidkurella/pipedream"
)

var ErrDataSourceNotFound = fmt.Errorf("data source not found")

// Queries a data source for data which is then saved to the pipeline context.
type QueryNode struct {
	// Data source to query from.
	DataSourceName string

	// Params to query by.
	Params pipedream.ValueBuilder

	// Name to save this into the pipeline context.
	SaveToName string
}

func (q QueryNode) Execute(pctx pipedream.PipelineContext) error {
	// Build input parameters.
	params, err := q.Params.Build(pctx)
	if err != nil {
		return fmt.Errorf("failed to build params to query %s: %w", q.DataSourceName, err)
	}

	// Get the data source with the defined name.
	// Assuming PipelineContext has a method to get data sources.
	dataSource, found := pctx.GetDataSource(q.DataSourceName)
	if !found {
		return fmt.Errorf("%w: %s", ErrDataSourceNotFound, q.DataSourceName)
	}

	// Query the data source.
	// Assuming the DataSource interface has a Query method.
	result, err := dataSource.Query(params)
	if err != nil {
		return fmt.Errorf("failed to query data source %s: %w", q.DataSourceName, err)
	}

	// Save the result to the context.
	pctx.SetValue(q.SaveToName, result)

	return nil
}
