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
