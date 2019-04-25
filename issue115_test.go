package mergo

import (
	"testing"
)

func TestIssue115MergeMapWithNilValueToMapWithOverride(t *testing.T) {
	p1 := map[string]interface{}{
		"A": 0, "B": 1, "C": 2,
	}
	p2 := map[string]interface{}{
		"D": nil,
	}
	if err := Map(&p1, p2, WithOverride); err != nil {
		t.Fatalf("Error during the merge: %v", err)
	}
	if _, ok := p1["D"]; !ok  {
		t.Errorf("p1 should contain D: %+v", p1)
	}
}
