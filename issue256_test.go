package mergo_test

import (
	"reflect"
	"testing"

	"dario.cat/mergo"
)

func TestIssue256MergeMapCaseInsensitively(t *testing.T) {
	dst := map[string]map[string]string{
		"Config": {
			"Host": "localhost",
		},
	}
	src := map[string]map[string]string{
		"config": {
			"host": "example.com",
			"Port": "443",
		},
	}

	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithCaseInsensitiveMapKeys); err != nil {
		t.Fatal(err)
	}

	expected := map[string]map[string]string{
		"Config": {
			"Host": "example.com",
			"Port": "443",
		},
	}
	if !reflect.DeepEqual(dst, expected) {
		t.Fatalf("got %#v, want %#v", dst, expected)
	}
}

func TestIssue256MergeMapCaseInsensitiveWithOverwriteEmpty(t *testing.T) {
	dst := map[string]string{"Host": "localhost", "User": "root"}
	src := map[string]string{"host": "example.com"}

	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithOverwriteWithEmptyValue, mergo.WithCaseInsensitiveMapKeys); err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{"Host": "example.com"}
	if !reflect.DeepEqual(dst, expected) {
		t.Fatalf("got %#v, want %#v", dst, expected)
	}
}

func TestIssue256MergeMapCaseSensitiveByDefault(t *testing.T) {
	dst := map[string]string{"Host": "localhost"}
	src := map[string]string{"host": "example.com"}

	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{"Host": "localhost", "host": "example.com"}
	if !reflect.DeepEqual(dst, expected) {
		t.Fatalf("got %#v, want %#v", dst, expected)
	}
}
