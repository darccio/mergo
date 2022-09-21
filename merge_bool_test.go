package mergo_test

import (
	"testing"

	"github.com/narvar/mergo"
)

type testMergeBoolStruct struct {
	BoolField bool
}

func TestMergeBool(t *testing.T) {
	dst := testMergeBoolStruct{
		BoolField: false,
	}
	src := testMergeBoolStruct{
		BoolField: true,
	}

	if err := mergo.Merge(&dst, src); err != nil {
		t.Error(err)
	}

	if dst.BoolField != src.BoolField {
		t.Error("dst is supposed to be the same as src")
	}
}

func TestSkipMergeBool(t *testing.T) {
	dst := testMergeBoolStruct{
		BoolField: false,
	}
	src := testMergeBoolStruct{
		BoolField: true,
	}

	if err := mergo.Merge(&dst, src, mergo.WithSkipMergingBool); err != nil {
		t.Error(err)
	}

	if dst.BoolField == src.BoolField {
		t.Error("dst is supposed to be different than src")
	}
}
