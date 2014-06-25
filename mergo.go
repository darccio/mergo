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

// Errors reported by Mergo when it finds invalid arguments.
var (
	ErrNilArguments            = errors.New("src and dst must not be nil")
	ErrDifferentArgumentsTypes = errors.New("src and dst must be of same type")
	ErrNotSupported            = errors.New("only structs and maps are supported")
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
	if !src.IsValid() {
		return nil
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
	case reflect.Map:
		for _, key := range src.MapKeys() {
			srcElement := src.MapIndex(key)
			if !srcElement.IsValid() {
				continue
			}
			dstElement := dst.MapIndex(key)
			switch reflect.TypeOf(srcElement.Interface()).Kind() {
			case reflect.Struct:
				fallthrough
			case reflect.Map:
				if err := deepMerge(dstElement, srcElement, visited, depth+1); err != nil {
					return err
				}
			}
			if !dstElement.IsValid() {
				dst.SetMapIndex(key, srcElement)
			}
		}
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if src.IsNil() {
			break
		} else if dst.IsNil() {
			if dst.CanSet() && isEmptyValue(dst) {
				dst.Set(src)
			}
		} else if err := deepMerge(dst.Elem(), src.Elem(), visited, depth+1); err != nil {
			return err
		}
	default:
		if dst.CanSet() && !isEmptyValue(src) {
			dst.Set(src)
		}
	}
	return nil
}

// Merge sets fields' values in dst from src if they have a zero
// value of their type.
// dst and src must be valid same-type structs and dst must be
// a pointer to struct.
// It won't merge unexported (private) fields and will do recursively
// any exported field.
func Merge(dst interface{}, src interface{}) error {
	if dst == nil || src == nil {
		return ErrNilArguments
	}
	vDst := reflect.ValueOf(dst).Elem()
	vSrc := reflect.ValueOf(src)
	if vDst.Type() != vSrc.Type() {
		return ErrDifferentArgumentsTypes
	}
	if vDst.Kind() != reflect.Struct && vDst.Kind() != reflect.Map {
		return ErrNotSupported
	}
	return deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0)
}
