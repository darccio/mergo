package mergo

import "reflect"

// Merge (WIP) merges src into dst recursively setting src values on dst
// if src values are not zero values and dst values are zero values.
// Breaking change: src can't be a T pointer anymore.
//
//go:noinline
func Merge[T any](dst *T, src T) {
	if dst == nil {
		return
	}

	elm := reflect.ValueOf(*dst)

	// If dst is an interface, we need to get the underlying value.
	if elm.Kind() == reflect.Interface {
		elm = elm.Elem()
	}

	// If dst is a non-zero value, we don't need to do anything.
	if !elm.IsZero() {
		return
	}

	*dst = src
}
