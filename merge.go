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
	"unsafe"
)

func isMergeableField(dst reflect.Value) (exported bool) {
	for i, n := 0, dst.NumField(); i < n; i++ {
		field := dst.Type().Field(i)
		if field.Anonymous && dst.Field(i).Kind() == reflect.Struct {
			exported = exported || isMergeableField(dst.Field(i))
		} else if isExportedComponent(&field) {
			exported = exported || len(field.PkgPath) == 0
		}
	}

	return false
}

func isExportedComponent(field *reflect.StructField) bool {
	pkgPath := field.PkgPath
	if len(pkgPath) > 0 {
		return false
	}
	c := field.Name[0]
	if 'a' <= c && c <= 'z' || c == '_' {
		return false
	}
	return true
}

type Config struct {
	Overwrite                    bool
	AppendSlice                  bool
	TypeCheck                    bool
	overwriteWithEmptyValue      bool
	overwriteSliceWithEmptyValue bool
	sliceDeepCopy                bool
}

// Transformers allow to merge specific types differently than in the default behavior.
// In other words, now you can customize how some types are merged.
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
	sliceDeepCopy := config.sliceDeepCopy

	if !src.IsValid() {
		return
	}

	if dst.CanAddr() {
		addr := dst.UnsafeAddr()
		h := 17 * addr
		typ := dst.Type()
		seen := visited[h]

		for p := seen; p != nil; p = p.next {
			if p.ptr == addr && p.typ == typ {
				return
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

	if typeCheck && dst.IsValid() && src.Type() != dst.Type() {
		err = fmt.Errorf("cannot append two different types (%s, %s)", src.Kind(), dst.Kind())
		return
	}

	switch dst.Kind() {
	case reflect.Struct:
		if isMergeableField(dst) {
			for i, n := 0, dst.NumField(); i < n; i++ {
				dstField := dst.Field(i)
				structField := dst.Type().Field(i)
				// copy un-exported struct fields
				if !isExportedComponent(&structField) {
					rf := dstCp.Field(i)
					rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem() //nolint:gosec

					dstRF := dst.Field(i)
					if !dstRF.CanAddr() {
						continue
					}

					dstRF = reflect.NewAt(dstRF.Type(), unsafe.Pointer(dstRF.UnsafeAddr())).Elem() //nolint:gosec
					rf.Set(dstRF)

					continue
				}

				dstField, err = deepMerge(dstField, src.Field(i), visited, depth+1, config)
				if err != nil {
					return
				}

				dstCp.Field(i).Set(dstField)
			}
		} else {
			if (isReflectNil(dst) || overwrite) && (!isEmptyValue(src) || overwriteWithEmptySrc) {
				dst.Set(src)
			}

			return
		}

		if (isReflectNil(dst) || overwrite) && (!isEmptyValue(src) || overwriteWithEmptySrc) {
			dst = src
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
			if dstElement.IsValid() {
				k := dstElement.Interface()
				dstElement = reflect.ValueOf(k)
			}

			switch srcElement.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Interface, reflect.Slice:
				if isReflectNil(srcElement) {
					if overwrite || isReflectNil(dstElement) {
						dst.SetMapIndex(key, srcElement)
					}
					continue
				}
			}

			if !srcElement.CanInterface() {
				continue
			}

			if srcElement.CanInterface() {
				srcElement = reflect.ValueOf(srcElement.Interface())

				if dstElement.IsValid() {
					dstElement = reflect.ValueOf(dstElement.Interface())
				}
			}

			if dstElement.IsValid() && dstElement.Type() == srcElement.Type() {
				newSlice := srcElement
				if overwrite || overwriteSliceWithEmptySrc || overwriteWithEmptySrc {
					if dstElement.Kind() == reflect.Slice && config.AppendSlice {
						if dstElement.Type().Elem().Kind() == srcElement.Type().Elem().Kind() {
							newSlice = reflect.AppendSlice(dstElement, srcElement)
						} else {
							err = fmt.Errorf("cannot override two slices with different type (%s, %s)", src.Type(), dst.Type())
							return
						}
					}
					dst.SetMapIndex(key, newSlice)
					dstElement = newSlice
				} else {
					if isEmptyValue(dstElement) {
						dst.SetMapIndex(key, newSlice)
						dstElement = newSlice
					}

					if (!isEmptyValue(src) || overwriteWithEmptySrc || overwriteSliceWithEmptySrc) && (overwrite || isEmptyValue(dst)) && !config.AppendSlice && !sliceDeepCopy {
						if typeCheck && srcSlice.Type() != dstSlice.Type() {
							return fmt.Errorf("cannot override two slices with different type (%s, %s)", srcSlice.Type(), dstSlice.Type())
						}
						dstSlice = srcSlice
					} else if config.AppendSlice {
						if srcSlice.Type() != dstSlice.Type() {
							return fmt.Errorf("cannot append two slices with different type (%s, %s)", srcSlice.Type(), dstSlice.Type())
						}
						dstSlice = reflect.AppendSlice(dstSlice, srcSlice)
					} else if sliceDeepCopy {
						i := 0
						for ; i < srcSlice.Len() && i < dstSlice.Len(); i++ {
							srcElement := srcSlice.Index(i)
							dstElement := dstSlice.Index(i)

							if srcElement.CanInterface() {
								srcElement = reflect.ValueOf(srcElement.Interface())
							}
							if dstElement.CanInterface() {
								dstElement = reflect.ValueOf(dstElement.Interface())
							}

							if err = deepMerge(dstElement, srcElement, visited, depth+1, config); err != nil {
								return
							}
						}

					}
				}
			}

			dstElement, err = deepMerge(dstElement, srcElement, visited, depth+1, config)
			if err != nil {
				return
			}

			dst.SetMapIndex(key, dstElement)
		}
	case reflect.Slice:
		if !dst.CanSet() {
			break
		}
		if (!isEmptyValue(src) || overwriteWithEmptySrc || overwriteSliceWithEmptySrc) && (overwrite || isEmptyValue(dst)) && !config.AppendSlice && !sliceDeepCopy {
			dst.Set(src)
		} else if config.AppendSlice {
			if typeCheck && src.Type() != dst.Type() {
				err = fmt.Errorf("cannot append two slice with different type (%s, %s)", src.Type(), dst.Type())
				return
			}
			newSlice = reflect.AppendSlice(dst, src)
		} else if config.deepMergeSlice {
			if typeCheck && src.Type() != dst.Type() {
				err = fmt.Errorf("cannot deep merge two slice with different type (%s, %s)", src.Type(), dst.Type())
				return
			}
			dst.Set(reflect.AppendSlice(dst, src))
		} else if sliceDeepCopy {
			for i := 0; i < src.Len() && i < dst.Len(); i++ {
				srcElement := src.Index(i)
				dstElement := dst.Index(i)
				if srcElement.CanInterface() {
					srcElement = reflect.ValueOf(srcElement.Interface())
				}
				if dstElement.CanInterface() {
					dstElement = reflect.ValueOf(dstElement.Interface())
				}

				if err = deepMerge(dstElement, srcElement, visited, depth+1, config); err != nil {
					return
				}
			}
		}
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		if src.IsNil() {
			break
		}

			if src.Len() > dst.Len() {
				newSlice = reflect.MakeSlice(dst.Type(), src.Len(), src.Len())

				for i := 0; i < dst.Len(); i++ {
					newSlice.Index(i).Set(dst.Index(i))
				}
			}

			for i := 0; i < src.Len(); i++ {
				srcElem := src.Index(i)
				dstElem := newSlice.Index(i)

				if srcElem.CanInterface() {
					srcElem = reflect.ValueOf(srcElem.Interface())
				}

				if dstElem.CanInterface() {
					dstElem = reflect.ValueOf(dstElem.Interface())
				}

				if newSlice.Index(i).IsZero() {
					newSlice.Index(i).Set(srcElem)
					continue
				}

				dstElem, err = deepMerge(dstElem, srcElem, visited, depth+1, config)

				if err != nil {
					return
				}

				newSlice.Index(i).Set(dstElem)
			}
		}

		if dst.CanSet() {
			dst.Set(newSlice)
		} else {
			dst = newSlice
		}
	case reflect.Ptr, reflect.Interface:
		if isReflectNil(src) {
			break
		}

		if src.Kind() != reflect.Interface {
			if dst.IsNil() || (src.Kind() != reflect.Ptr && overwrite) {
				if dst.CanSet() && (overwrite || isEmptyValue(dst)) {
					dst.Set(src)
				}
			} else if src.Kind() == reflect.Ptr {
				if dst, err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
					return
				}
				dst = dst.Addr()
			} else if dst.Elem().Type() == src.Type() {
				if dst, err = deepMerge(dst.Elem(), src, visited, depth+1, config); err != nil {
					return
				}
			} else {
				err = ErrDifferentArgumentsTypes
				return
			}

			break
		}

		if dst.IsNil() || overwrite {
			if (overwrite || isEmptyValue(dst)) && (overwriteWithEmptySrc || !isEmptyValue(src)) {
				if dst.CanSet() {
					dst.Set(src)
				} else {
					dst = src
				}
			}
			break
		}

		if dst.Elem().Kind() == src.Elem().Kind() {
			if dst, err = deepMerge(dst.Elem(), src.Elem(), visited, depth+1, config); err != nil {
				return
			}
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
			}

			break
		}
	default:
		mustSet := (!isEmptyValue(src) || overwriteWithEmptySrc) && (overwrite || isEmptyValue(dst))
		if mustSet {
			if dst.CanSet() {
				dst.Set(src)
			} else {
				dst = src
			}
		}
	}

	return dst, nil
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

// WithOverwriteWithEmptyValue will make merge override non empty dst attributes with empty src attributes values.
func WithOverwriteWithEmptyValue(config *Config) {
	config.Overwrite = true
	config.overwriteWithEmptyValue = true
}

// WithOverrideEmptySlice will make merge override empty dst slice with empty src slice.
func WithOverrideEmptySlice(config *Config) {
	config.overwriteSliceWithEmptyValue = true
}

// WithAppendSlice will make merge append slices instead of overwriting it.
func WithAppendSlice(config *Config) {
	config.TypeCheck = true
	config.AppendSlice = true
}

// WithTypeCheck will make merge check types while overwriting it (must be used with WithOverride).
func WithTypeCheck(config *Config) {
	config.TypeCheck = true
}

// WithSliceDeepCopy will merge slice element one by one with Overwrite flag.
func WithSliceDeepCopy(config *Config) {
	config.sliceDeepCopy = true
	config.Overwrite = true
}

func merge(dst, src interface{}, opts ...func(*Config)) error {
	if dst != nil && reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return ErrNonPointerAgument
	}
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
