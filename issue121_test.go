package mergo

import (
	"testing"
)

func TestIssue121WithSliceDeepMerge(t *testing.T) {
	dst := map[string]interface{}{
		"a": "1",
		"b": []map[string]interface{}{
			{"c": "2"},
		},
	}
	src := map[string]interface{}{
		"b": []map[string]interface{}{
			{"c": "3", "d": "1"},
			{"e": "1", "f": "1", "g": []string{"1", "2"}},
		},
	}
	if err := Merge(&dst, src, WithDeepMergeSlice); err != nil {
		t.Fatalf("Error during the merge: %v", err)
	}
	if dst["a"].(string) != "1" {
		t.Error("a should equal '1'")
	}
	if dst["b"].([]map[string]interface{})[0]["c"] != "2" {
		t.Error("b.[0].c should equal '2'")
	}
	if dst["b"].([]map[string]interface{})[0]["d"] != "1" {
		t.Error("b.[0].d should equal '2'")
	}
	if dst["b"].([]map[string]interface{})[1]["e"] != "1" {
		t.Error("b.[1].e should equal '1'")
	}
	if dst["b"].([]map[string]interface{})[1]["g"].([]string)[0] != "1" {
		t.Error("b.[1].g[0] should equal '1'")
	}
}

func TestIssue121WithSliceDeepMergeFromPR126(t *testing.T) {
	dst := map[string]interface{}{
		"inter": map[string]interface{}{
			"a": "1",
			"b": "2",
		},
	}
	src := map[string]interface{}{
		"inter": map[string]interface{}{
			"a": "3",
			"c": "4",
		},
	}
	if err := Merge(&dst, src, WithDeepMergeSlice); err != nil {
		t.Fatalf("Error during the merge: %v", err)
	}
	if dst["inter"].(map[string]interface{})["a"].(string) != "1" {
		t.Error("inter.a should equal '1'")
	}
	if dst["inter"].(map[string]interface{})["b"].(string) != "2" {
		t.Error("inter.a should equal '2'")
	}
	if dst["inter"].(map[string]interface{})["c"].(string) != "4" {
		t.Error("inter.c should equal '4'")
	}
}

type order struct {
	A       string
	B       int64
	Details []detail
}

type detail struct {
	A string
	B string
}

func TestIssue121(t *testing.T) {
	src := order{
		A: "one",
		B: 2,
		Details: []detail{
			{
				B: "B",
			},
		},
	}
	dst := order{
		A: "two",
		Details: []detail{
			{
				A: "one",
			},
		},
	}
	if err := Merge(&dst, src, WithDeepMergeSlice); err != nil {
		t.Fatalf("Error during the merge: %v", err)
	}
	if len(dst.Details) != 1 {
		t.Fatalf("B was not properly merged: %+v", dst.Details)
	}
	if dst.Details[0].B != "B" {
		t.Fatalf("B was not properly merged: %+v", dst.Details[0].B)
	}
}
