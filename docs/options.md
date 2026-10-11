<!--
title: Merge options
summary: Choose overwrite, empty-value, slice, pointer, type-check, and transformer behavior.
language: en
feed: false
-->
# Merge options

Pass options after `dst` and `src` to `mergo.Merge` or `mergo.Map`. Options configure a merge; they do not validate every possible input. All functions below are exported in [v1.0.2](https://pkg.go.dev/dario.cat/mergo@v1.0.2).

## At a glance

| Option | Effect | Important boundary |
| --- | --- | --- |
| `WithOverride` | Gives non-empty source scalar struct fields priority | Map entries and opaque structs have different rules |
| `WithOverwriteWithEmptyValue` | Also copies empty source values; enables override | Deletes destination map keys absent from the source |
| `WithOverrideEmptySlice` | Allows nil or empty source slices to be assigned | Needs `WithOverride` to clear a non-empty destination slice |
| `WithoutDereference` | Treats a non-nil pointer as non-empty | Does not replace an existing pointer to a struct |
| `WithAppendSlice` | Appends source slice elements after destination elements | No deduplication; takes precedence over slice replacement |
| `WithTypeCheck` | Rejects different concrete slice types when replacing map-held slices | Use with `WithOverride`; not a general type checker |
| `WithSliceDeepCopy` | Visits overlapping slice positions recursively; enables override | Not a deep clone; does not extend slices or assign ordinary value elements |
| `WithTransformers(t)` | Installs custom behavior for selected types | A transformer replaces normal merging for its matched type |

## Override and empty values

For ordinary scalar struct fields, `WithOverride` still ignores an empty source. Use `WithOverwriteWithEmptyValue` when a source `false`, `0`, empty string, or nil pointer is meant to clear the destination.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_emptyValues() {
	type Config struct {
		Enabled bool
		Retries int
	}
	src := Config{}
	dst := Config{Enabled: true, Retries: 3}
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	if err := mergo.Merge(&dst, src, mergo.WithOverwriteWithEmptyValue); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	// Output:
	// {true 3}
	// {false 0}
}
```

`WithOverwriteWithEmptyValue` already enables override, so adding `WithOverride` is unnecessary. On maps it also removes destination keys missing from the source, including in nested maps. Treat it as a potentially destructive operation when you only intend to update a few fields.

Map entries use different rules: `WithOverride` alone can assign a zero-valued source entry. It leaves destination-only keys in place, unlike `WithOverwriteWithEmptyValue`.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_mapOverwrite() {
	dst := map[string]int{"retries": 3, "keep": 1}
	src := map[string]int{"retries": 0}
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst["retries"], dst["keep"], len(dst))
	if err := mergo.Merge(&dst, src, mergo.WithOverwriteWithEmptyValue); err != nil {
		panic(err)
	}
	_, kept := dst["keep"]
	fmt.Println(dst["retries"], kept, len(dst))
	// Output:
	// 0 1 2
	// 0 false 1
}
```

Opaque structs such as `time.Time` are another exception: without a [transformer](transformers.md), `WithOverride` can replace them with a zero struct. Do not apply the scalar rule to every type.

## Empty slices

By default, a non-empty source slice can fill an empty destination slice, but an empty source does not clear a non-empty destination. `WithOverrideEmptySlice` permits assignment of an empty or nil source slice. It does not enable override by itself: combine it with `WithOverride` to clear a non-empty slice.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_emptySlice() {
	type Config struct{ Names []string }
	dst := Config{Names: []string{"local"}}
	src := Config{Names: []string{}}
	if err := mergo.Merge(&dst, src, mergo.WithOverrideEmptySlice); err != nil {
		panic(err)
	}
	fmt.Println(dst.Names)
	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithOverrideEmptySlice); err != nil {
		panic(err)
	}
	fmt.Println(len(dst.Names), dst.Names == nil)
	src.Names = nil
	if err := mergo.Merge(&dst, src, mergo.WithOverrideEmptySlice); err != nil {
		panic(err)
	}
	fmt.Println(dst.Names == nil)
	// Output:
	// [local]
	// 0 false
	// true
}
```

This also distinguishes an allocated empty slice from a nil slice, which can matter when encoding JSON. `WithOverwriteWithEmptyValue` already allows empty-slice assignment and clearing. These descriptions apply to top-level slices and slice fields. Map-held slices follow the map branch's rules: `WithOverride` alone can replace them with empty source slices. Check the behavior for your map's concrete value types rather than treating the struct-field rule as universal.

## Append slices

`WithAppendSlice` appends destination elements followed by source elements, preserving duplicates. It applies even when the destination is non-empty, and adding `WithOverride` does not turn it into replacement.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_appendSlice() {
	dst := []string{"local", "shared"}
	src := []string{"shared", "default"}
	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithAppendSlice); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	// Output: [local shared shared default]
}
```

