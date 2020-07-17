package mergo

import "testing"

type issue100s struct {
	Member interface{}
}

func TestIssue100(t *testing.T) {
	m := make(map[string]interface{})
	m["Member"] = "anything"

	st := &issue100s{}
	if err := Map(st, m); err != nil {
		t.Error(err)
	}
}
