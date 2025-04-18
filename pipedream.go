package pipedream

import "context"

type PipelineExecutor struct {
	dataSources map[string]DataSource
}

func NewPipelineExecutor() *PipelineExecutor {
	return &PipelineExecutor{}
}

func (p *PipelineExecutor) RegisterDataSource(
	source DataSource,
) {
	p.dataSources[source.Name] = source
}

func (p *PipelineExecutor) Execute(ctx context.Context, pipeline Pipeline) (any, error) {
	// TODO
	return nil, nil
}
