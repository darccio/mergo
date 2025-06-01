package mergo_test

import (
	"testing"

	"dario.cat/mergo"
)

func BenchmarkMerge(b *testing.B) {
	type ts struct {
		Field string
	}

	var (
		dst = ts{}
		src = ts{
			Field: "src",
		}
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := mergo.Merge(&dst, src)
		if err != nil {
			b.Fatal(err)
		}
		dst.Field = ""
	}
}
