package mergo_test

import (
	"github.com/imdario/mergo"
	"reflect"
	"testing"
)

const (
	canaryString = "canary"
	canaryInt64  = 0xBEEF
)

type issue210s struct {
	CanaryString string
	MapVal       map[string]int64
	SliceVal     []int64
	Int64Ptr     *int64
	BoolPtr      *bool
	StringPtr    *string
	Canary       int64
}

func TestIssue210(t *testing.T) {
	var (
		boolVal      bool  = true
		int64Val     int64 = 100
		emptyInt64   int64 = 0
		emptyString        = ""
		canaryString       = "canary"
	)

	cases := []struct {
		dst    *issue210s
		src    *issue210s
		expect *issue210s
		policy []func(*mergo.Config)
	}{
		{
			dst: &issue210s{
				CanaryString: canaryString,
				MapVal: map[string]int64{
					canaryString: canaryInt64,
				},
				SliceVal:  []int64{100, 200},
				Int64Ptr:  &int64Val,
				BoolPtr:   &boolVal,
				StringPtr: &canaryString,
				Canary:    canaryInt64,
			},
			src: &issue210s{
				CanaryString: canaryString,
				StringPtr:    &emptyString,
				Int64Ptr:     &emptyInt64,
				Canary:       canaryInt64,
			},
			expect: &issue210s{
				CanaryString: canaryString,
				MapVal: map[string]int64{
					canaryString: canaryInt64,
				},
				SliceVal:  []int64{100, 200},
				Int64Ptr:  &emptyInt64,
				BoolPtr:   &boolVal,
				StringPtr: &emptyString,
				Canary:    canaryInt64,
			},
			policy: []func(*mergo.Config){
				mergo.WithSkipReflectNilSource,
				mergo.WithOverwriteWithEmptyValue,
			},
		},
	}

	for _, testCase := range cases {
		dst := testCase.dst
		src := testCase.src

		err := mergo.Merge(dst, src, testCase.policy...)

		if err != nil {
			t.Errorf("mergo TestIssue210 merge failed, %v", err)
		}

		if !reflect.DeepEqual(testCase.expect.MapVal, dst.MapVal) {
			t.Errorf("mergo TestIssue210 merge failed, map val not equal")
		}

		if !reflect.DeepEqual(testCase.expect.SliceVal, dst.SliceVal) {
			t.Errorf("mergo TestIssue210 merge failed, slice val not equal")
		}

		if *testCase.expect.BoolPtr != *dst.BoolPtr {
			t.Errorf("mergo TestIssue210 merge failed, bool val not equal")
		}

		if *testCase.expect.Int64Ptr != *dst.Int64Ptr {
			t.Errorf("mergo TestIssue210 merge failed, int64 val not equal, got %v", *dst.Int64Ptr)
		}

		if *testCase.expect.StringPtr != *dst.StringPtr {
			t.Errorf("mergo TestIssue210 merge failed, string val not equal, got %v", *dst.StringPtr)
		}
	}

}
