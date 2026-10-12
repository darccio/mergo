package mergo_test

import (
	"reflect"
	"testing"

	"dario.cat/mergo"
)

type sliceStructIssue257 struct {
	Values []string
}

func TestIssue257WithOverride(t *testing.T) {
	dst := sliceStructIssue257{
		Values: []string{"value1", "value2", "value3"},
	}
	src := map[string]interface{}{
		"Values": []interface{}{"newValue1", "newValue2", "newValue3"},
	}

	err := mergo.Map(&dst, src, mergo.WithOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"newValue1", "newValue2", "newValue3"}
	if !reflect.DeepEqual(dst.Values, expected) {
		t.Errorf("expected %+v, got %+v", expected, dst.Values)
	}
}

func TestIssue257WithAppendSlice(t *testing.T) {
	dst := sliceStructIssue257{
		Values: []string{"value1", "value2", "value3"},
	}
	src := map[string]interface{}{
		"Values": []interface{}{"newValue1", "newValue2", "newValue3"},
	}

	err := mergo.Map(&dst, src, mergo.WithAppendSlice)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"value1", "value2", "value3", "newValue1", "newValue2", "newValue3"}
	if !reflect.DeepEqual(dst.Values, expected) {
		t.Errorf("expected %+v, got %+v", expected, dst.Values)
	}
}

func TestIssue257SafeAndLossyConversions(t *testing.T) {
	type customString string
	type customTarget struct {
		Strings []customString
		Ints    []int64
	}

	var dst customTarget
	src := map[string]interface{}{
		"Strings": []interface{}{"a", "b"},
		"Ints":    []interface{}{10, float64(20)},
	}

	if err := mergo.Map(&dst, src); err != nil {
		t.Fatalf("unexpected error on safe conversion: %v", err)
	}

	if !reflect.DeepEqual(dst.Strings, []customString{"a", "b"}) {
		t.Errorf("expected [a b], got %+v", dst.Strings)
	}
	if !reflect.DeepEqual(dst.Ints, []int64{10, 20}) {
		t.Errorf("expected [10 20], got %+v", dst.Ints)
	}

	lossySrc := map[string]interface{}{
		"Ints": []interface{}{float64(20.5)},
	}
	if err := mergo.Map(&dst, lossySrc, mergo.WithOverride); err == nil {
		t.Fatal("expected error for fractional float to int64, got nil")
	}
}
