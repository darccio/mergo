package mergo_test

import (
	"testing"

	"github.com/narvar/mergo"
	"github.com/stretchr/testify/assert"
)

type testMergeSliceStruct struct {
	SliceField []string
}

func TestMergeSlice(t *testing.T) {
	for _, test := range []struct {
		Name                 string
		Src                  *testMergeSliceStruct
		Dst                  *testMergeSliceStruct
		ExpectedErr          error
		ExpectedResult       []string
		AppendSlice          bool
		AppendSliceReversely bool
	}{
		{
			Name:           "Merge src to dst with AppendSlice option",
			Src:            &testMergeSliceStruct{[]string{"a", "b"}},
			Dst:            &testMergeSliceStruct{[]string{"c", "d"}},
			ExpectedErr:    nil,
			ExpectedResult: []string{"c", "d", "a", "b"},
			AppendSlice:    true,
		},
		{
			Name:           "Merge src to dst with AppendSlice option and src is nil",
			Src:            &testMergeSliceStruct{nil},
			Dst:            &testMergeSliceStruct{[]string{"c", "d"}},
			ExpectedErr:    nil,
			ExpectedResult: []string{"c", "d"},
			AppendSlice:    true,
		},
		{
			Name:           "Merge src to dst with AppendSlice option and src is empty",
			Src:            &testMergeSliceStruct{[]string{}},
			Dst:            &testMergeSliceStruct{[]string{"c", "d"}},
			ExpectedErr:    nil,
			ExpectedResult: []string{"c", "d"},
			AppendSlice:    true,
		},
		{
			Name:           "Merge src to dst with AppendSlice option and dst is nil",
			Src:            &testMergeSliceStruct{[]string{"a", "b"}},
			Dst:            &testMergeSliceStruct{nil},
			ExpectedErr:    nil,
			ExpectedResult: []string{"a", "b"},
			AppendSlice:    true,
		},
		{
			Name:           "Merge src to dst with AppendSlice option and dst is empty",
			Src:            &testMergeSliceStruct{[]string{"a", "b"}},
			Dst:            &testMergeSliceStruct{[]string{}},
			ExpectedErr:    nil,
			ExpectedResult: []string{"a", "b"},
			AppendSlice:    true,
		},
		{
			Name:           "Merge src to dst with AppendSlice option, src and dst are nil",
			Src:            &testMergeSliceStruct{nil},
			Dst:            &testMergeSliceStruct{nil},
			ExpectedErr:    nil,
			ExpectedResult: nil,
			AppendSlice:    true,
		},
		{
			Name:           "Merge src to dst with AppendSlice option, src and dst are empty",
			Src:            &testMergeSliceStruct{[]string{}},
			Dst:            &testMergeSliceStruct{[]string{}},
			ExpectedErr:    nil,
			ExpectedResult: []string{},
			AppendSlice:    true,
		},
		{
			Name:           "Merge src to dst with AppendSlice option, src is empty and dst is nil",
			Src:            &testMergeSliceStruct{[]string{}},
			Dst:            &testMergeSliceStruct{nil},
			ExpectedErr:    nil,
			ExpectedResult: nil,
			AppendSlice:    true,
		},
		{
			Name:           "Merge src to dst with AppendSlice option, src is nil and dst is empty",
			Src:            &testMergeSliceStruct{nil},
			Dst:            &testMergeSliceStruct{[]string{}},
			ExpectedErr:    nil,
			ExpectedResult: []string{},
			AppendSlice:    true,
		},
		{
			Name:                 "Merge src to dst with AppendSliceReversely option",
			Src:                  &testMergeSliceStruct{[]string{"a", "b"}},
			Dst:                  &testMergeSliceStruct{[]string{"c", "d"}},
			ExpectedErr:          nil,
			ExpectedResult:       []string{"a", "b", "c", "d"},
			AppendSliceReversely: true,
		},
		{
			Name:                 "Merge src to dst with AppendSliceReversely option and src is nil",
			Src:                  &testMergeSliceStruct{nil},
			Dst:                  &testMergeSliceStruct{[]string{"c", "d"}},
			ExpectedErr:          nil,
			ExpectedResult:       []string{"c", "d"},
			AppendSliceReversely: true,
		},
		{
			Name:                 "Merge src to dst with AppendSliceReversely option and src is empty",
			Src:                  &testMergeSliceStruct{[]string{}},
			Dst:                  &testMergeSliceStruct{[]string{"c", "d"}},
			ExpectedErr:          nil,
			ExpectedResult:       []string{"c", "d"},
			AppendSliceReversely: true,
		},
		{
			Name:                 "Merge src to dst with AppendSliceReversely option and dst is nil",
			Src:                  &testMergeSliceStruct{[]string{"a", "b"}},
			Dst:                  &testMergeSliceStruct{nil},
			ExpectedErr:          nil,
			ExpectedResult:       []string{"a", "b"},
			AppendSliceReversely: true,
		},
		{
			Name:                 "Merge src to dst with AppendSliceReversely option and dst is empty",
			Src:                  &testMergeSliceStruct{[]string{"a", "b"}},
			Dst:                  &testMergeSliceStruct{[]string{}},
			ExpectedErr:          nil,
			ExpectedResult:       []string{"a", "b"},
			AppendSliceReversely: true,
		},
		{
			Name:                 "Merge src to dst with AppendSliceReversely option, src and dst are nil",
			Src:                  &testMergeSliceStruct{nil},
			Dst:                  &testMergeSliceStruct{nil},
			ExpectedErr:          nil,
			ExpectedResult:       nil,
			AppendSliceReversely: true,
		},
		{
			Name:                 "Merge src to dst with AppendSliceReversely option, src and dst are empty",
			Src:                  &testMergeSliceStruct{[]string{}},
			Dst:                  &testMergeSliceStruct{[]string{}},
			ExpectedErr:          nil,
			ExpectedResult:       []string{},
			AppendSliceReversely: true,
		},
		{
			Name:                 "Merge src to dst with AppendSliceReversely option, src is empty and dst is nil",
			Src:                  &testMergeSliceStruct{[]string{}},
			Dst:                  &testMergeSliceStruct{nil},
			ExpectedErr:          nil,
			ExpectedResult:       []string{},
			AppendSliceReversely: true,
		},
		{
			Name:                 "Merge src to dst with AppendSliceReversely option, src is nil and dst is empty",
			Src:                  &testMergeSliceStruct{nil},
			Dst:                  &testMergeSliceStruct{[]string{}},
			ExpectedErr:          nil,
			ExpectedResult:       nil,
			AppendSliceReversely: true,
		},
	} {
		var err error
		if test.AppendSlice {
			err = mergo.Merge(test.Dst, test.Src, mergo.WithAppendSlice)
		} else if test.AppendSliceReversely {
			err = mergo.Merge(test.Dst, test.Src, mergo.WithAppendSliceReversely)
		} else {
			err = mergo.Merge(test.Dst, test.Src)
		}

		if test.ExpectedErr != nil {
			assert.NotEmpty(t, err)
			assert.Equal(t, test.ExpectedErr.Error(), err.Error())
		} else {
			assert.Equal(t, test.ExpectedResult, test.Dst.SliceField)
		}
	}
}
