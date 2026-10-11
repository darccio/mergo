<!--
title: Mergo
summary: Merge Go structs and maps with configuration defaults.
language: en
feed: false
-->
# Mergo

Mergo merges Go structs and maps of the same type. It fills zero-value destination fields from the source, which is useful for configuration defaults.

Mergo is stable and frozen for production use. New features are reserved for a future major version.

## Install

```sh
go get dario.cat/mergo
```

Import the library as `dario.cat/mergo`. Use the [API reference](https://pkg.go.dev/dario.cat/mergo) for function signatures and available options.

## Use

Pass a pointer to the destination and a source of the same type to `mergo.Merge`. The default merge preserves non-zero destination fields. `mergo.WithOverride` allows source values to replace them.

Only exported fields are merged. Struct values inside maps cannot be merged because Go reflection cannot address them. See the [repository README](https://github.com/darccio/mergo#usage) for usage examples and transformers.

## Project links

- [API reference](https://pkg.go.dev/dario.cat/mergo)
- [Source repository](https://github.com/darccio/mergo)
- [Releases](https://github.com/darccio/mergo/releases)
- [Report a bug](https://github.com/darccio/mergo/issues)
- [Sponsor maintenance](https://github.com/sponsors/darccio)

## License

Mergo and these docs are distributed under the [BSD 3-Clause license](LICENSE.txt).
