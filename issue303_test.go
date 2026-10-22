package mergo

import (
	"errors"
	"reflect"
	"testing"
)

type issue303Transformer struct {
	err error
}

func (t issue303Transformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ.Kind() != reflect.Int {
		return nil
	}
	return func(dst, src reflect.Value) error {
		return t.err
	}
}

func TestIssue303IsEmptyValuePointerWithoutDereference(t *testing.T) {
	value := 0
	if isEmptyValue(reflect.ValueOf(&value), false) {
		t.Fatal("non-nil pointer should not be empty when dereferencing is disabled")
	}
}

func TestIssue303MergePropagatesTransformerErrorsFromStructFields(t *testing.T) {
	want := errors.New("issue303 transformer error")
	dst := struct {
		Value int
	}{}
	src := struct {
		Value int
	}{Value: 1}

	got := Merge(&dst, src, WithTransformers(issue303Transformer{err: want}))
	if !errors.Is(got, want) {
		t.Fatalf("expected transformer error %v, got %v", want, got)
	}
}

func TestIssue303MapErrors(t *testing.T) {
	t.Run("non-pointer destination", func(t *testing.T) {
		dst := map[string]any{}
		if got := Map(dst, struct{}{}); got != ErrNonPointerArgument {
			t.Fatalf("expected ErrNonPointerArgument, got %v", got)
		}
	})

	t.Run("nil source", func(t *testing.T) {
		dst := map[string]any{}
		if got := Map(&dst, nil); got != ErrNilArguments {
			t.Fatalf("expected ErrNilArguments, got %v", got)
		}
	})

	t.Run("struct source into non-map destination", func(t *testing.T) {
		dst := []string{}
		if got := Map(&dst, struct{}{}); got != ErrExpectedMapAsDestination {
			t.Fatalf("expected ErrExpectedMapAsDestination, got %v", got)
		}
	})

	t.Run("map source into non-struct destination", func(t *testing.T) {
		dst := []string{}
		src := map[string]any{}
		if got := Map(&dst, src); got != ErrExpectedStructAsDestination {
			t.Fatalf("expected ErrExpectedStructAsDestination, got %v", got)
		}
	})

	t.Run("unsupported source", func(t *testing.T) {
		dst := map[string]any{}
		if got := Map(&dst, 1); got != ErrNotSupported {
			t.Fatalf("expected ErrNotSupported, got %v", got)
		}
	})
}

func TestIssue303ChangeInitialCaseEmptyString(t *testing.T) {
	if got := changeInitialCase("", func(r rune) rune { return r + 1 }); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestIssue303MapRecursesIntoNestedMap(t *testing.T) {
	type child struct {
		Name string
	}
	type parent struct {
		Child child
	}

	dst := parent{}
	src := map[string]any{
		"child": map[string]any{
			"name": "nested",
		},
	}

	if err := Map(&dst, src); err != nil {
		t.Fatalf("unexpected map error: %v", err)
	}
	if dst.Child.Name != "nested" {
		t.Fatalf("expected nested child to be mapped, got %#v", dst.Child)
	}
}
