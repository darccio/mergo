package mergo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPtr struct {
	Str    string
	Number int
	Nested *Nested
}

type Test struct {
	Str    string
	Number int
	Nested Nested
}

type Nested struct {
	NestedStr    string
	NestedNumber int
}

func TestPtr2NestedStruct(t *testing.T) {
	dstNested := Nested{
		NestedStr:    "This is a non empty string in a nested struct",
		NestedNumber: 42,
	}
	dst := TestPtr{
		Str:    "This is a non empty test string",
		Nested: &dstNested,
	}
	src := TestPtr{
		Str:    "This is also a string",
		Number: 4,
		Nested: &Nested{
			NestedStr:    "Bad string",
			NestedNumber: 7,
		},
	}
	// the nested struct is not empty in dst and a merge without override
	// should not override the dst, with the value from src
	if err := Merge(&dst, src); err != nil {
		t.FailNow()
	}
	assert.Equalf(t, &dstNested, dst.Nested, "dst.Nested has changed without using option 'WithOverride'")

	// the nested struct in dst is still not empty, however merging with "WithOverride"
	// should overwrite the dst.Nested with src.Nested
	if err := Merge(&dst, src, WithOverride); err != nil {
		t.FailNow()
	}
	assert.Equalf(t, src.Nested, dst.Nested, "dst.Nested and src.Nested are not equal but should, because we merged with option 'WithOverride'")

	// we set src.Nested to a zero value ptr and reset dst.Nested to its original value
	// this should never be merged to dst also when using option "WithOverride"
	src.Nested = &Nested{}
	dst.Nested = &dstNested
	if err := Merge(&dst, src); err != nil {
		t.FailNow()
	}
	assert.Equalf(t, &dstNested, dst.Nested, "src.Nested was merged even though it had a zero value ptr")
	if err := Merge(&dst, src, WithOverride); err != nil {
		t.FailNow()
	}
	assert.Equal(t, &dstNested, dst.Nested, "src.Nested was merged using option 'WithOverride' even though it had a zero value ptr")
}

func TestNonPtrToNestedStruct(t *testing.T) {
	dstNested := Nested{
		NestedStr:    "This is a non empty string in a nested struct",
		NestedNumber: 42,
	}

	dst := Test{
		Str:    "This is a non empty test string",
		Nested: dstNested,
	}

	src := Test{
		Str:    "This is also a string",
		Number: 4,
		Nested: Nested{
			NestedStr:    "New string",
			NestedNumber: 7,
		},
	}
	// the nested struct is not empty in dst and a merge without override
	// should not override the dst, with the value from src
	if err := Merge(&dst, src); err != nil {
		t.FailNow()
	}
	assert.Equalf(t, dst.Nested, dstNested, "dst.Nested has changed without using option 'WithOverride'")

	// the nested struct in dst is still not empty, however merging with "WithOverride"
	// should overwrite the dst.Nested with src.Nested
	if err := Merge(&dst, src, WithOverride); err != nil {
		t.FailNow()
	}
	assert.Equalf(t, dst, src, "dst.Nested and src.Nested are not equal but should, because we merged with option 'WithOverride'")

	// we set src.Nested to the zero value of its type and reset dst.Nested to its original value
	// this should never be merged to dst also when using option "WithOverride"
	src.Nested = Nested{}
	dst.Nested = dstNested
	if err := Merge(&dst, src); err != nil {
		t.FailNow()
	}
	assert.Equalf(t, dstNested, dst.Nested, "src.Nested was merged even though it had a zero value")
	if err := Merge(&dst, src, WithOverride); err != nil {
		t.FailNow()
	}
	assert.Equal(t, dstNested, dst.Nested, "src.Nested was merged using option 'WithOverride' even though it had a zero value")
}
