package ng_test

import (
	"errors"
	"testing"

	mergo "dario.cat/mergo/internal/ng"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	type ts struct {
		Field int
	}

	t.Run("merge", func(t *testing.T) {
		t.Parallel()

		var (
			dst = ts{
				Field: 0,
			}
			src = ts{
				Field: 1,
			}
		)

		if err := mergo.Merge(&dst, src); err != nil {
			t.Errorf("Error while merging %s", err)
		}

		if dst.Field != src.Field {
			t.Errorf("Expected dst.Field to be %d, got %d", src.Field, dst.Field)
		}
	})

	t.Run("merge from pointer", func(t *testing.T) {
		t.Parallel()

		var (
			dst = &ts{
				Field: 0,
			}
			src = &ts{
				Field: 1,
			}
		)

		if err := mergo.Merge(dst, src); err != nil {
			t.Errorf("Error while merging %s", err)
		}

		if dst.Field != src.Field {
			t.Errorf("Expected dst.Field to be %d, got %d", src.Field, dst.Field)
		}
	})
}

func TestMergePredefinedType(t *testing.T) {
	t.Parallel()

	t.Run("merge ints", func(t *testing.T) {
		t.Parallel()

		var (
			dst = 0
			src = 1
		)

		if err := mergo.Merge(&dst, src); err != nil {
			t.Errorf("Error while merging %s", err)
		}

		if dst != src {
			t.Errorf("Expected dst to be %d, got %d", src, dst)
		}
	})

	t.Run("merge strings", func(t *testing.T) {
		t.Parallel()

		var (
			dst = ""
			src = "src"
		)

		if err := mergo.Merge(&dst, src); err != nil {
			t.Errorf("Error while merging %s", err)
		}

		if dst != src {
			t.Errorf("Expected dst to be %q, got %q", src, dst)
		}
	})
}

func TestMergeNil(t *testing.T) {
	t.Parallel()

	t.Run("both nil", func(t *testing.T) {
		t.Parallel()

		var naerr *mergo.NilArgumentsError

		if err := mergo.Merge(nil, nil); !errors.As(err, &naerr) {
			t.Error(err)
		}
	})

	t.Run("dst nil", func(t *testing.T) {
		t.Parallel()

		var naerr *mergo.NilArgumentsError

		if err := mergo.Merge(nil, struct{}{}); !errors.As(err, &naerr) {
			t.Error(err)
		}
	})

	t.Run("src nil", func(t *testing.T) {
		t.Parallel()

		var (
			dst   = struct{}{}
			naerr *mergo.NilArgumentsError
		)

		if err := mergo.Merge(&dst, nil); !errors.As(err, &naerr) {
			t.Error(err)
		}
	})
}

