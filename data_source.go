package pipedream

import (
	"context"
	"fmt"
	"reflect"
)

var ErrParamTypeDoesNotMatch = fmt.Errorf("cannot convert provided params to appropriate type for this data source")

type DataGetter[T any, U any] func(ctx context.Context, params T) (U, error)

type DataSource struct {
	Name string

	getter any

	paramType    reflect.Type
	responseType reflect.Type
}

func NewDataSource[T any, U any](
	name string,
	getter DataGetter[T, U],
) DataSource {
	return DataSource{
		Name:         name,
		getter:       getter,
		paramType:    reflect.TypeFor[T](),
		responseType: reflect.TypeFor[U](),
	}
}

func (d DataSource) Get(ctx context.Context, params any) (any, error) {
	// Since this value is provided we need to check that it is of the correct type.
	paramsValue := reflect.ValueOf(params)
	if !paramsValue.CanConvert(d.paramType) {
		return nil, ErrParamTypeDoesNotMatch
	}

	// Don't really need to check these since it was provided at construction in a type-safe way.
	f := reflect.ValueOf(d.getter)
	resps := f.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		paramsValue.Convert(d.paramType),
	})

	resp, err := resps[0].Convert(d.responseType).Interface(), resps[1].Convert(reflect.TypeFor[error]()).Interface().(error)

	return resp, err
}
