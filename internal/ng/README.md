# Mergo

## Why next gen?

1. Cleaner code.
2. Reduce `interface{}`/`any` usage in the API.
3. Allow the compiler to optimize the code through generics.
4. Reduce allocations: v1 does 4 allocations per merge.
5. Reduce `reflect` usage.
6. Migrate from sentinel errors to [concrete error types](https://jub0bs.com/posts/2025-03-31-why-concrete-error-types-are-superior-to-sentinel-errors/).
