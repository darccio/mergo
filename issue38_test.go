package mergo_test

import (
	"testing"
	"time"

	"dario.cat/mergo"
)

type structWithoutTimePointer struct {
	Created time.Time
}

func TestIssue38Merge(t *testing.T) {
	dst := structWithoutTimePointer{
		time.Now(),
	}

	expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	src := structWithoutTimePointer{
		expected,
	}

	if err := mergo.Merge(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if dst.Created.Equal(src.Created) {
		t.Errorf("Created merged unexpectedly: dst.Created(%v) == src.Created(%v)", dst.Created, src.Created)
	}
}

func TestIssue38MergeEmptyStruct(t *testing.T) {
	dst := structWithoutTimePointer{}

	expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	src := structWithoutTimePointer{
		expected,
	}

	if err := mergo.Merge(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if dst.Created.Equal(src.Created) {
		t.Errorf("Created merged unexpectedly: dst.Created(%v) == src.Created(%v)", dst.Created, src.Created)
	}
}

func TestIssue38MergeWithOverwrite(t *testing.T) {
	dst := structWithoutTimePointer{
		time.Now(),
	}

	expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	src := structWithoutTimePointer{
		expected,
	}

	if err := mergo.MergeWithOverwrite(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if !dst.Created.Equal(src.Created) {
		t.Errorf("Created not merged in properly: dst.Created(%v) != src.Created(%v)", dst.Created, src.Created)
	}
}
