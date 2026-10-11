<!--
title: Mergo FAQ
summary: Understand map values, aliasing, empty values, mapping types, and merge errors.
language: en
feed: false
-->
# Mergo FAQ

These answers describe v1.0.2 and link to implementation or regression tests where the behavior is easy to misread. See [Options](options.md) for the full flag reference and [Transformers](transformers.md) for type-specific policies.

## Why does a default merge ignore false, zero, or an empty string?

Those values count as empty in ordinary scalar fields, so Mergo cannot distinguish an explicit zero from an omitted value. `WithOverride` still skips an empty scalar source field. `WithOverwriteWithEmptyValue` allows clearing fields, but also deletes map keys absent from the source. For a field where an explicit `false` or `0` means something different from absence, use a pointer and `WithoutDereference` during default merging. [The options examples](options.md#pointer-defaults-and-identity) demonstrate both target mutation and pointer assignment.

The emptiness rules are implemented in [mergo.go](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/mergo.go). Mergo does not automatically call an arbitrary type's `IsZero` method.

## Why aren't struct fields inside maps merged?

A struct returned by a map lookup is not addressable, so Mergo cannot set its fields through reflection. An existing struct map value is preserved by a default merge; `WithOverride` can replace the whole map entry instead of combining its fields. [Closed issue #90](https://github.com/darccio/mergo/issues/90) and [its regression test](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/issue90_test.go) illustrate the addressability boundary.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_mapStructValues() {
	type Item struct {
		Name string
		Port int
	}
	dst := map[string]Item{"service": {Name: "local"}}
	src := map[string]Item{"service": {Name: "source", Port: 8080}}
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(dst["service"])
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst["service"])
	// Output:
	// {local 0}
	// {source 8080}
}
```

A map of non-nil pointers to structs can expose addressable targets, but then it also introduces pointer sharing and mutation. For independent values, explicitly copy a map element to a local variable, merge that variable, then assign it back.

Nested maps can be merged recursively:

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_nestedMaps() {
	dst := map[string]map[string]int{"service": {"port": 8080}}
	src := map[string]map[string]int{"service": {"port": 9000, "retries": 3}}
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(dst["service"]["port"], dst["service"]["retries"])
	// Output: 8080 3
}
```

## Does a merge copy data or share it?

Mergo is a merge library, not a deep-cloning library. Assigning a slice copies its descriptor and can share its backing array. Assigning a pointer shares its target. Map values that contain maps, pointers, or slices can also share their referenced data.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_aliasing() {
	type Config struct {
		Names []string
		Item  *string
	}
	name := "source"
	src := Config{Names: []string{"source"}, Item: &name}
	var dst Config
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	dst.Names[0] = "changed"
	*dst.Item = "changed"
	fmt.Println(src.Names[0], name, dst.Item == src.Item)
	// Output: changed changed true
}
```

`WithAppendSlice` may reuse destination slice capacity, and pointer/map/slice elements remain references. `WithSliceDeepCopy` does not allocate an independent clone: it merges existing overlapping positions and may retain references. Use explicit copying when isolation is required, and synchronize concurrent access when data is shared.

## Why doesn't WithSliceDeepCopy change my slice of structs or integers?

In v1.0.2, the element-wise path rewraps interfaceable elements as reflected values. Ordinary scalar and value-struct elements then lack the settable destination required for assignment. Pointer targets and map contents can still be mutated. The option also visits only indices present in both slices; it never grows the destination. The [executable options example](options.md#traverse-slice-positions) shows the limitation alongside a working pointer-slice case.

Choose append, replacement, or an explicit application-specific loop according to the desired result. Do not rely on this option as a safe general merge for heterogeneous slices.

## Why is time.Time special?

`time.Time` has unexported fields. A default merge does not fill a zero time field, while `WithOverride` can copy a zero time over a non-zero time. Other opaque structs can have similar behavior. A [time transformer](transformers.md#merge-timetime-using-iszero) uses `IsZero` and makes the source-priority rule explicit. See [closed issue #52](https://github.com/darccio/mergo/issues/52).

## Can Map convert JSON numbers or use JSON tags?

`Map` uses Go field names, not JSON tags. It capitalizes the first letter of map keys to locate struct fields; struct-to-map mapping lowercases the first letter of field names. Nested struct values mapped into a map remain struct values; they are not recursively converted into maps.

No automatic numeric conversion takes place. For example, standard `encoding/json` decoding into `map[string]interface{}` produces `float64` numbers, which do not match an `int` or `uint16` struct field. Decode directly into a typed struct, or convert and validate values before mapping. [Closed issue #138](https://github.com/darccio/mergo/issues/138) and [its type-mismatch test](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/issue138_test.go) capture this case.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMap() {
	type Config struct {
		Host string
		Port int
	}
	dst := Config{Host: "local.example"}
	src := map[string]interface{}{"host": "default.example", "port": 8080}
	if err := mergo.Map(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	mapped := map[string]interface{}{}
	if err := mergo.Map(&mapped, dst); err != nil {
		panic(err)
	}
	fmt.Println(mapped["host"], mapped["port"])
	err := mergo.Map(&dst, map[string]interface{}{"port": float64(8080)})
	fmt.Println(err != nil)
	// Output:
	// {local.example 8080}
	// local.example 8080
	// true
}
```

Initialize a destination map before struct-to-map mapping. Map-to-struct handling is not identical to `Merge` in every corner case: for example, an explicit zero map value can clear a scalar field with `WithOverride`. Use typed inputs and test the behavior your application needs. The conversion paths are in [map.go](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/map.go).

## Can Mergo merge any type or private field?

Top-level `Merge` accepts same-type structs, maps, and slices. It does not accept a standalone scalar as the destination. The destination must be a pointer to a usable value. Use a value or a pointer to the same type for the source.

Normal field traversal does not merge unexported fields. An opaque struct may be assigned as a whole under override; that is distinct from editing its private fields. Transformers can replace an entire settable value but cannot bypass Go reflection's access rules. Mergo does not apply JSON tags, perform general conversions, or define array element-wise merging.

Avoid passing typed nil pointers as top-level arguments. Validate inputs before merging when their types or presence are not already known.

## What happens when a merge returns an error?

Check every returned error. Common argument errors include `ErrNonPointerArgument`, `ErrNilArguments`, `ErrDifferentArgumentsTypes`, and `ErrNotSupported`. `Map` can also report field-type mismatches; a transformer may return application-specific errors.

A failure does not roll back earlier assignments. [The transformer error example](transformers.md#return-an-error) demonstrates a partially changed destination. Validate first, or work on an independent copy if success must be atomic. A shallow copy is insufficient when fields reference maps, slices, or pointer targets.

## Will v1 gain new features?

Mergo v1 is stable and frozen. Report reproducible bugs with the Mergo version, Go version, input types, options, and a small test. New features are reserved for a future major version. Use [the issue tracker](https://github.com/darccio/mergo/issues) for library problems, and the [migration guide](migration.md) for old module-path errors.

[Back to Mergo](index.md) · [Options](options.md) · [Transformers](transformers.md)
