package pipedream

type PipelineContext struct {
	values map[string]any
}

func (p PipelineContext) GetValue(k string) (any, bool) {
	v, b := p.values[k]
	return v, b
}
