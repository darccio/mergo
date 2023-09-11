package mergo

import "reflect"

// Merge (WIP) merges src into dst recursively setting src values on dst
// if src values are not zero values and dst values are zero values.
func Merge[T any](dst *T, src T) {
	if dst == nil {
		return
	}

	vDst := reflect.ValueOf(dst)

	e := vDst.Elem()
	if !e.CanSet() {
		return
	}

	if !e.IsZero() {
		return
	}

	*dst = src
}
