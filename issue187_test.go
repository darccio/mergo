package mergo_test

import (
	"dario.cat/mergo"
	"testing"
)

func TestIssue187MergeStructToMap(t *testing.T) {
	dst := map[string]interface{}{
		"empty": "data",
	}

	src := struct {
		Foo   string
		Bar   int
		Empty string
	}{
		Foo: "hello",
		Bar: 42,
	}
	if err := mergo.Map(&dst, src); err != nil {
		t.Error(err)
	}
	if dst["foo"] != "hello" || dst["bar"] != 42 || dst["empty"] != "data" {
		t.Errorf("expected dst to be {foo: hello, bar: 42, empty: data}, got {foo: %v, bar: %v, empty: %v}", dst["foo"], dst["bar"], dst["empty"])
	}
}

func TestIssue187MergeStructToMapWithOverwrite(t *testing.T) {
	dst := map[string]interface{}{
		"foo":   "initial",
		"bar":   1,
		"empty": "data",
	}
	src := struct {
		Foo   string
		Bar   int
		Empty string
	}{
		Foo: "hello",
		Bar: 42,
	}
	if err := mergo.Map(&dst, src, mergo.WithOverride); err != nil {
		t.Error(err)
	}
	if dst["foo"] != "hello" || dst["bar"] != 42 || dst["empty"] != "data" {
		t.Errorf("expected dst to be {foo: hello, bar: 42, empty: data}, got {foo: %v, bar: %v, empty: %v}", dst["foo"], dst["bar"], dst["empty"])
	}
}

func TestIssue187MergeStructToMapWithOverwriteWithEmptyValue(t *testing.T) {
	dst := map[string]interface{}{
		"foo":   "initial",
		"bar":   1,
		"empty": "data",
	}
	src := struct {
		Foo   string
		Bar   int
		Empty string
	}{
		Foo: "hello",
		Bar: 42,
	}
	if err := mergo.Map(&dst, src, mergo.WithOverwriteWithEmptyValue); err != nil {
		t.Error(err)
	}
	if dst["foo"] != "hello" || dst["bar"] != 42 || dst["empty"] != "" {
		t.Errorf("expected dst to be {foo: hello, bar: 42, empty: }, got {foo: %v, bar: %v, empty: %v}", dst["foo"], dst["bar"], dst["empty"])
	}
}
