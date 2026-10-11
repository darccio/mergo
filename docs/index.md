<!--
title: Mergo
summary: Merge Go structs and maps, choose options, and migrate to dario.cat/mergo.
language: en
feed: false
-->
# Mergo

Mergo fills empty destination values from a source of the same type. It is useful for layering configuration defaults over values an application already has. It mutates the destination; it is not a general deep-copy or validation library.

These guides describe **v1.0.2**. Mergo v1 is stable and frozen for production use; new features are reserved for a future major version. The [API reference](https://pkg.go.dev/dario.cat/mergo@v1.0.2) remains the reference for signatures.

## Install

```sh
go get dario.cat/mergo@v1.0.2
```

Import `dario.cat/mergo`, without a `/v1` suffix. Existing users of `github.com/imdario/mergo` should read the [migration guide](migration.md) before updating.

## Fill defaults, then override

Pass a pointer to the destination and a source of the same type to `mergo.Merge`. A pointer to that source type also works. By default, non-empty destination values win. `WithOverride` gives non-empty source values priority for ordinary scalar struct fields.

This runnable example imports `fmt` and `dario.cat/mergo`:

```go
func ExampleMerge_defaults() {
	type Config struct {
		Host string
		Port int
	}
	src := Config{Host: "default.example", Port: 8080}
	dst := Config{Host: "custom.example"}
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	// Output:
	// {custom.example 8080}
	// {default.example 8080}
}
```

An empty scalar is `false`, a numeric zero, or an empty string. Maps and slices with length zero count as empty, and pointers normally inherit the emptiness of their targets. This cannot distinguish an omitted `false` or `0` from one explicitly configured by the caller. See [pointer defaults and empty values](options.md).

## Guides

- [Options](options.md): all eight options, precedence, empty slices, pointer behavior, and map caveats.
- [Transformers](transformers.md): merge `time.Time` using `IsZero`, choose an overwrite policy, and return validation errors.
- [Migration](migration.md): move direct imports to `dario.cat/mergo`, or temporarily pin an indirect legacy dependency.
- [FAQ](faq.md): map values, aliases, JSON numbers, errors, and other boundaries.

## Merge and Map

`Merge` accepts same-type structs, maps, and slices. It recursively visits exported struct fields and maps, subject to the [limitations in the FAQ](faq.md). A merge does not delete destination map keys by default.

`Map` maps between a struct and `map[string]interface{}`. It changes the first letter of a key to match a field, or lowercases the first letter of a field to make a map key. It does not interpret JSON tags or convert number types. Start with an initialized map when mapping a struct to a map. [The mapping example and boundaries](faq.md#can-map-convert-json-numbers-or-use-json-tags) show both directions.

Use `Merge(..., WithOverride)` or `Map(..., WithOverride)` in new code. `MergeWithOverwrite` and `MapWithOverwrite` are deprecated convenience wrappers.

## Executable examples and sources

The Go code blocks are examples from [docs_example_test.go](https://github.com/darccio/mergo/blob/master/docs_example_test.go). Put them in a `_test.go` file with `package mergo_test` and the imports noted above each example, or run the repository examples:

```sh
go test ./... -run '^Example' -count=1
```

Behavior in these guides is checked against the [v1.0.2 implementation](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/merge.go) and [emptiness rules](https://github.com/darccio/mergo/blob/7b33b2b01026fbbbbfcfbb1ee2c9c0a5e0c9a9f7/mergo.go). Code examples use deterministic output instead of the current time or map iteration order.

## Project links

- [API reference](https://pkg.go.dev/dario.cat/mergo@v1.0.2)
- [Source repository](https://github.com/darccio/mergo)
- [Releases](https://github.com/darccio/mergo/releases)
- [Report a bug](https://github.com/darccio/mergo/issues)
- [Sponsor maintenance](https://github.com/sponsors/darccio)

## License

Mergo and these docs are distributed under the [BSD 3-Clause license](LICENSE.txt).
