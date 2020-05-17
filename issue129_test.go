package mergo

import (
	"testing"
)

func TestIssue129Boolean(t *testing.T) {
	type Foo struct {
		A bool
		B bool
	}

	src := Foo{
		A: true,
		B: false,
	}
	dst := Foo{
		A: false,
		B: true,
	}

	// Standard behavior
	if err := Merge(&dst, src); err != nil {
		t.Fatal(err)
	}
	if dst.A == false {
		t.Fatalf("expected true, got false")
	}
	if dst.B == false {
		t.Fatalf("expected true, got false")
	}

	// Expected behavior
	dst = Foo{
		A: false,
		B: true,
	}
	if err := Merge(&dst, src, WithOverwriteWithEmptyValue); err != nil {
		t.Fatal(err)
	}
	if dst.A == false {
		t.Fatalf("expected true, got false")
	}
	if dst.B == true {
		t.Fatalf("expected false, got true")
	}
}
