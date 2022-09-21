package mergo_test

import (
	"testing"
	"time"

	"github.com/narvar/mergo"
)

type document struct {
	Created *time.Time
}

func TestIssue23MergeWithOverwrite(t *testing.T) {
	now := time.Now()
	dst := document{
		&now,
	}
	expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	src := document{
		&expected,
	}

	if err := mergo.MergeWithOverwrite(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if !dst.Created.Equal(*src.Created) { //--> https://golang.org/pkg/time/#pkg-overview
		t.Errorf("Created not merged in properly: dst.Created(%v) != src.Created(%v)", dst.Created, src.Created)
	}
}
