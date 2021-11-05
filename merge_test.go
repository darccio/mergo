package mergo_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/imdario/mergo"
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
	s   string
	Bar *bar
}

type bar struct {
	i int
	s map[string]string
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
	want := mergo.ErrNonPointerAgument

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
	want := mergo.ErrNonPointerAgument
	if got := mergo.Merge(dst, src); got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

type fromStringTransformer struct {
}

func (s *fromStringTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(net.IP{}) {
		return func(dst, src reflect.Value) error {
			strValue := src.Interface().(string)
			address := net.ParseIP(strValue)
			dst.Set(reflect.ValueOf(address))

			return nil
		}
	}

	return nil
}

func TestMergeWithTypeTransformer(t *testing.T) {
	type myStruct struct {
		Address net.IP
		Name    string
	}

	var dst myStruct

	src := map[string]interface{}{
		"name":    "some name",
		"address": "11.22.33.44",
	}

	if err := mergo.Map(&dst, src, mergo.WithTransformers(&fromStringTransformer{})); err != nil {
		t.Errorf("err should be nil")
	}

	if !dst.Address.Equal(net.ParseIP("11.22.33.44")) {
		t.Errorf("invalid IP")
	}

	if dst.Name != "some name" {
		t.Errorf("invalid name")
	}
}
