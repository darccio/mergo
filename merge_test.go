package mergo_test

import (
	"reflect"
	"testing"

	"dario.cat/mergo"
)

type transformer struct {
	m map[reflect.Type]func(dst, src reflect.Value) error
}

func (s *transformer) Transformer(t reflect.Type) func(dst, src reflect.Value) error {
	if fn, ok := s.m[t]; ok {
		return fn
	}
	return nil
}

type foo struct {
	Bar *bar
	s   string
}

type bar struct {
	s map[string]string
	i int
}

func TestMergeWithTransformerNilStruct(t *testing.T) {
	a := foo{s: "foo"}
	b := foo{Bar: &bar{i: 2, s: map[string]string{"foo": "bar"}}}

	if err := mergo.Merge(&a, &b, mergo.WithOverride, mergo.WithTransformers(&transformer{
		m: map[reflect.Type]func(dst, src reflect.Value) error{
			reflect.TypeOf(&bar{}): func(dst, src reflect.Value) error {
				// Do sthg with Elem
				t.Log(dst.Elem().FieldByName("i"))
				t.Log(src.Elem())
				return nil
			},
		},
	})); err != nil {
		t.Error(err)
	}

	if a.s != "foo" {
		t.Errorf("b not merged in properly: a.s.Value(%s) != expected(%s)", a.s, "foo")
	}

	if a.Bar == nil {
		t.Errorf("b not merged in properly: a.Bar shouldn't be nil")
	}
}

func TestMergeNonPointer(t *testing.T) {
	dst := bar{
		i: 1,
	}
	src := bar{
		i: 2,
		s: map[string]string{
			"a": "1",
		},
	}
	want := mergo.ErrNonPointerArgument

	if got := mergo.Merge(dst, src); got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

func TestMapNonPointer(t *testing.T) {
	dst := make(map[string]bar)
	src := map[string]bar{
		"a": {
			i: 2,
			s: map[string]string{
				"a": "1",
			},
		},
	}
	want := mergo.ErrNonPointerArgument
	if got := mergo.Merge(dst, src); got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

type structWithSlices struct {
	Values []string
}

func TestMergeInterfaceSliceIntoStringSlice(t *testing.T) {
	dst := structWithSlices{
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

func TestMergeInterfaceSliceAppendIntoStringSlice(t *testing.T) {
	dst := structWithSlices{
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

func TestMapMergeSliceIntoNonSliceWithoutOverridePreservesValue(t *testing.T) {
	dst := map[string]interface{}{
		"key": "existing",
	}
	src := map[string]interface{}{
		"key": []string{"a", "b"},
	}

	if err := mergo.Map(&dst, src); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst["key"] != "existing" {
		t.Errorf("expected existing value to be preserved, got %#v", dst["key"])
	}
}

func TestConvertSliceRejectsLossyConversions(t *testing.T) {
	type target struct {
		Strings []string
		Ints    []int
		Int8s   []int8
		Uints   []uint
		IntPtrs []*int
	}

	cases := []struct {
		name string
		src  map[string]interface{}
	}{
		{
			name: "numeric to string",
			src:  map[string]interface{}{"Strings": []interface{}{65}},
		},
		{
			name: "fractional float to int",
			src:  map[string]interface{}{"Ints": []interface{}{1.5}},
		},
		{
			name: "out of range float to int",
			src:  map[string]interface{}{"Ints": []interface{}{1e30}},
		},
		{
			name: "narrowing int overflow",
			src:  map[string]interface{}{"Int8s": []interface{}{int64(300)}},
		},
		{
			name: "negative int to uint",
			src:  map[string]interface{}{"Uints": []interface{}{-1}},
		},
		{
			name: "pointer element lossy float to int",
			src:  map[string]interface{}{"IntPtrs": []interface{}{2.5}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name+"_override", func(t *testing.T) {
			dst := target{
				Strings: []string{"init"},
				Ints:    []int{1},
				Int8s:   []int8{1},
				Uints:   []uint{1},
			}
			if err := mergo.Map(&dst, tc.src, mergo.WithOverride); err == nil {
				t.Fatalf("expected error for %s with override, got nil", tc.name)
			}
		})

		t.Run(tc.name+"_append", func(t *testing.T) {
			dst := target{
				Strings: []string{"init"},
				Ints:    []int{1},
				Int8s:   []int8{1},
				Uints:   []uint{1},
			}
			if err := mergo.Map(&dst, tc.src, mergo.WithAppendSlice); err == nil {
				t.Fatalf("expected error for %s with append, got nil", tc.name)
			}
		})
	}
}
