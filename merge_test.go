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

