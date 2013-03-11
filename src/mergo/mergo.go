// Copyright 2013 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on src/pkg/reflect/deepequal.go from official
// golang's stdlib.

package mergo

import (
	"errors"
	"reflect"
)

var (
	InvalidArgumentsErr        = errors.New("src and dst must be valid")
	NilArgumentsErr            = errors.New("src and dst must not be nil")
	DifferentArgumentsTypesErr = errors.New("src and dst must be of same type")
	OnlyStructSupportedErr     = errors.New("only structs are supported")
)

// During deepMerge, must keep track of checks that are
// in progress.  The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited are stored in a map indexed by 17 * a1 + a2;
type visit struct {
	ptr  uintptr
	typ  reflect.Type
	next *visit
}

// From src/pkg/encoding/json.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deepMerge(dst, src reflect.Value, visited map[uintptr]*visit, depth int) error {
	if !dst.IsValid() || !src.IsValid() {
		return InvalidArgumentsErr
	}
	if dst.CanAddr() {
		addr := dst.UnsafeAddr()
		h := 17 * addr
		seen := visited[h]
		typ := dst.Type()
		for p := seen; p != nil; p = p.next {
			if p.ptr == addr && p.typ == typ {
				return nil
			}
		}
		// Remember, remember...
		visited[h] = &visit{addr, typ, seen}
	}
	switch dst.Kind() {
	case reflect.Struct:
		for i, n := 0, dst.NumField(); i < n; i++ {
			if err := deepMerge(dst.Field(i), src.Field(i), visited, depth+1); err != nil {
				return err
			}
		}
	default:
		if dst.CanSet() && isEmptyValue(dst) {
			dst.Set(src)
		}
	}
	return nil
}

func Merge(dst interface{}, src interface{}) error {
	if dst == nil || src == nil {
		return NilArgumentsErr
	}
	vDst := reflect.ValueOf(dst).Elem()
	vSrc := reflect.ValueOf(src)
	if vDst.Type() != vSrc.Type() {
		return DifferentArgumentsTypesErr
	}
	if vDst.Kind() != reflect.Struct {
		return OnlyStructSupportedErr
	}
	return deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0)
}
