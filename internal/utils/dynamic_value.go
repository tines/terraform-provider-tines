package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetUnderlyingDynamicValue converts a Terraform Dynamic value to a Go value.
// It handles all supported Terraform types and recursively converts nested structures.
//
// Supported conversions:
// - String -> Go string
// - Number -> Go int64 (for integers) or float64 (for decimals)
// - Bool -> Go bool
// - List/Tuple/Set -> Go []any (with recursive element conversion)
// - Object/Map -> Go map[string]any (with recursive value conversion)
// - Null -> Go nil
//
// Returns an error for Unknown values and unsupported types.
func GetUnderlyingDynamicValue(ctx context.Context, res *types.Dynamic) (any, diag.Diagnostics) {
	// Handle the underlying value in the dynamic type.
	var diags diag.Diagnostics

	switch val := res.UnderlyingValue().(type) {
	case types.String:
		tflog.Info(ctx, "Dynamic value is a string")
		return val.ValueString(), diags

	case types.Number:
		tflog.Info(ctx, "Dynamic value is a number")
		// Try to convert to int64 first, then fall back to float64
		if val.IsNull() {
			return nil, diags
		}
		bigFloat := val.ValueBigFloat()
		if bigFloat.IsInt() {
			intVal, _ := bigFloat.Int64()
			return intVal, diags
		} else {
			floatVal, _ := bigFloat.Float64()
			return floatVal, diags
		}

	case types.Bool:
		tflog.Info(ctx, "Dynamic value is a boolean")
		return val.ValueBool(), diags

	case types.List:
		tflog.Info(ctx, "Dynamic value is a list")
		var result []any
		for _, elem := range val.Elements() {
			elemDynamic := types.DynamicValue(elem)
			convertedElem, elemDiags := GetUnderlyingDynamicValue(ctx, &elemDynamic)
			diags.Append(elemDiags...)
			if diags.HasError() {
				return nil, diags
			}
			result = append(result, convertedElem)
		}
		return result, diags

	case types.Tuple:
		tflog.Info(ctx, "Dynamic value is a tuple")
		var result []any
		for _, elem := range val.Elements() {
			elemDynamic := types.DynamicValue(elem)
			convertedElem, elemDiags := GetUnderlyingDynamicValue(ctx, &elemDynamic)
			diags.Append(elemDiags...)
			if diags.HasError() {
				return nil, diags
			}
			result = append(result, convertedElem)
		}
		return result, diags

	case types.Object:
		tflog.Info(ctx, "Dynamic value is an object")
		result := make(map[string]any)
		for attrName, attrValue := range val.Attributes() {
			attrDynamic := types.DynamicValue(attrValue)
			convertedAttr, attrDiags := GetUnderlyingDynamicValue(ctx, &attrDynamic)
			diags.Append(attrDiags...)
			if diags.HasError() {
				return nil, diags
			}
			result[attrName] = convertedAttr
		}
		return result, diags

	case types.Map:
		tflog.Info(ctx, "Dynamic value is a map")
		result := make(map[string]any)
		for key, mapValue := range val.Elements() {
			mapDynamic := types.DynamicValue(mapValue)
			convertedValue, mapDiags := GetUnderlyingDynamicValue(ctx, &mapDynamic)
			diags.Append(mapDiags...)
			if diags.HasError() {
				return nil, diags
			}
			result[key] = convertedValue
		}
		return result, diags

	case types.Set:
		tflog.Info(ctx, "Dynamic value is a set")
		var result []any
		for _, elem := range val.Elements() {
			elemDynamic := types.DynamicValue(elem)
			convertedElem, elemDiags := GetUnderlyingDynamicValue(ctx, &elemDynamic)
			diags.Append(elemDiags...)
			if diags.HasError() {
				return nil, diags
			}
			result = append(result, convertedElem)
		}
		return result, diags
	}

	// Handle null/unknown values
	if res.IsNull() {
		tflog.Info(ctx, "Dynamic value is null")
		return nil, diags
	}

	if res.IsUnknown() {
		tflog.Info(ctx, "Dynamic value is unknown")
		diags.AddError("Dynamic value is unknown", "Cannot convert unknown values to Go types")
		return nil, diags
	}

	// Unsupported type
	tflog.Info(ctx, fmt.Sprintf("Dynamic value is an unsupported type: %T", res.UnderlyingValue()))
	diags.AddError("Dynamic value is an unsupported type",
		fmt.Sprintf("Type %T is not supported for conversion to Go values", res.UnderlyingValue()))
	return nil, diags
}

