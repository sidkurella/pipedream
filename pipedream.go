package pipedream

import (
	"context"
	"fmt"
)

var ErrDataSourceNotFound = fmt.Errorf("data source not found")

type ExecutionContext struct {
	dataSources map[string]DataSource
	returnValue any
}

func (e ExecutionContext) GetDataSource(name string) (DataSource, error) {
	dataSource, ok := e.dataSources[name]
	if !ok {
		return DataSource{}, ErrDataSourceNotFound
	}

	return dataSource, nil
}

type PipelineExecutor struct {
	ectx ExecutionContext
}

func NewPipelineExecutor() *PipelineExecutor {
	return &PipelineExecutor{
		ectx: ExecutionContext{dataSources: map[string]DataSource{}},
	}
}

func (p *PipelineExecutor) RegisterDataSource(
	source DataSource,
) {
	p.ectx.dataSources[source.Name] = source
}

func (p *PipelineExecutor) Execute(ctx context.Context, pipeline Pipeline) (any, error) {
	// TODO
	return nil, nil
}
