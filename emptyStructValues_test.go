package mergo_test

import (
	"testing"

	"github.com/imdario/mergo"
)

type parent struct {
	S     string
	Child *child
}

type child struct {
	S string
}

func TestDoNotOverwriteEmptyValueWithinStruct(t *testing.T) {
	src := parent{
		Child: &child{S: "src"},
		S:     "src",
	}
	dest := parent{
		Child: &child{S: ""},
		S:     "dest",
	}
	if err := mergo.Merge(&dest, src, mergo.WithNoOverrideEmptyStructValues); err != nil {
		t.Error(err)
	}
	if dest.Child.S != "" {
		t.Errorf("dest.Child.S overwritten")
	}
}

func TestOverwriteEmptyStruct(t *testing.T) {
	src := parent{
		Child: &child{S: "src"},
		S:     "src",
	}
	dest := parent{
		S: "dest",
	}
	if err := mergo.Merge(&dest, src, mergo.WithNoOverrideEmptyStructValues); err != nil {
		t.Error(err)
	}
	if dest.Child.S != "src" {
		t.Errorf("dest.Child.S not overwritten")
	}
}
