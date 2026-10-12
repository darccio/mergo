package mergo_test

import (
	"reflect"
	"testing"

	"dario.cat/mergo"
)

type sliceTarget struct {
	Values []string
}

type multiSliceTarget struct {
	Strings []string
	Ints    []int
	Floats  []float64
	Bools   []bool
}

type nestedSliceTarget struct {
	Inner sliceTarget
}

type itemStruct struct {
	Name  string
	Count int
}

type structSliceTarget struct {
	Items []itemStruct
}

func TestMapInterfaceSliceToStringSliceOverride(t *testing.T) {
	dst := sliceTarget{
		Values: []string{"value1", "value2", "value3"},
	}
	src := map[string]interface{}{
		"Values": []interface{}{"newValue1", "newValue2", "newValue3"},
	}

	if err := mergo.Map(&dst, src, mergo.WithOverride); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"newValue1", "newValue2", "newValue3"}
	if !reflect.DeepEqual(dst.Values, expected) {
		t.Errorf("expected %v, got %v", expected, dst.Values)
	}
}

func TestMapInterfaceSliceToStringSliceAppend(t *testing.T) {
	dst := sliceTarget{
		Values: []string{"value1", "value2", "value3"},
	}
	src := map[string]interface{}{
		"Values": []interface{}{"newValue1", "newValue2", "newValue3"},
	}

	if err := mergo.Map(&dst, src, mergo.WithAppendSlice); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"value1", "value2", "value3", "newValue1", "newValue2", "newValue3"}
	if !reflect.DeepEqual(dst.Values, expected) {
		t.Errorf("expected %v, got %v", expected, dst.Values)
	}
}

func TestMapInterfaceSliceToStringSliceEmptyDst(t *testing.T) {
	var dst sliceTarget
	src := map[string]interface{}{
		"Values": []interface{}{"a", "b", "c"},
	}

	if err := mergo.Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(dst.Values, expected) {
		t.Errorf("expected %v, got %v", expected, dst.Values)
	}
}

func TestMapInterfaceSliceToStringSliceNoOverridePreserves(t *testing.T) {
	dst := sliceTarget{
		Values: []string{"original"},
	}
	src := map[string]interface{}{
		"Values": []interface{}{"new"},
	}

	if err := mergo.Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"original"}
	if !reflect.DeepEqual(dst.Values, expected) {
		t.Errorf("expected %v, got %v", expected, dst.Values)
	}
}

func TestMapMultipleTypedSlices(t *testing.T) {
	var dst multiSliceTarget
	src := map[string]interface{}{
		"Strings": []interface{}{"alpha", "beta"},
		"Ints":    []interface{}{1, 2, 3},
		"Floats":  []interface{}{1.1, 2.2},
		"Bools":   []interface{}{true, false},
	}

	if err := mergo.Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(dst.Strings, []string{"alpha", "beta"}) {
		t.Errorf("Strings: expected %v, got %v", []string{"alpha", "beta"}, dst.Strings)
	}
	if !reflect.DeepEqual(dst.Ints, []int{1, 2, 3}) {
		t.Errorf("Ints: expected %v, got %v", []int{1, 2, 3}, dst.Ints)
	}
	if !reflect.DeepEqual(dst.Floats, []float64{1.1, 2.2}) {
		t.Errorf("Floats: expected %v, got %v", []float64{1.1, 2.2}, dst.Floats)
	}
	if !reflect.DeepEqual(dst.Bools, []bool{true, false}) {
		t.Errorf("Bools: expected %v, got %v", []bool{true, false}, dst.Bools)
	}
}

func TestMapNestedStructSlice(t *testing.T) {
	var dst nestedSliceTarget
	src := map[string]interface{}{
		"Inner": map[string]interface{}{
			"Values": []interface{}{"x", "y"},
		},
	}

	if err := mergo.Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"x", "y"}
	if !reflect.DeepEqual(dst.Inner.Values, expected) {
		t.Errorf("expected %v, got %v", expected, dst.Inner.Values)
	}
}

func TestMapInterfaceSliceOfMapsToStructSlice(t *testing.T) {
	var dst structSliceTarget
	src := map[string]interface{}{
		"Items": []interface{}{
			map[string]interface{}{
				"Name":  "Item1",
				"Count": 10,
			},
			map[string]interface{}{
				"Name":  "Item2",
				"Count": 20,
			},
		},
	}

	if err := mergo.Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []itemStruct{
		{Name: "Item1", Count: 10},
		{Name: "Item2", Count: 20},
	}
	if !reflect.DeepEqual(dst.Items, expected) {
		t.Errorf("expected %v, got %v", expected, dst.Items)
	}
}

func TestMapIncompatibleInterfaceSliceFails(t *testing.T) {
	var dst sliceTarget
	src := map[string]interface{}{
		"Values": []interface{}{true, false},
	}

	err := mergo.Map(&dst, src)
	if err == nil {
		t.Fatal("expected error converting bools to strings, got nil")
	}
}

func TestMapInterfaceSliceWithEmptyAndNil(t *testing.T) {
	dst := sliceTarget{
		Values: []string{"existing"},
	}
	src := map[string]interface{}{
		"Values": []interface{}{},
	}

	if err := mergo.Map(&dst, src, mergo.WithOverwriteWithEmptyValue); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dst.Values) != 0 {
		t.Errorf("expected empty slice, got %v", dst.Values)
	}
}
