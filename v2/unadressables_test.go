package mergo_test

import (
	"testing"

	"dario.cat/mergo/v2"
	"github.com/google/go-cmp/cmp"
)

func ifc[T any](v T) interface{} {
	return v
}

func TestInterfaceMerge(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		dst  *interface{}
		src  interface{}
		want interface{}
	}{
		{dst: litPtr(ifc(0)), src: 1, want: 1},
		{dst: litPtr(ifc(2)), src: 1, want: 2},
		{dst: nil, src: 1, want: nil},
		{dst: litPtr(ifc(3)), src: 0, want: 3},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run("", func(t *testing.T) {
			t.Parallel()

			mergo.Merge(tc.dst, tc.src)
			if tc.dst == nil {
				if tc.want != nil {
					t.Errorf("expected %v, got %v", tc.want, tc.dst)
				}

				return
			}
			if !cmp.Equal(*tc.dst, tc.want) {
				t.Errorf("expected %v, got %v", tc.want, *tc.dst)
			}
		})
	}
}
