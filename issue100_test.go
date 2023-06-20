package mergo_test

import (
	"testing"

	"dario.cat/mergo"
)

type issue100s struct {
	Member interface{}
}

func TestIssue100(t *testing.T) {
	m := make(map[string]interface{})
	m["Member"] = "anything"

	st := &issue100s{}
	if err := mergo.Map(st, m); err != nil {
		t.Error(err)
	}
}
