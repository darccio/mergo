package mergo_test

import (
	"testing"

	"dario.cat/mergo"
)

type structWithBlankField struct {
	_ struct{}
	A struct{}
}

func TestIssue174(t *testing.T) {
	dst := structWithBlankField{}
	src := structWithBlankField{}

	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		t.Error(err)
	}
}
