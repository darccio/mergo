// Copyright 2014 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on src/pkg/reflect/deepequal.go from official
// golang's stdlib.

package mergo

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deepMap(dst reflect.Value, src map[string]interface{}, visited map[uintptr]*visit, depth int) (err error) {
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
	zeroValue := reflect.Value{}
	for key := range src {
		srcValue := src[key]
		fieldName := capitalize(key)
		dstElement := dst.FieldByName(fieldName)
		if (dstElement == zeroValue) {
			// We discard it because the field doesn't exist.
			continue
		}
	        srcElement := reflect.ValueOf(srcValue)
		dstKind := reflect.TypeOf(dstElement.Interface()).Kind()
		srcKind := reflect.TypeOf(srcElement.Interface()).Kind()
		// TODO What happens if dstElements is a pointer and srcElement isn't?
		if srcKind == reflect.Ptr && dstKind != reflect.Ptr {
		        srcElement = srcElement.Elem()
			srcKind = reflect.TypeOf(srcElement.Interface()).Kind()
		}
		if !srcElement.IsValid() {
			continue
		}
		if srcKind == dstKind {
			if err = deepMerge(dstElement, srcElement, visited, depth+1); err != nil {
				return
			}
		} else {
			if srcKind == reflect.Map {
				if err = deepMap(dstElement, srcValue.(map[string]interface{}), visited, depth+1); err != nil {
					return
				}
			} else {
				return fmt.Errorf("Type mismatch on %s field: found %v, expected %v", fieldName, dstKind, srcKind)
			}
		}
	}
	return
}

// Map sets fields' values in dst from src.
// src must be a map with string keys, usually coming from a third
// party: HTTP request, database query result, etc.
// dst must be a valid pointer to struct.
// It won't merge unexported (private) fields and will do recursively
// any exported field.
// Missing fields in dst from src's keys will be skipped.
func Map(dst interface{}, src map[string]interface{}) error {
	if dst == nil || src == nil {
		return ErrNilArguments
	}
	vDst := reflect.ValueOf(dst).Elem()
	if vDst.Kind() != reflect.Struct {
		return ErrNotSupported
	}
	return deepMap(vDst, src, make(map[uintptr]*visit), 0)
}
