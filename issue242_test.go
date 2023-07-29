package mergo_test

import (
	"reflect"
	"testing"

	"dario.cat/mergo"
)

type StudentPinned struct {
	Name               string     `mergo:"pinned"`
	Birthplace         Birthplace `mergo:"pinned"`
	Books              []string
	IsMemberOfBookClub bool
}

type Birthplace struct {
	Country string
	City    string
}

func TestMergeWithPinnedFields(t *testing.T) {
	dst := StudentPinned{
		Name: "Artur",
		Birthplace: Birthplace{
			Country: "Ireland",
			City:    "Cork",
		},
		Books:              []string{"Clean Architecture", "Crime and Punishment"},
		IsMemberOfBookClub: false,
	}

	src := StudentPinned{
		Name: "Changed Name",
		Birthplace: Birthplace{
			Country: "Changed",
			City:    "Birthplace",
		},
		Books:              []string{"Mathematical analysis"},
		IsMemberOfBookClub: true,
	}

	expected := StudentPinned{
		Name: "Artur",
		Birthplace: Birthplace{
			Country: "Ireland",
			City:    "Cork",
		},
		Books:              []string{"Clean Architecture", "Crime and Punishment", "Mathematical analysis"},
		IsMemberOfBookClub: true,
	}

	err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithAppendSlice)
	if err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if !reflect.DeepEqual(dst, expected) {
		t.Errorf("expected: %#v\ngot: %#v", expected, dst)
	}
}
