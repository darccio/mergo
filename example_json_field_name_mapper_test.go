// Copyright 2013 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mergo_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/imdario/mergo"
)

// NOTE: this is a very simplistic and quite poor reimplementation of the private
// parts of 'encoding/json' that handle the tags

var ErrMissingJSONTags = errors.New(`Missing JSON tags`)

type JSONTagsFieldNameMapper struct{}

func (m *JSONTagsFieldNameMapper) FromStructField(field reflect.StructField) (string, error) {
	tags := field.Tag.Get("json")

	if len(tags) == 0 {
		return "", ErrMissingJSONTags
	}

	split := strings.Split(tags, ",")

	if split[0] != "" {
		return split[0], nil
	}

	return "", ErrMissingJSONTags
}

func (m *JSONTagsFieldNameMapper) FromKeyName(dst reflect.Value, key string) (reflect.Value, string, error) {
	numField := dst.NumField()

	for i := 0; i < numField; i++ {
		fieldValue := dst.Field(i)
		field := dst.Type().Field(i)

		tags := field.Tag.Get("json")

		if len(tags) == 0 {
			continue
		}

		split := strings.Split(tags, ",")

		if split[0] == key {
			return fieldValue, field.Name, nil
		}
	}

	return reflect.Value{}, "", mergo.ErrNotSupported
}

func TestMapToStructWithJsonFieldNameMapper(t *testing.T) {
	var dst struct {
		Name string `json:"custom_name"`
		Age  int    `json:"custom_age"`
	}

	src := map[string]interface{}{
		"custom_name": "some name",
		"custom_age":  42,
	}

	if err := mergo.Map(&dst, &src, mergo.WithFieldNameMapper(&JSONTagsFieldNameMapper{})); err != nil {
		t.Errorf("expected nil, got %q", err)
	}

	if dst.Name != "some name" {
		t.Errorf("expected %s, got %s", "some name", dst.Name)
	}

	if dst.Age != 42 {
		t.Errorf("expected %d, got %d", 42, dst.Age)
	}
}

func TestMapFromStructWithJsonFieldNameMapper(t *testing.T) {
	src := struct {
		Name string `json:"custom_name"`
		Age  int    `json:"custom_age"`
	}{
		Name: "some name",
		Age:  42,
	}

	dst := map[string]interface{}{}

	if err := mergo.Map(&dst, src, mergo.WithFieldNameMapper(&JSONTagsFieldNameMapper{})); err != nil {
		t.Errorf("expected nil, got %q", err)
	}

	if dst["custom_name"].(string) != "some name" {
		t.Errorf("expected %s, got %v", "some name", dst["custom_name"])
	}

	if dst["custom_age"].(int) != 42 {
		t.Errorf("expected %d, got %v", 42, dst["custom_age"])
	}
}