// SetUnderlyingDynamicValue converts a stringified JSON value to a Terraform Dynamic value
// with the correct underlying type based on the JSON structure.
//
// Supported conversions:
// - JSON string -> types.Dynamic with types.String underlying
// - JSON number -> types.Dynamic with types.Number underlying
// - JSON boolean -> types.Dynamic with types.Bool underlying
// - JSON array -> types.Dynamic with types.Tuple underlying (with recursive conversion)
// - JSON object -> types.Dynamic with types.Object underlying (with recursive conversion)
// - JSON null -> types.Dynamic with null value
// - Plain string (invalid JSON) -> types.Dynamic with types.String underlying (fallback)
//
// If the input is not valid JSON, it will be treated as a plain string value.
// Returns an error only for unsupported parsed JSON types.
func SetUnderlyingDynamicValue(ctx context.Context, jsonValue string) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Parse the JSON value
	var parsed any
	if err := json.Unmarshal([]byte(jsonValue), &parsed); err != nil {
		// If JSON parsing fails, treat the input as a plain string value
		tflog.Info(ctx, fmt.Sprintf("Input is not valid JSON, treating as plain string: %s", err.Error()))
		return types.DynamicValue(types.StringValue(jsonValue)), diags
	}

	// Convert the parsed value to a Terraform dynamic value
	dynamicValue, convertDiags := convertToTerraformValue(ctx, parsed)
	diags.Append(convertDiags...)

	return dynamicValue, diags
}

// convertToTerraformValue recursively converts a parsed JSON value to a Terraform dynamic value
func convertToTerraformValue(ctx context.Context, value any) (types.Dynamic, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch v := value.(type) {
	case nil:
		tflog.Info(ctx, "Converting null JSON value")
		return types.DynamicNull(), diags

	case string:
		tflog.Info(ctx, "Converting string JSON value")
		return types.DynamicValue(types.StringValue(v)), diags

	case bool:
		tflog.Info(ctx, "Converting boolean JSON value")
		return types.DynamicValue(types.BoolValue(v)), diags

	case float64:
		tflog.Info(ctx, "Converting number JSON value")
		// JSON numbers are always parsed as float64, but we need to preserve integer precision
		bigFloat := big.NewFloat(v)
		if bigFloat.IsInt() {
			// If it's an integer, create a Number value that preserves the integer nature
			intVal, _ := bigFloat.Int64()
			return types.DynamicValue(types.NumberValue(big.NewFloat(float64(intVal)))), diags
		} else {
			return types.DynamicValue(types.NumberValue(bigFloat)), diags
		}

	case []interface{}:
		tflog.Info(ctx, "Converting array JSON value")
		var tupleElements []attr.Value
		var tupleTypes []attr.Type

		for _, elem := range v {
			elemDynamic, elemDiags := convertToTerraformValue(ctx, elem)
			diags.Append(elemDiags...)
			if diags.HasError() {
				return types.DynamicNull(), diags
			}

			tupleElements = append(tupleElements, elemDynamic)
			tupleTypes = append(tupleTypes, elemDynamic.Type(ctx))
		}

		tupleType := types.TupleType{ElemTypes: tupleTypes}
		tupleValue, tupleDiags := types.TupleValue(tupleType.ElemTypes, tupleElements)
		diags.Append(tupleDiags...)
		if diags.HasError() {
			return types.DynamicNull(), diags
		}

		return types.DynamicValue(tupleValue), diags

	case map[string]any:
		tflog.Info(ctx, "Converting object JSON value")
		objectAttrs := make(map[string]attr.Value)
		objectTypes := make(map[string]attr.Type)

		for key, val := range v {
			attrDynamic, attrDiags := convertToTerraformValue(ctx, val)
			diags.Append(attrDiags...)
			if diags.HasError() {
				return types.DynamicNull(), diags
			}

			objectAttrs[key] = attrDynamic
			objectTypes[key] = attrDynamic.Type(ctx)
		}

		objectType := types.ObjectType{AttrTypes: objectTypes}
		objectValue, objectDiags := types.ObjectValue(objectType.AttrTypes, objectAttrs)
		diags.Append(objectDiags...)
		if diags.HasError() {
			return types.DynamicNull(), diags
		}

		return types.DynamicValue(objectValue), diags

	default:
		tflog.Error(ctx, fmt.Sprintf("Unsupported JSON value type: %T", value))
		diags.AddError("Unsupported JSON value type",
			fmt.Sprintf("JSON value of type %T is not supported for conversion to Terraform types", value))
		return types.DynamicNull(), diags
	}
}