func TestMergeNotZero(t *testing.T) {
	t.Parallel()

	type ts struct {
		Field int
	}

	var (
		dst = ts{
			Field: 1,
		}
		src = ts{
			Field: 2,
		}
	)

	if err := mergo.Merge(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if dst.Field == src.Field {
		t.Errorf("Expected dst.Field to be 1, got %d", dst.Field)
	}
}

func TestMergeOnlyExportedFields(t *testing.T) {
	t.Parallel()

	type ts struct {
		Field      int
		unexported int
	}

	var (
		dst = ts{
			Field:      0,
			unexported: 2,
		}
		src = ts{
			Field:      3,
			unexported: 4,
		}
	)

	if err := mergo.Merge(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	if dst.Field != src.Field {
		t.Errorf("Expected dst.Field to be %d, got %d", src.Field, dst.Field)
	}

	if dst.unexported == src.unexported {
		t.Errorf("Expected dst.unexported to be 2, got %d", dst.unexported)
	}
}

func TestMergeNotPointer(t *testing.T) {
	t.Parallel()

	var npeerr *mergo.InvalidDestinationError
	if err := mergo.Merge(struct{}{}, struct{}{}); !errors.As(err, &npeerr) {
		t.Error(err)
	}
}

func TestMergeDifferentTypes(t *testing.T) {
	t.Parallel()

	dst := struct {
		Field int
	}{
		Field: 1,
	}

	src := struct {
		Field string
	}{
		Field: "src",
	}

	var dteerr *mergo.DifferentArgumentTypesError
	if err := mergo.Merge(&dst, src); !errors.As(err, &dteerr) {
		t.Error(err)
	}
}

func TestMergeValue(t *testing.T) {
	t.Parallel()

	type ts struct {
		Field int
	}

	dst := ts{
		Field: 0,
	}

	src := ts{
		Field: 1,
	}

	if err := mergo.MergeValue(&dst, src); err != nil {
		t.Error(err)
	}

	if dst.Field != src.Field {
		t.Errorf("Expected dst.Field to be %d, got %d", src.Field, dst.Field)
	}
}

func TestMergeValueNil(t *testing.T) {
	t.Parallel()

	var naerr *mergo.NilArgumentsError

	if err := mergo.MergeValue(nil, struct{}{}); !errors.As(err, &naerr) {
		t.Error(err)
	}
}

func TestMergePtr(t *testing.T) {
	t.Parallel()

	type ts struct {
		Field int
	}

	var (
		dst = &ts{
			Field: 0,
		}
		src = &ts{
			Field: 1,
		}
	)

	if err := mergo.MergePtr(dst, src); err != nil {
		t.Error(err)
	}

	if dst.Field != src.Field {
		t.Errorf("Expected dst.Field to be %d, got %d", src.Field, dst.Field)
	}
}

func TestMergePtrNil(t *testing.T) {
	t.Parallel()

	type ts struct {
		Field int
	}

	t.Run("dst nil", func(t *testing.T) {
		t.Parallel()

		var (
			src = &ts{
				Field: 0,
			}
			naerr *mergo.NilArgumentsError
		)

		if err := mergo.MergePtr(nil, src); !errors.As(err, &naerr) {
			t.Error(err)
		}
	})

	t.Run("src nil", func(t *testing.T) {
		t.Parallel()

		var (
			dst = &ts{
				Field: 0,
			}
			naerr *mergo.NilArgumentsError
		)

		if err := mergo.MergePtr(dst, nil); !errors.As(err, &naerr) {
			t.Error(err)
		}
	})

	t.Run("both nil", func(t *testing.T) {
		t.Parallel()

		var (
			dst   *ts
			src   *ts
			naerr *mergo.NilArgumentsError
		)

		if err := mergo.MergePtr(dst, src); !errors.As(err, &naerr) {
			t.Error(err)
		}
	})
}

func TestErrorMessages(t *testing.T) {
	t.Parallel()

	t.Run("NilArgumentsError", func(t *testing.T) {
		t.Parallel()

		msg := "src and dst must not be nil"
		if err := mergo.Merge(nil, nil); err.Error() != msg {
			t.Errorf("Expected error to be %q, got %q", msg, err.Error())
		}
	})

	t.Run("InvalidDestinationError", func(t *testing.T) {
		t.Parallel()

		msg := "dst must be a pointer"
		if err := mergo.Merge(struct{}{}, struct{}{}); err.Error() != msg {
			t.Errorf("Expected error to be %q, got %q", msg, err.Error())
		}
	})

	t.Run("DifferentArgumentTypesError", func(t *testing.T) {
		t.Parallel()

		dst := struct {
			Field int
		}{
			Field: 1,
		}

		src := struct {
			Field string
		}{
			Field: "src",
		}

		msg := "dst and src must have the same type"
		if err := mergo.Merge(&dst, src); err.Error() != msg {
			t.Errorf("Expected error to be %q, got %q", msg, err.Error())
		}
	})
}

func BenchmarkMerge(b *testing.B) {
	b.ReportAllocs()

	type ts struct {
		Field int
	}

	src := ts{
		Field: 2,
	}

	b.Run("Merge", func(b *testing.B) {
		dst := ts{
			Field: 0,
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if err := mergo.Merge(&dst, src); err != nil {
				b.Fatal(err)
			}

			dst.Field = 0
		}
	})

	b.Run("MergeValue", func(b *testing.B) {
		dst := ts{
			Field: 0,
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if err := mergo.MergeValue(&dst, src); err != nil {
				b.Fatal(err)
			}

			dst.Field = 0
		}
	})

	b.Run("MergePtr", func(b *testing.B) {
		dst := &ts{
			Field: 0,
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if err := mergo.MergePtr(dst, &src); err != nil {
				b.Fatal(err)
			}

			dst.Field = 0
		}
	})
}
