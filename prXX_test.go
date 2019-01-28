package mergo

import (
	"testing"
)

type flagStruct struct {
	Save 	bool  `mergo:"savedst"`
}

type invalidFlagStruct struct {
	Save 	bool 
}


func Test_MergoWithoutFlags(t *testing.T) {
	src := invalidFlagStruct{false}
	dst := invalidFlagStruct{true}
	if err := Merge(&dst, src); err != nil {
		t.FailNow()
	}

	if dst.Save == false {
		t.Fatalf("dst.Save was saved which is wrong")
	}
}


func Test_MergoWithFlags(t *testing.T) {
	src := flagStruct{false}
	dst := flagStruct{true}
	if err := Merge(&dst, src); err != nil {
		t.FailNow()
	}

	if dst.Save == false {
		t.Fatalf("dst.Save was not saved but overriden")
	}
}