Appending can reuse the destination backing array. Pointer, slice, and map elements can still share data; see [aliasing](faq.md#does-a-merge-copy-data-or-share-it).

## Traverse slice positions

Despite its name, `WithSliceDeepCopy` is an element-wise merge, not a deep-copy operation. It enables override for the whole merge and visits indices up to the shorter slice length. It neither appends source-only positions nor removes destination-only positions.

A useful case is a slice of non-nil pointers to structs: the existing destination objects can be modified without changing pointer identity. Empty scalar source fields still preserve destination fields.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_sliceDeepCopy() {
	type Item struct{ Name string }
	first := &Item{Name: "local"}
	dst := []*Item{first, {Name: "keep"}}
	src := []*Item{{Name: "source"}, {Name: ""}, {Name: "ignored"}}
	if err := mergo.Merge(&dst, src, mergo.WithSliceDeepCopy); err != nil {
		panic(err)
	}
	fmt.Println(len(dst), dst[0].Name, dst[1].Name)
	fmt.Println(dst[0] == first, dst[0] == src[0])
	values := []int{1}
	if err := mergo.Merge(&values, []int{2, 3}, mergo.WithSliceDeepCopy); err != nil {
		panic(err)
	}
	fmt.Println(values)
	// Output:
	// 2 source keep
	// true false
	// [1]
}
```

The last line shows a limitation: ordinary scalar slice elements are not assigned by this option in v1.0.2. Value-struct elements have the same addressability problem. Maps and pointer targets can be modified, but this option does not promise independent data or support arbitrary heterogeneous slices. Use ordinary slice replacement or explicit copying when that is the intended behavior. [The implementation](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/merge.go) defines these bounds.

If `WithAppendSlice` and `WithSliceDeepCopy` are both set, append wins for slices. Prefer choosing one slice policy per call.

## Pointer defaults and identity

By default, a non-nil pointer to a zero value is considered empty, and merging can mutate its target. `WithoutDereference` makes a non-nil pointer count as present, preserving a caller's explicit `false` or `0` during a default merge.

For scalar pointers, `WithOverride` plus `WithoutDereference` assigns the source pointer rather than modifying the old destination target. The two values then share that pointer.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_pointerDefaults() {
	type Config struct{ Enabled *bool }
	local, fallback := false, true
	dst, src := Config{Enabled: &local}, Config{Enabled: &fallback}
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(local, dst.Enabled == &local)
	local = false
	if err := mergo.Merge(&dst, src, mergo.WithoutDereference); err != nil {
		panic(err)
	}
	fmt.Println(local, dst.Enabled == &local)
	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithoutDereference); err != nil {
		panic(err)
	}
	fmt.Println(*dst.Enabled, dst.Enabled == src.Enabled, local)
	// Output:
	// true true
	// false true
	// true true false
}
```

There is an important boundary for struct pointers: `WithoutDereference` does not replace two existing non-nil pointers to structs, even with `WithOverride`. An ordinary override instead merges their exported fields while retaining destination pointer identity.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_pointerStruct() {
	type Item struct{ Name string }
	type Config struct{ Item *Item }
	dst := Config{Item: &Item{Name: "local"}}
	src := Config{Item: &Item{Name: "source"}}
	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithoutDereference); err != nil {
		panic(err)
	}
	fmt.Println(dst.Item.Name, dst.Item == src.Item)
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst.Item.Name, dst.Item == src.Item)
	// Output:
	// local false
	// source false
}
```

A nil destination pointer can be assigned the source pointer. A nil source does not clear an existing pointer unless `WithOverwriteWithEmptyValue` is enabled. See the [pointer regression tests](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/issue131_test.go) and [nil-source test](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/issue149_test.go).

## Check slice types in maps

Top-level `Merge` requires identical destination and source types regardless of options. A `map[string]interface{}`, however, can hold different concrete slice types under the same key. `WithOverride` plus `WithTypeCheck` rejects replacing, for example, `[]int` with `[]string`.

Imports: `fmt` and `dario.cat/mergo`.

```go
func ExampleMerge_typeCheck() {
	dst := map[string]interface{}{"items": []int{1}}
	src := map[string]interface{}{"items": []string{"one"}}
	err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTypeCheck)
	fmt.Println(err != nil)
	fmt.Printf("%T\n", dst["items"])
	// Output:
	// true
	// []int
}
```

This check is specific to the map-held slice replacement branch. It does not enforce a schema for all interface or map values, perform conversions, or make `WithSliceDeepCopy` safe for arbitrary mixed element types. Slice append checks concrete slice types even without `WithTypeCheck`. See [the v1.0.2 slice type tests](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/mergo_test.go).

## Transformers and precedence

Use `WithTransformers(yourTransformer)` to choose behavior for a type such as `time.Time`. A matched callback runs instead of normal traversal; `WithOverride` and the empty-value options do not automatically alter its policy. [The transformer cookbook](transformers.md) explains how to implement that policy and return errors.

The built-in options above set flags, so their order does not change these precedence rules. Multiple `WithTransformers` options replace the same configuration field: the last one wins. Combine type handlers in one implementation rather than expecting them to be chained. A custom option that writes `Config` fields can introduce its own ordering behavior.

[Back to Mergo](index.md) · [FAQ](faq.md)
