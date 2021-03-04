package mergo

import (
	"encoding/json"
	"fmt"
	"testing"
)

func pp(i interface{}) string {
	b, _ := json.MarshalIndent(i, "", "    ")
	return string(b)
}

func TestIssue121WithSliceDeepMerge(t *testing.T) {
	dst := map[string]interface{}{
		"a": "1",
		"b": []map[string]interface{}{
			map[string]interface{}{"c": "2"},
		},
	}

	src := map[string]interface{}{
		"b": []map[string]interface{}{
			map[string]interface{}{"c": "3", "d": "1"},
			map[string]interface{}{"e": "1", "f": "1", "g": []string{"1", "2"}},
		},
	}

	if err := Merge(&dst, src, WithSliceDeepMerge); err != nil {
		t.Fatalf("Error during the merge: %v", err)
	}

	fmt.Println(pp(dst))

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
