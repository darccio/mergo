package mergo

import (
	"reflect"
)

// Copyright 2020 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on src/pkg/reflect/deepequal.go from official
// golang's stdlib.

// Merge recursively copies non-zero values from src to exported and embedded fields with zero value in dst.
// dst must be a pointer. It assumes that src and dst are compatible, defined as follows.
//
// Values are compatible when they have the same type, or it is possible to traverse them similarly.
// I.e., a map and a struct share traversal semantics as they are key-value data structures.
func Merge(dst, src interface{}) error {
	if dst == nil {
		return ErrNilDestination
	} else if src == nil {
		return ErrNilSource
	}
	dstValue, srcValue, err := prepare(dst, src)
	if err != nil {
		return err
	}
	return merge(dstValue, srcValue)
}

func prepare(dst, src interface{}) (reflect.Value, reflect.Value, error) {
	var (
		dstValue = reflect.ValueOf(dst)
		srcValue = reflect.ValueOf(src)
	)
	if dstValue.Kind() == reflect.Ptr {
		dstValue = dstValue.Elem()
	} else {
		return zero(dstValue), zero(srcValue), ErrNonPointerDestination
	}
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}
	if !compatible(dstValue, srcValue) {
		return zero(dstValue), zero(srcValue), ErrIncompatibleTypes
	}
	return dstValue, srcValue, nil
}

func compatible(dst, src reflect.Value) bool {
	if dst.Type() == src.Type() {
		return true
	}
	switch dst.Kind() {
	case reflect.Map:
		switch src.Kind() {
		case reflect.Map:
			if dst.Type().Key().Kind() != src.Type().Key().Kind() {
				return false
			}
			return compatible(dst.Elem(), src.Elem())
		case reflect.Struct:
			if dst.Type().Key().Kind() != reflect.String {
				return false
			}
			// TODO: if value type is equal -> true
			// TODO: if value type is interface{} -> true
		}
	case reflect.Struct:

	case reflect.Slice:
		if src.Kind() != reflect.Slice {
			return false
		}
		return compatible(dst.Elem(), src.Elem())
	}
	return false
}

func zero(value reflect.Value) reflect.Value {
	return reflect.Zero(value.Type())
}

func merge(dst, src reflect.Value) error {
	// TODO
	return nil
}
