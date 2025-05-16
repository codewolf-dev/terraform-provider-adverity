package adverity

import (
	"encoding/json"
	"reflect"
)

type Parameter struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// FlattenedMarshal flattens a slice of Parameter into the base struct.
func FlattenedMarshal(base interface{}, params *[]Parameter) ([]byte, error) {
	// Use reflection to get the underlying value
	v := reflect.ValueOf(base)

	// If it's a pointer, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Create a copy of the base struct as an anonymous struct (no methods attached) to avoid recursion
	c := reflect.New(v.Type()).Elem()
	c.Set(v)

	// Marshal base struct (excluding parameters)
	b, err := json.Marshal(c.Interface())
	if err != nil {
		return nil, err
	}

	// Unmarshal into map to merge with parameters
	merged := make(map[string]interface{})
	if err := json.Unmarshal(b, &merged); err != nil {
		return nil, err
	}

	// Add parameters to map
	if params != nil {
		for _, p := range *params {
			merged[p.Key] = p.Value
		}
	}

	return json.Marshal(merged)
}
