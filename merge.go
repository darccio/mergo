// Copyright 2013 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on src/pkg/reflect/deepequal.go from official
// golang's stdlib.

package mergo

import (
	"fmt"
	"reflect"
)

func hasExportedField(dst reflect.Value) (exported bool) {
	for i, n := 0, dst.NumField(); i < n; i++ {
		field := dst.Type().Field(i)
		if field.Anonymous && dst.Field(i).Kind() == reflect.Struct {
			exported = exported || hasExportedField(dst.Field(i))
		} else {
			exported = exported || isFieldExported(field)
		}
	}
	return
}

func isFieldExported(field reflect.StructField) bool {
	return isExportedComponent(field.Name, field.PkgPath)
}

func isTypeExported(v reflect.Type) bool {
	return isExportedComponent(v.Name(), v.PkgPath())
}

func isExportedComponent(name, pkgPath string) bool {
	if len(pkgPath) > 0 {
		return false
	}
	c := name[0]
	if 'a' <= c && c <= 'z' || c == '_' {
		return false
	}
	return true
}

type Config struct {
	Overwrite                    bool
	AppendSlice                  bool
	TypeCheck                    bool
	Transformers                 Transformers
	overwriteWithEmptyValue      bool
	overwriteSliceWithEmptyValue bool
}

type Transformers interface {
	Transformer(reflect.Type) func(dst, src reflect.Value) error
}

