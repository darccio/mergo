package mergo

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestMergoSimpleMap(t *testing.T) {
	dst := map[string]string{"key1": "loosethis", "key2": "keepthis"}
	src := map[string]string{"key1": "key10"}
	exp := map[string]string{"key1": "key10", "key2": "keepthis"}
	err := Merge(&dst, src, WithAppendSlice, WithOverride)
	assert.Equal(t, err, nil)
	assert.Equal(t, dst, exp)
}

type CustomStruct struct {
	SomeMap map[string]string
}

func TestMergoStructMap(t *testing.T) {
	var testData = []struct {
		name string
		src  map[string]CustomStruct
		dst  map[string]CustomStruct
		exp  map[string]CustomStruct
	}{
		{name: "Normal",
			dst: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{"key1": "loosethis", "key2": "keepthis"}}},
			src: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{"key1": "key10"}}},
			exp: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{"key1": "key10", "key2": "keepthis"}}},
		},
		{name: "Init of struct key", dst: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{}}},
			src: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{"key1": "key10"}}},
			exp: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{"key1": "key10"}}},
		},
		{name: "Not Init of struct key", dst: map[string]CustomStruct{},
			src: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{"key1": "key10"}}},
			exp: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{"key1": "key10"}}},
		},
		{name: "Nil struct key", dst: map[string]CustomStruct{"a": CustomStruct{SomeMap: nil}},
			src: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{"key1": "key10"}}},
			exp: map[string]CustomStruct{"a": CustomStruct{SomeMap: map[string]string{"key1": "key10"}}}},
	}

	for _, data := range testData {
		dst := data.dst
		src := data.src
		exp := data.exp

		err := Merge(&dst, src, WithAppendSlice, WithOverride)
		assert.Equal(t, err, nil)
		assert.Equal(t, dst, exp)
	}
}
