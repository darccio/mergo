package mergo

import (
	"testing"
)

func TestIssue89Boolean(t *testing.T) {
	type Foo struct {
		Bar bool `json:"bar"`
	}

	src := Foo{Bar: true}
	dst := Foo{Bar: false}

	if err := Merge(&dst, src); err != nil {
		t.Fatal(err)
	}
	if dst.Bar == false {
		t.Fatalf("expected true, got false")
	}
}

func TestIssue89MergeWithEmptyValue(t *testing.T) {
	p1 := map[string]interface{}{
		"A": 3, "B": "note", "C": true,
	}
	p2 := map[string]interface{}{
		"B": "", "C": false,
	}
	if err := Merge(&p1, p2, WithOverwriteWithEmptyValue); err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		key      string
		expected interface{}
	}{
		{
			"A",
			3,
		},
		{
			"B",
			"",
		},
		{
			"C",
			false,
		},
	}
	for _, tC := range testCases {
		if p1[tC.key] != tC.expected {
			t.Fatalf("expected %v in p1[%q], got %v", tC.expected, tC.key, p1[tC.key])
		}
	}
}
