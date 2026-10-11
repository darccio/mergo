<!--
title: Transformer cookbook
summary: Merge time.Time with explicit overwrite rules and propagate transformer errors.
language: en
feed: false
-->
# Transformer cookbook

A transformer customizes merging for an exact Go type. Pass an implementation of `mergo.Transformers` to `mergo.WithTransformers`. Its `Transformer(reflect.Type)` method returns a callback for a type it handles, and `nil` for types that should use normal merging.

A matching callback replaces normal processing for that value, including traversal into its children. It owns the decision to keep, replace, append, or reject the source. The callback receives reflected destination and source values, but does not receive Mergo's `Config`. This behavior follows the [v1.0.2 dispatch code](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/merge.go).

## Merge time.Time using IsZero

`time.Time` has unexported fields and its own `IsZero` method. A default merge does not fill a zero `time.Time` field, while `WithOverride` can replace a non-zero time with a zero source. A transformer makes the intended rule explicit. [Issue #52](https://github.com/darccio/mergo/issues/52) and its [regression tests](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/issue52_test.go) document this distinction.

This transformer fills a zero destination with a non-zero source. When its own `Overwrite` field is true, it also replaces a non-zero destination. It always preserves the destination if the source is zero. It checks `CanSet` before assigning and compares the exact type, so named types based on `time.Time` need their own handler.

Imports for the following definition and example: `fmt`, `reflect`, `time`, and `dario.cat/mergo`.

```go
type docsTimeTransformer struct {
	Overwrite bool
}

func (t docsTimeTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ != reflect.TypeOf(time.Time{}) {
		return nil
	}
	return func(dst, src reflect.Value) error {
		if !dst.CanSet() {
			return nil
		}
		source := src.Interface().(time.Time)
		destination := dst.Interface().(time.Time)
		if !source.IsZero() && (destination.IsZero() || t.Overwrite) {
			dst.Set(src)
		}
		return nil
	}
}
```

```go
func ExampleMerge_timeTransformer() {
	type Snapshot struct{ Time time.Time }
	first := time.Date(2025, 5, 7, 12, 0, 0, 0, time.UTC)
	later := first.Add(time.Hour)
	var dst Snapshot
	if err := mergo.Merge(&dst, Snapshot{Time: first}, mergo.WithTransformers(docsTimeTransformer{})); err != nil {
		panic(err)
	}
	fmt.Println(dst.Time.Format(time.RFC3339))
	if err := mergo.Merge(&dst, Snapshot{Time: later}, mergo.WithOverride, mergo.WithTransformers(docsTimeTransformer{})); err != nil {
		panic(err)
	}
	fmt.Println(dst.Time.Equal(first))
	if err := mergo.Merge(&dst, Snapshot{Time: later}, mergo.WithTransformers(docsTimeTransformer{Overwrite: true})); err != nil {
		panic(err)
	}
	fmt.Println(dst.Time.Equal(later))
	if err := mergo.Merge(&dst, Snapshot{}, mergo.WithTransformers(docsTimeTransformer{Overwrite: true})); err != nil {
		panic(err)
	}
	fmt.Println(dst.Time.Equal(later))
	// Output:
	// 2025-05-07T12:00:00Z
	// true
	// true
	// true
}
```

The second merge passes `WithOverride`, but the callback's `Overwrite` is false, so the first timestamp survives. The third merge sets the callback's policy to true without requiring `WithOverride`. This is deliberate: ordinary options do not automatically override a transformer's decisions.

To allow a zero source to clear a time, change the transformer's condition deliberately. Do not expect `WithOverwriteWithEmptyValue` to change the callback's rule. Handle multiple special types in the same `Transformer` method; passing `WithTransformers` twice installs only the last implementation.

## Return an error

Callbacks can return errors, and `Merge` returns that error to the caller. This recipe defines a named port type so the transformer only validates that field type; returning `nil` for ordinary `int` avoids validating unrelated integer fields.

Imports for the following definition and example: `errors`, `fmt`, `reflect`, and `dario.cat/mergo`.

```go
type docsPort int

var errDocsPort = errors.New("port must be positive")

type docsPortTransformer struct{}

func (docsPortTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ != reflect.TypeOf(docsPort(0)) {
		return nil
	}
	return func(dst, src reflect.Value) error {
		if src.Int() <= 0 {
			return errDocsPort
		}
		if dst.CanSet() && dst.Int() == 0 {
			dst.Set(src)
		}
		return nil
	}
}
```

```go
func ExampleMerge_transformerError() {
	type Config struct {
		Name string
		Port docsPort
	}
	var dst Config
	src := Config{Name: "service", Port: -1}
	err := mergo.Merge(&dst, src, mergo.WithTransformers(docsPortTransformer{}))
	fmt.Println(errors.Is(err, errDocsPort), dst.Name, dst.Port)
	// Output: true service 0
}
```

The example also shows that merging is **not transactional**: `Name` was assigned before the port transformer rejected the source. Map iteration order is unspecified, so do not rely on a particular partial result from a failed map merge. Validate before merging, or merge into an independent copy and replace the original only after success when atomic changes are required. A shallow struct copy does not isolate pointers, maps, or slice backing arrays.

## Reflection boundaries

Match the exact type before calling type-specific reflection methods. Use `CanSet` before assigning. A value returned from a map lookup may not be addressable or settable, so a transformer does not automatically make struct map values mutable. See the [FAQ](faq.md#why-arent-struct-fields-inside-maps-merged).

Mergo skips transformer dispatch for a nil destination pointer, slice, map, or interface. Normal merging may assign a source value in that case. If you need uniform validation, validate the source before merging rather than depending on a callback running for every value.

[Back to Mergo](index.md) · [Options](options.md)
