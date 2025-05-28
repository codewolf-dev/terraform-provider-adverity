package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-adverity/internal/adverity"
)

// ConvertValue handles values.
func ConvertValue(value attr.Value) (interface{}, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}
	switch v := value.(type) {
	case types.String:
		return v.ValueString(), nil
	case types.Bool:
		return v.ValueBool(), nil
	case types.Int64:
		return v.ValueInt64(), nil
	case types.Float64:
		return v.ValueFloat64(), nil
	case types.Number:
		return ConvertNumber(v)
	case types.List:
		return ConvertList(v)
	case types.Tuple:
		return ConvertTuple(v)
	case types.Map:
		return ConvertMap(v)
	case types.Set:
		return ConvertSet(v)
	case types.Object:
		return ConvertObject(v)
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}

// ConvertNumber handles numbers.
func ConvertNumber(n types.Number) (interface{}, error) {
	if n.IsNull() || n.IsUnknown() {
		return nil, nil
	}
	bigFloat := n.ValueBigFloat()
	if bigFloat.IsInt() {
		intVal, _ := bigFloat.Int64()
		return intVal, nil
	} else {
		floatVal, _ := bigFloat.Int64()
		return floatVal, nil
	}
}

// ConvertList handles lists.
func ConvertList(l types.List) ([]interface{}, error) {
	if l.IsNull() || l.IsUnknown() {
		return nil, nil
	}
	elements := make([]interface{}, 0, len(l.Elements()))
	for _, elem := range l.Elements() {
		converted, err := ConvertValue(elem)
		if err != nil {
			return nil, fmt.Errorf("failed to convert list element: %w", err)
		}
		elements = append(elements, converted)
	}
	return elements, nil
}

// ConvertSet handle sets.
func ConvertSet(s types.Set) ([]interface{}, error) {
	if s.IsNull() || s.IsUnknown() {
		return nil, nil
	}
	elements := make([]interface{}, 0, len(s.Elements()))
	for _, elem := range s.Elements() {
		converted, err := ConvertValue(elem)
		if err != nil {
			return nil, fmt.Errorf("failed to convert set element: %w", err)
		}
		elements = append(elements, converted)
	}
	return elements, nil
}

// ConvertTuple handle tuples.
func ConvertTuple(t types.Tuple) ([]interface{}, error) {
	if t.IsNull() || t.IsUnknown() {
		return nil, nil
	}
	elements := make([]interface{}, 0, len(t.Elements()))
	for _, elem := range t.Elements() {
		converted, err := ConvertValue(elem)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tuple element: %w", err)
		}
		elements = append(elements, converted)
	}
	return elements, nil
}

// ConvertMap handle maps.
func ConvertMap(m types.Map) (map[string]interface{}, error) {
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}
	elements := make(map[string]interface{})
	for k, v := range m.Elements() {
		converted, err := ConvertValue(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert map element %s: %w", k, err)
		}
		elements[k] = converted
	}
	return elements, nil
}

// ConvertObject handle objects.
func ConvertObject(o types.Object) (map[string]interface{}, error) {
	if o.IsNull() || o.IsUnknown() {
		return nil, nil
	}
	elements := make(map[string]interface{})
	for k, v := range o.Attributes() {
		converted, err := ConvertValue(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert object attribute %s: %w", k, err)
		}
		elements[k] = converted
	}
	return elements, nil
}

func ExpandParameters(params attr.Value, path path.Path, diags *diag.Diagnostics) []adverity.Parameter {
	if params.IsNull() || params.IsUnknown() {
		return nil
	}

	object, ok := params.(types.Object)
	if !ok {
		diags.AddAttributeError(
			path,
			"Parameters must be an object",
			fmt.Sprintf("Expected parameters to be an object, got: %T", params),
		)
		return nil
	}

	converted, err := ConvertValue(object)
	if err != nil {
		diags.AddAttributeError(
			path,
			"Invalid parameters",
			"Failed to convert parameters: "+err.Error(),
		)
		return nil
	}

	result := make([]adverity.Parameter, 0, len(object.Attributes()))
	parameters, _ := converted.(map[string]interface{}) // we know it is a map[string]interface{} since params is an object
	for k, v := range parameters {
		result = append(result, adverity.Parameter{Key: k, Value: v})
	}

	return result
}
