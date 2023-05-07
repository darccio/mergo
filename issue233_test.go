package mergo_test

import (
	"reflect"
	"testing"

	"github.com/imdario/mergo"
)

// SimpleStructTest233 is a simple struct with fields of a base type
type SimpleStructTest233 struct {
	Field1 string
	Field2 string
	Field3 string
}

// StructWithSliceOfSimpleStructsTest233 has to have slice of structs with fields
type StructWithSliceOfSimpleStructsTest233 struct {
	SliceOfSimpleStructs []SimpleStructTest233
}

// makeSrcDst makes source and destination structs for test
func makeSrcDst() (src StructWithSliceOfSimpleStructsTest233, dst StructWithSliceOfSimpleStructsTest233) {
	src = StructWithSliceOfSimpleStructsTest233{
		SliceOfSimpleStructs: []SimpleStructTest233{
			{
				Field1: "src:Slice[0].Field1",
				Field2: "src:Slice[0].Field2",
				Field3: "",
			},
			{
				Field1: "src:Slice[1].Field1",
				Field2: "src:Slice[1].Field2",
				Field3: "",
			},
		},
	}
	dst = StructWithSliceOfSimpleStructsTest233{
		SliceOfSimpleStructs: []SimpleStructTest233{
			{
				Field1: "dst:Slice[0].Field1",
				Field2: "",
				Field3: "dst:Slice[0].Field3",
			},
		},
	}
	return
}

// TestNestedStructsFieldsAreMergedWithDeepMerge test base mergo.WithSliceDeepMerge usage
func TestNestedStructsFieldsAreMergedWithDeepMerge(t *testing.T) {
	src, dst := makeSrcDst()
	expected := StructWithSliceOfSimpleStructsTest233{
		SliceOfSimpleStructs: []SimpleStructTest233{
			{
				// Original dst field is expected not to be overwritten by value
				Field1: "dst:Slice[0].Field1",
				// Empty dst field is expected to be filled with src value
				Field2: "src:Slice[0].Field2",
				// Original dst field is expected not to be overwritten by empty value
				Field3: "dst:Slice[0].Field3",
			},
			// Expected dst being appended
			{
				Field1: "src:Slice[1].Field1",
				Field2: "src:Slice[1].Field2",
				Field3: "",
			},
		},
	}

	err := mergo.Merge(&dst, src, mergo.WithSliceDeepMerge)
	if err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if !reflect.DeepEqual(dst, expected) {
		t.Errorf("expected: %#v\ngot: %#v", expected, dst)
	}
}

// TestNestedStructsFieldsAreMergedWithDeepMergeWithOverride test combination of
// mergo.WithSliceDeepMerge and mergo.WithOverride
func TestNestedStructsFieldsAreMergedWithDeepMergeWithOverride(t *testing.T) {
	src, dst := makeSrcDst()
	expected := StructWithSliceOfSimpleStructsTest233{
		SliceOfSimpleStructs: []SimpleStructTest233{
			{
				// Original dst field is expected to be overwritten by value
				Field1: "src:Slice[0].Field1",
				// Empty dst field is expected to be filled with src value
				Field2: "src:Slice[0].Field2",
				// Original dst field is expected not to be overwritten by empty value
				Field3: "dst:Slice[0].Field3",
			},
			// Expected dst being appended
			{
				Field1: "src:Slice[1].Field1",
				Field2: "src:Slice[1].Field2",
				Field3: "",
			},
		},
	}

	err := mergo.Merge(&dst, src, mergo.WithSliceDeepMerge, mergo.WithOverride)
	if err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if !reflect.DeepEqual(dst, expected) {
		t.Errorf("expected: %#v\ngot: %#v", expected, dst)
	}
}

// TestNestedStructsFieldsAreMergedWithDeepMergeWithOverwriteWithEmptyValue test combination of
// mergo.WithSliceDeepMerge and mergo.WithOverwriteWithEmptyValue
func TestNestedStructsFieldsAreMergedWithDeepMergeWithOverwriteWithEmptyValue(t *testing.T) {
	src, dst := makeSrcDst()
	expected := StructWithSliceOfSimpleStructsTest233{
		SliceOfSimpleStructs: []SimpleStructTest233{
			{
				// Original dst field is expected to be overwritten by value
				Field1: "src:Slice[0].Field1",
				// Empty dst field is expected to be filled with src value
				Field2: "src:Slice[0].Field2",
				// Original dst field is expected to be overwritten by empty value
				Field3: "",
			},
			// Expected dst being appended
			{
				Field1: "src:Slice[1].Field1",
				Field2: "src:Slice[1].Field2",
				Field3: "",
			},
		},
	}

	err := mergo.Merge(&dst, src, mergo.WithSliceDeepMerge, mergo.WithOverwriteWithEmptyValue)
	if err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if !reflect.DeepEqual(dst, expected) {
		t.Errorf("expected: %#v\ngot: %#v", expected, dst)
	}
}