// Traverses recursively both values, assigning src's fields values to dst.
// The map argument tracks comparisons that have already been seen, which allows
// short circuiting on recursive types.
func deepMerge(dstIn, src reflect.Value, visited map[uintptr]*visit, depth int, config *Config) (dst reflect.Value, err error) {
	dst = dstIn
	overwrite := config.Overwrite
	typeCheck := config.TypeCheck
	overwriteWithEmptySrc := config.overwriteWithEmptyValue
	overwriteSliceWithEmptySrc := config.overwriteSliceWithEmptyValue
	config.overwriteWithEmptyValue = false

	if !src.IsValid() {
		return
	}
	if dst.CanAddr() {
		addr := dst.UnsafeAddr()
		h := 17 * addr
		seen := visited[h]
		typ := dst.Type()
		for p := seen; p != nil; p = p.next {
			if p.ptr == addr && p.typ == typ {
				return dst, nil
			}
		}
		// Remember, remember...
		visited[h] = &visit{addr, typ, seen}
	}

	if config.Transformers != nil && !isEmptyValue(dst) {
		if fn := config.Transformers.Transformer(dst.Type()); fn != nil {
			err = fn(dst, src)
			return
		}
	}

	switch dst.Kind() {
	case reflect.Struct:
		if hasExportedField(dst) {
			dstCp := reflect.New(dst.Type()).Elem()
			for i, n := 0, dst.NumField(); i < n; i++ {
				dstField := dst.Field(i)
				if !isFieldExported(dst.Type().Field(i)) {
					continue
				}
				if dst.Field(i).IsValid() {
					k := dstField.Interface()
					dstField = reflect.ValueOf(k)
				}
				dstField, err = deepMerge(dstField, src.Field(i), visited, depth+1, config)
				if err != nil {
					return
				}
				dstCp.Field(i).Set(dstField)
			}
			if dst.CanSet() {
				dst.Set(dstCp)
			} else {
				dst = dstCp
			}

		} else {
			if dst.CanSet() && (!isEmptyValue(src) || overwriteWithEmptySrc) && (overwrite || isEmptyValue(dst)) {
				dst.Set(src)
			}
		}
	case reflect.Map:
		if dst.IsNil() && !src.IsNil() {
			if dst.CanSet() {
				dst.Set(reflect.MakeMap(dst.Type()))
			} else {
				dst = src
				return
			}
		}
		for _, key := range src.MapKeys() {
			srcElement := src.MapIndex(key)
			if !srcElement.IsValid() {
				continue
			}
			dstElement := dst.MapIndex(key)
			if dst.MapIndex(key).IsValid() {
				k := dstElement.Interface()
				dstElement = reflect.ValueOf(k)
			}
			// dstElement.Set(reflect.ValueOf())

			switch srcElement.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Interface, reflect.Slice:
				if srcElement.IsNil() {
					continue
				}
				fallthrough
			default:
				if !srcElement.CanInterface() {
					continue
				}
				switch reflect.TypeOf(srcElement.Interface()).Kind() {
				case reflect.Struct:
					fallthrough
				case reflect.Ptr:
					fallthrough
				case reflect.Map:
					srcMapElm := srcElement
					dstMapElm := dstElement
					if srcMapElm.CanInterface() {
						srcMapElm = reflect.ValueOf(srcMapElm.Interface())
						if dstMapElm.IsValid() {
							dstMapElm = reflect.ValueOf(dstMapElm.Interface())
						}
					}
					dstMapElm, err = deepMerge(dstMapElm, srcMapElm, visited, depth+1, config)
					if err != nil {
						return
					}
					dst.SetMapIndex(key, dstMapElm)
				case reflect.Slice:
					srcSlice := reflect.ValueOf(srcElement.Interface())

					var dstSlice reflect.Value
					if !dstElement.IsValid() || dstElement.IsNil() {
						dstSlice = reflect.MakeSlice(srcSlice.Type(), 0, srcSlice.Len())
					} else {
						dstSlice = reflect.ValueOf(dstElement.Interface())
					}

					if (!isEmptyValue(src) || overwriteWithEmptySrc || overwriteSliceWithEmptySrc) && (overwrite || isEmptyValue(dst)) && !config.AppendSlice {
						if typeCheck && srcSlice.Type() != dstSlice.Type() {
							return fmt.Errorf("cannot override two slices with different type (%s, %s)", srcSlice.Type(), dstSlice.Type())
						}
						dstSlice = srcSlice
					} else if config.AppendSlice {
						if srcSlice.Type() != dstSlice.Type() {
							return fmt.Errorf("cannot append two slices with different type (%s, %s)", srcSlice.Type(), dstSlice.Type())
						}
						dstSlice = reflect.AppendSlice(dstSlice, srcSlice)
					}
					dst.SetMapIndex(key, dstSlice)
				}
			}
			if dstElement.IsValid() && !isEmptyValue(dstElement) && (reflect.TypeOf(srcElement.Interface()).Kind() == reflect.Map || reflect.TypeOf(srcElement.Interface()).Kind() == reflect.Slice) || (reflect.TypeOf(srcElement.Interface()).Kind() == reflect.Struct) {
				continue
			}

			if srcElement.IsValid() && ((srcElement.Kind() != reflect.Ptr && overwrite) || !dstElement.IsValid() || isEmptyValue(dstElement)) {
				if dst.IsNil() {
					dst.Set(reflect.MakeMap(dst.Type()))
				}
				dst.SetMapIndex(key, srcElement)
			}
		}
	case reflect.Slice:
		if !dst.CanSet() {
			break
		}
		if (!isEmptyValue(src) || overwriteWithEmptySrc || overwriteSliceWithEmptySrc) && (overwrite || isEmptyValue(dst)) && !config.AppendSlice {
			dst.Set(src)
		} else if config.AppendSlice {
			if src.Type() != dst.Type() {
				err = fmt.Errorf("cannot append two slice with different type (%s, %s)", src.Type(), dst.Type())
				return
			}
			dst.Set(reflect.AppendSlice(dst, src))
		}
	case reflect.Ptr, reflect.Interface:
		if isReflectNil(src) {
			break
		}

		if dst.Kind() != reflect.Ptr && src.Type().AssignableTo(dst.Type()) {
			if dst.IsNil() || overwrite {
				if overwrite || isEmptyValue(dst) {
					if dst.CanSet() {
						dst.Set(src)
					} else {
						dst = src
					}
				}
<<<<<<< HEAD
			}
			break
		}

		if src.Kind() != reflect.Interface {
			if dst.IsNil() || (src.Kind() != reflect.Ptr && overwrite) {
				if dst.CanSet() && (overwrite || isEmptyValue(dst)) {
					dst.Set(src)
				}
=======

>>>>>>> semi working
			} else if src.Kind() == reflect.Ptr {
				if dst, err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
					return
				}
			} else if dst.Elem().Type() == src.Type() {
				if dst, err = deepMerge(dst.Elem(), src, visited, depth+1, config); err != nil {
					return
				}
			} else {
				return dst, ErrDifferentArgumentsTypes
			}
			break
		}
		if dst.IsNil() || overwrite {
			if dst.CanSet() && (overwrite || isEmptyValue(dst)) {
				dst.Set(src)
			}
			// TODO HERE
		} else if _, err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
			return
		}
		// dst.Set()
	default:
		overwrite := (!isEmptyValue(src) || overwriteWithEmptySrc) && (overwrite || isEmptyValue(dst))
		if dst.CanSet() && overwrite {
			dst.Set(src)
		} else {
			dst = src
		}
	}

	return
}

