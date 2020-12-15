package cty

import "github.com/zclconf/go-cty/cty"

func ValueAsString(v cty.Value) string {
	return v.AsString()
}

func ValueAsBool(v cty.Value) bool {
	return v.True()
}

func ValueAsInt64(v cty.Value) int64 {
	r, _ := v.AsBigFloat().Int64()
	return r
}

func ValueAsObject(v cty.Value) map[string]cty.Value {
	return v.AsValueMap()
}

func ValueAsMap(v cty.Value) map[string]cty.Value {
	return v.AsValueMap()
}

func ValueAsSet(v cty.Value) []cty.Value {
	return v.AsValueSet().Values()
}

func ValueAsList(v cty.Value) []cty.Value {
	return v.AsValueSlice()
}
