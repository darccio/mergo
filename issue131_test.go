package mergo

import (
	"testing"
)

type foz struct {
	A *bool
	B string
}

func TestIssue131MergeWithOverwriteWithEmptyValue(t *testing.T) {
	src := foz{
		A: func(v bool) *bool { return &v }(false),
		B: "src",
	}
	dest := foz{
		A: func(v bool) *bool { return &v }(true),
		B: "dest",
	}
	Merge(&dest, src, WithOverwriteWithEmptyValue)
	if *src.A != *dest.A {
		t.Fatalf("dest.A not merged in properly: %v != %v", *src.A, *dest.A)
	}
	if src.B != dest.B {
		t.Fatalf("dest.B not merged in properly: %v != %v", src.B, dest.B)
	}
}
