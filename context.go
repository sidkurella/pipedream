package pipedream

import "maps"

type PipelineContext struct {
	values map[string]any
}

// GetValue gets the value with the specified name, if it exists.
func (p PipelineContext) GetValue(k string) (any, bool) {
	v, b := p.values[k]
	return v, b
}

// SetValue sets the value provided to the specified name, returning any value that was already there.
func (p PipelineContext) SetValue(k string, v any) (any, bool) {
	oldV, b := p.values[k]
	p.values[k] = v

	return oldV, b
}

func (p PipelineContext) Clone() PipelineContext {
	return PipelineContext{
		values: maps.Clone(p.values),
	}
}
