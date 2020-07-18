package mergo_test

// Copyright 2020 Dario Castañé. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"testing"

	"github.com/imdario/mergo"
)

func TestMergeErrors(t *testing.T) {
	var (
		stub    = struct{}{}
		mapStub = make(map[string]interface{})
		must    = func(dst, src interface{}) func(*testing.T) {
			return func(t *testing.T) {
				if err := mergo.Merge(dst, src); err != nil {
					unexpected(t, nil, err)
				}
			}
		}
		mustFail = func(dst, src interface{}, expected error) func(*testing.T) {
			return func(t *testing.T) {
				if err := mergo.Merge(dst, src); err == nil {
					unexpected(t, expected, nil)
				}
			}
		}
	)
	testCases := []testCase{
		{
			desc: "nil dst",
			test: mustFail(nil, stub, mergo.ErrNilDestination),
		},
		{
			desc: "nil src",
			test: mustFail(&stub, nil, mergo.ErrNilSource),
		},
		{
			desc: "non-pointer dst",
			test: mustFail(stub, stub, mergo.ErrNonPointerDestination),
		},
		{
			desc: "equal types",
			test: must(&stub, stub),
		},
	}
	runTests(t, testCases)
}
