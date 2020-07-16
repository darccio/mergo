package mergo

import (
	"reflect"
	"testing"
)

type CustomStruct struct {
	SomeMap map[string]string
}

type issue90TestData struct {
	name string
	src  map[string]CustomStruct
	dst  map[string]CustomStruct
	exp  map[string]CustomStruct
}

func issue90Data() []issue90TestData {
	return []issue90TestData{
		{
			name: "Normal",
			dst: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{
						"key1": "loosethis", "key2": "keepthis",
					},
				},
			},
			src: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{
						"key1": "key10",
					},
				},
			},
			exp: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{
						"key1": "key10", "key2": "keepthis",
					},
				},
			},
		},
		{
			name: "Init of struct key",
			dst: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{},
				},
			},
			src: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{
						"key1": "key10",
					},
				},
			},
			exp: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{
						"key1": "key10",
					},
				},
			},
		},
		{
			name: "Not Init of struct key",
			dst:  map[string]CustomStruct{},
			src: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{
						"key1": "key10",
					},
				},
			},
			exp: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{
						"key1": "key10",
					},
				},
			},
		},
		{
			name: "Nil struct key",
			dst: map[string]CustomStruct{
				"a": {
					SomeMap: nil,
				},
			},
			src: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{
						"key1": "key10",
					},
				},
			},
			exp: map[string]CustomStruct{
				"a": {
					SomeMap: map[string]string{
						"key1": "key10",
					},
				},
			},
		},
	}
}

func TestMergoStructMap(t *testing.T) {
	for _, data := range issue90Data() {
		dst := data.dst
		src := data.src
		exp := data.exp

		err := Merge(&dst, src, WithAppendSlice, WithOverride)
		if err != nil {
			t.Errorf("mergo error was not nil, %v", err)
		}

		if !reflect.DeepEqual(dst, exp) {
			t.Errorf("Actual: %#v did not match \nExpected: %#v", dst, exp)
		}
	}
}
