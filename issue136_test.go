package mergo

import (
	"testing"
)

type embeddedTestA struct {
	Name string
	Age  uint8
}

type embeddedTestB struct {
	embeddedTestA
	Address string
}

func TestMergeEmbedded(t *testing.T) {
	a := &embeddedTestA{
		"Suwon", 16,
	}

	b := &embeddedTestB{}

	err := Merge(&b.embeddedTestA, *a)

	if b.Name != "Suwon" {
		t.Errorf("%v %v", b.Name, err)
	}
}