// Merge will fill any empty for value type attributes on the dst struct using corresponding
// src attributes if they themselves are not empty. dst and src must be valid same-type structs
// and dst must be a pointer to struct.
// It won't merge unexported (private) fields and will do recursively any exported field.
func Merge(dst, src interface{}, opts ...func(*Config)) error {
	return merge(dst, src, opts...)
}

// MergeWithOverwrite will do the same as Merge except that non-empty dst attributes will be overridden by
// non-empty src attribute values.
// Deprecated: use Merge(…) with WithOverride
func MergeWithOverwrite(dst, src interface{}, opts ...func(*Config)) error {
	return merge(dst, src, append(opts, WithOverride)...)
}

// WithTransformers adds transformers to merge, allowing to customize the merging of some types.
func WithTransformers(transformers Transformers) func(*Config) {
	return func(config *Config) {
		config.Transformers = transformers
	}
}

// WithOverride will make merge override non-empty dst attributes with non-empty src attributes values.
func WithOverride(config *Config) {
	config.Overwrite = true
}

// WithOverride will make merge override empty dst slice with empty src slice.
func WithOverrideEmptySlice(config *Config) {
	config.overwriteSliceWithEmptyValue = true
}

// WithAppendSlice will make merge append slices instead of overwriting it.
func WithAppendSlice(config *Config) {
	config.AppendSlice = true
}

// WithTypeCheck will make merge check types while overwriting it (must be used with WithOverride).
func WithTypeCheck(config *Config) {
	config.TypeCheck = true
}

func merge(dst, src interface{}, opts ...func(*Config)) error {
	var (
		vDst, vSrc reflect.Value
		err        error
	)

	config := &Config{}

	for _, opt := range opts {
		opt(config)
	}

	if vDst, vSrc, err = resolveValues(dst, src); err != nil {
		return err
	}
	if !vDst.CanSet() {
		return fmt.Errorf("cannot set dst, needs reference")
	}
	if vDst.Type() != vSrc.Type() {
		return ErrDifferentArgumentsTypes
	}
	_, err = deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0, config)
	return err
}

func handleNil(dst reflect.Value) reflect.Value {
	if !dst.CanSet() {
		t := reflect.Indirect(reflect.ValueOf(dst)).Type()
		dsttmp := reflect.New(t).Elem()
		fmt.Println(dsttmp)
		initiallize(t, dsttmp)
		if !isTypeExported(t) {
			return dsttmp
		}
		dsttmp.Set(dst)
		dst = dsttmp
	}
	return dst
}

func initiallize(t reflect.Type, v reflect.Value) {
	// fmt.Println("initializeStruct(1.2.1)", t)
	switch t.Kind() {
	case reflect.Map:
		v.Set(reflect.MakeMap(t))
	case reflect.Slice:
		v.Set(reflect.MakeSlice(t, 0, 0))
	case reflect.Chan:
		v.Set(reflect.MakeChan(t, 0))
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			ft := t.Field(i)
			initiallize(ft.Type, f)
		}
	case reflect.Ptr:
		ft := t.Elem()
		fv := reflect.New(ft)
		initiallize(ft, fv.Elem())
		if isTypeExported(ft) {
			v.Set(fv)
		}
	default:
	}
}

func handleEmpty(dst, src reflect.Value, overwrite, appendSlice, overwriteWithEmptySrc bool) reflect.Value {
	dst = handleNil(dst)
	overwrite = (overwrite && !appendSlice || isEmptyValue(dst))
	if (!isEmptyValue(src) || overwriteWithEmptySrc) && overwrite {
		if dst.CanSet() {
			dst.Set(src)
		} else {
			dst = src
		}
	}
	if dst.CanInterface() {
		dst = reflect.ValueOf(dst.Interface())
	}
	return dst
}

// IsReflectNil is the reflect value provided nil
func isReflectNil(v reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.Interface, reflect.Slice, reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr:
		// Both interface and slice are nil if first word is 0.
		// Both are always bigger than a word; assume flagIndir.
		return v.IsNil()
	default:
		return false
	}
}
