package mergo_test

import (
	"reflect"
	"testing"

	"dario.cat/mergo"
)

type simpleStruct233 struct {
	Field1 string
	Field2 string
	Field3 string
}

type structWithSlice233 struct {
	Items []simpleStruct233
}

// TestWithSliceDeepMerge verifies that WithSliceDeepMerge merges overlapping
// slice elements without overriding non-empty destination fields and appends
// trailing source elements.
func TestWithSliceDeepMerge(t *testing.T) {
	src := structWithSlice233{
		Items: []simpleStruct233{
			{Field1: "src0_f1", Field2: "src0_f2", Field3: ""},
			{Field1: "src1_f1", Field2: "src1_f2", Field3: ""},
		},
	}
	dst := structWithSlice233{
		Items: []simpleStruct233{
			{Field1: "dst0_f1", Field2: "", Field3: "dst0_f3"},
		},
	}
	expected := structWithSlice233{
		Items: []simpleStruct233{
			{Field1: "dst0_f1", Field2: "src0_f2", Field3: "dst0_f3"},
			{Field1: "src1_f1", Field2: "src1_f2", Field3: ""},
		},
	}

	if err := mergo.Merge(&dst, src, mergo.WithSliceDeepMerge); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(dst, expected) {
		t.Fatalf("expected %+v, got %+v", expected, dst)
	}
}

// TestWithSliceDeepMergeAndOverride verifies that WithSliceDeepMerge combined
// with WithOverride overrides overlapping non-empty destination fields while
// appending trailing source elements.
func TestWithSliceDeepMergeAndOverride(t *testing.T) {
	src := structWithSlice233{
		Items: []simpleStruct233{
			{Field1: "src0_f1", Field2: "src0_f2", Field3: ""},
			{Field1: "src1_f1", Field2: "src1_f2", Field3: "src1_f3"},
		},
	}
	dst := structWithSlice233{
		Items: []simpleStruct233{
			{Field1: "dst0_f1", Field2: "", Field3: "dst0_f3"},
		},
	}
	expected := structWithSlice233{
		Items: []simpleStruct233{
			{Field1: "src0_f1", Field2: "src0_f2", Field3: "dst0_f3"},
			{Field1: "src1_f1", Field2: "src1_f2", Field3: "src1_f3"},
		},
	}

	if err := mergo.Merge(&dst, src, mergo.WithSliceDeepMerge, mergo.WithOverride); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(dst, expected) {
		t.Fatalf("expected %+v, got %+v", expected, dst)
	}
}

// TestWithSliceDeepMergeInterfaceStructs verifies that WithSliceDeepMerge merges
// struct values stored inside map-held []interface{} slices.
func TestWithSliceDeepMergeInterfaceStructs(t *testing.T) {
	src := map[string]interface{}{
		"items": []interface{}{
			simpleStruct233{Field1: "src0_f1", Field2: "src0_f2", Field3: ""},
			simpleStruct233{Field1: "src1_f1", Field2: "src1_f2", Field3: ""},
		},
	}
	dst := map[string]interface{}{
		"items": []interface{}{
			simpleStruct233{Field1: "dst0_f1", Field2: "", Field3: "dst0_f3"},
		},
	}
	expected := map[string]interface{}{
		"items": []interface{}{
			simpleStruct233{Field1: "dst0_f1", Field2: "src0_f2", Field3: "dst0_f3"},
			simpleStruct233{Field1: "src1_f1", Field2: "src1_f2", Field3: ""},
		},
	}

	if err := mergo.Merge(&dst, src, mergo.WithSliceDeepMerge); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(dst, expected) {
		t.Fatalf("expected %+v, got %+v", expected, dst)
	}
}

// TestWithSliceDeepMergeDirectInterfaceStructs verifies that WithSliceDeepMerge
// fills empty fields in struct values stored in a direct []interface{} struct field
// and appends trailing source elements.
func TestWithSliceDeepMergeDirectInterfaceStructs(t *testing.T) {
	type holder struct {
		Items []interface{}
	}
	dst := holder{Items: []interface{}{
		simpleStruct233{Field1: "dst0_f1", Field2: "", Field3: "dst0_f3"},
		nil,
	}}
	src := holder{Items: []interface{}{
		simpleStruct233{Field1: "src0_f1", Field2: "src0_f2", Field3: ""},
		"filled_nil_slot",
		simpleStruct233{Field1: "src1_f1", Field2: "src1_f2", Field3: ""},
	}}
	expected := holder{Items: []interface{}{
		simpleStruct233{Field1: "dst0_f1", Field2: "src0_f2", Field3: "dst0_f3"},
		"filled_nil_slot",
		simpleStruct233{Field1: "src1_f1", Field2: "src1_f2", Field3: ""},
	}}

	if err := mergo.Merge(&dst, src, mergo.WithSliceDeepMerge); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(dst, expected) {
		t.Fatalf("expected %+v, got %+v", expected, dst)
	}
}

// TestWithSliceDeepMergeIncompatibleInterfaceStructError verifies that an incompatible
// source element in a []interface{} slice returns an error instead of panicking.
func TestWithSliceDeepMergeIncompatibleInterfaceStructError(t *testing.T) {
	type holder struct {
		Items []interface{}
	}
	dst := holder{Items: []interface{}{
		simpleStruct233{Field1: "dst0_f1"},
	}}
	src := holder{Items: []interface{}{
		"scalar_instead_of_struct",
	}}
	if err := mergo.Merge(&dst, src, mergo.WithSliceDeepMerge); err == nil {
		t.Fatal("expected error when merging scalar into struct inside []interface{}, got nil")
	}
}
