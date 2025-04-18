package pipedream

import (
	"fmt"
	"reflect"
)

var ErrValueNotFound = fmt.Errorf("value with specified key does not exist")
var ErrKeyTypeInvalid = fmt.Errorf("key type is invalid for this input type")
var ErrImproperValueKind = fmt.Errorf("cannot get values from this input")
var ErrInputIsNil = fmt.Errorf("input is nil")
var ErrKeyIsEmpty = fmt.Errorf("key is empty or nil")
var ErrFieldIsUnexported = fmt.Errorf("field is unexported")

type ValueGetter interface {
	GetValue(input any, valueKey any) (any, error)
}

// If the input is a struct or pointer to a struct, uses reflection to get the field with the name specified.
// If the input is a map, gets the value for the key provided.
// If the input is an array or slice, gets the value at the specified index.
// If the input is a primitive, gets the value no matter what value key is provided.
// Otherwise, returns an error.
type DefaultValueGetter struct {
}

func (d DefaultValueGetter) GetValue(input any, valueKey any) (any, error) {
	if input == nil {
		return nil, ErrInputIsNil
	}
	v := reflect.ValueOf(input)

	if valueKey == nil {
		return nil, ErrKeyIsEmpty
	}
	key := reflect.ValueOf(valueKey)

	return getValueFromReflectValue(v, key)
}

func getValueFromReflectValue(v reflect.Value, key reflect.Value) (any, error) {
	switch v.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		// This is a primitive. Just return the value.
		return v.Interface(), nil
	case reflect.Array, reflect.Slice:
		// This is a list of values. Return the value at the specified index.
		indexType := reflect.TypeFor[int]()
		if !key.CanConvert(indexType) {
			return nil, ErrKeyTypeInvalid
		}
		index := key.Convert(indexType)
		i := index.Interface().(int)

		if i < 0 || i >= v.Len() {
			return nil, ErrValueNotFound
		}

		vi := v.Index(i)
		return vi.Interface(), nil // Return interface value
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// These values aren't introspectable in any way. Return an error.
		return nil, ErrImproperValueKind
	case reflect.Pointer, reflect.Interface:
		// These types hold another type as their underlying element. Use that instead, if they're not nil.
		if v.IsNil() {
			return nil, ErrInputIsNil
		}
		return getValueFromReflectValue(v.Elem(), key)
	case reflect.Map:
		// Return the value for the provided key.
		keyType := v.Type().Key()
		if !key.CanConvert(keyType) {
			return nil, ErrKeyTypeInvalid
		}

		val := v.MapIndex(key.Convert(keyType))
		if !val.IsValid() {
			return nil, ErrValueNotFound
		}
		return val.Interface(), nil
	case reflect.Struct:
		// Return the value for the field with the provided name.
		keyType := reflect.TypeFor[string]()
		if !key.CanConvert(keyType) {
			return nil, ErrKeyTypeInvalid
		}
		fieldName := key.Convert(keyType).Interface().(string)

		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, ErrValueNotFound // Field not found
		}
		if !field.CanInterface() {
			// Field is unexported
			return nil, ErrFieldIsUnexported
		}
		return field.Interface(), nil
	default:
		// This should only happen if a brand-new Kind is ever added.
		return nil, ErrImproperValueKind
	}
}

type StaticValueGetter[T any] struct {
	Value T
}

func NewStaticValueGetter[T any](v T) StaticValueGetter[T] {
	return StaticValueGetter[T]{Value: v}
}

func (s StaticValueGetter[T]) GetValue(_ any, _ any) (any, error) {
	return s.Value, nil
}

// TODO: Chainable value getter?
