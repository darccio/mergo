package mergo

import (
	"testing"
)

func TestIssue121WithSliceDeepCopy(t *testing.T) {
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

	if err := Merge(&dst, src, WithSliceDeepCopy); err != nil {
		t.Errorf("Error during the merge: %v", err)
	}

	if dst["inter"].(map[string]interface{})["a"].(string) != "3" {
		t.Error("inter.a should equal '3'")
	}

	if dst["inter"].(map[string]interface{})["c"].(string) != "4" {
		t.Error("inter.c should equal '4'")
	}
}
