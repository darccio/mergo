package mergo_test

import (
	"testing"

	"dario.cat/mergo/v2"
)

func litPtr[T any](v T) *T {
	return &v
}

func TestIntMerge(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		dst  *int
		src  int
		want *int
	}{
		{dst: litPtr(0), src: 1, want: litPtr(1)},
		{dst: litPtr(2), src: 1, want: litPtr(2)},
		{dst: nil, src: 1, want: nil},
		{dst: litPtr(3), src: 0, want: litPtr(3)},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run("", func(t *testing.T) {
			t.Parallel()

			mergo.Merge(tc.dst, tc.src)
			if tc.dst == nil {
				if tc.want != nil {
					t.Errorf("expected %v, got %v", *tc.want, tc.dst)
				}

				return
			}
			if *tc.dst != *tc.want {
				t.Errorf("expected %v, got %v", *tc.want, *tc.dst)
			}
		})
	}
}
