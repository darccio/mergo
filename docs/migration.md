<!--
title: Migrate to dario.cat/mergo
summary: Move direct imports to the vanity module path or pin an indirect legacy dependency.
language: en
feed: false
-->
# Migrate to dario.cat/mergo

Since v1.0.0, Mergo declares its module path as `dario.cat/mergo`. The GitHub repository location and the module's import path are different things. New code should import the declared module path.

## Recognize the module-path mismatch

If an import or dependency still requires `github.com/imdario/mergo` but selects a v1 release, Go reports a module-path mismatch. [Issue #248](https://github.com/darccio/mergo/issues/248) records the transition failure. Its [dependency discussion](https://github.com/darccio/mergo/issues/248#issuecomment-1985399728) also explains why replacing the old path with the new one can fail when both are in the module graph.

A `go get` alone does not rewrite your source imports. Choose the direct or indirect dependency path below.

## Direct dependency: change your imports

For code you maintain, replace imports of `github.com/imdario/mergo` with `dario.cat/mergo`, install the new module, and tidy your dependencies. The public API remains available at the new path, but run your application's tests to verify the behavior you depend on.

From the application's module root:

```sh
rg -l --null 'github\.com/imdario/mergo' -g '*.go' -g '!vendor/**' \
  | xargs -0 -r sed -i 's#github.com/imdario/mergo#dario.cat/mergo#g'
go get dario.cat/mergo@v1.0.2
go mod tidy
go test ./...
```

The rewrite command uses ripgrep and GNU `sed`/`xargs`, as available on Linux. On other platforms, use the editor's replace-in-files command over your own `.go` files, excluding vendored dependencies.

After `go mod tidy`, `github.com/imdario/mergo` should disappear from `go.mod` unless another dependency still imports it. Check the result:

```sh
go list -m all
rg 'github\.com/imdario/mergo' go.mod go.sum
rg 'github\.com/imdario/mergo' -g '*.go' -g '!vendor/**'
```

The old and new module can temporarily coexist if a dependency still uses the old path. Inspect the dependency chain with `go mod why -m github.com/imdario/mergo` and update that dependency when it supports the vanity path.

Do not add `/v1` to the import. Go's major-version suffix convention applies from v2 onward. See [Go's module path rules](https://go.dev/ref/mod#module-path).

## Indirect dependency: pin the old path temporarily

If a dependency you cannot immediately change still imports the old path, retain that path and pin it to **v0.3.16**, the final release declaring `github.com/imdario/mergo`. Add this replacement to your application's `go.mod`:

```text
replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.16
```

Equivalently, from the module root:

```sh
go mod edit -replace=github.com/imdario/mergo=github.com/imdario/mergo@v0.3.16
go mod tidy
go test ./...
```

This keeps the legacy dependency on the legacy module. It does not upgrade that dependency to v1. Your own code may still use the `dario.cat/mergo` module at v1.0.2 alongside it. Avoid redirecting the old path to `dario.cat/mergo` with a replacement: when the graph also uses the vanity module directly, Go can reject one module version being used for two module paths.

A replacement in a library's `go.mod` does not propagate to applications that depend on that library. Put a temporary replacement in the main application's module, or its active workspace as appropriate. This follows [Go's replacement rules](https://go.dev/ref/mod#go-mod-file-replace).

Once the dependency updates its imports, remove the replacement and tidy again:

```sh
go mod edit -dropreplace=github.com/imdario/mergo
go mod tidy
go test ./...
```

## Verify resolution

`go list -m dario.cat/mergo@v1.0.2` confirms the selected version. To explicitly check resolution without a module proxy:

```sh
GOPROXY=direct go list -m dario.cat/mergo@v1.0.2
```

When Go performs direct discovery for the vanity path, it requests the site's `?go-get=1` endpoint and reads its `go-import` metadata to locate GitHub. The documentation body does not participate in that discovery. See [Go's repository discovery rules](https://go.dev/ref/mod#vcs-find).

[Back to Mergo](index.md) · [FAQ](faq.md)
