package mergo_test

import (
	"errors"
	"reflect"
	"testing"

	"dario.cat/mergo"
)

func TestMergeMapSliceTransformer(t *testing.T) {
	appendStrings := &transformer{m: map[reflect.Type]func(reflect.Value, reflect.Value) error{
		reflect.TypeOf([]string{}): func(dst, src reflect.Value) error {
			dst.Set(reflect.AppendSlice(dst, src))
			return nil
		},
	}}

	t.Run("typed map", func(t *testing.T) {
		dst := map[string][]string{"values": {"root"}}
		src := map[string][]string{"values": {"import"}}
		if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTransformers(appendStrings)); err != nil {
			t.Fatal(err)
		}
		if want := []string{"root", "import"}; !reflect.DeepEqual(dst["values"], want) {
			t.Fatalf("got %v, want %v", dst["values"], want)
		}
	})

	t.Run("nested untyped map", func(t *testing.T) {
		dst := map[string]interface{}{"plugin": map[string]interface{}{"values": []string{"root"}, "other": "keep"}}
		src := map[string]interface{}{"plugin": map[string]interface{}{"values": []string{"import"}}}
		if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTransformers(appendStrings)); err != nil {
			t.Fatal(err)
		}
		want := map[string]interface{}{"plugin": map[string]interface{}{"values": []string{"root", "import"}, "other": "keep"}}
		if !reflect.DeepEqual(dst, want) {
			t.Fatalf("got %#v, want %#v", dst, want)
		}
	})

	t.Run("array of tables and duplicate entries", func(t *testing.T) {
		amd64 := map[string]interface{}{"platform": "linux/amd64", "snapshotter": "overlayfs"}
		arm64 := map[string]interface{}{"platform": "linux/arm64", "snapshotter": "native"}
		dst := map[string]interface{}{"plugin": map[string]interface{}{"unpack_config": []interface{}{amd64}}}
		src := map[string]interface{}{"plugin": map[string]interface{}{"unpack_config": []interface{}{amd64, arm64, arm64}}}
		tf := &transformer{m: map[reflect.Type]func(reflect.Value, reflect.Value) error{
			reflect.TypeOf([]interface{}{}): func(dst, src reflect.Value) error {
				for i := 0; i < src.Len(); i++ {
					found := false
					for j := 0; j < dst.Len(); j++ {
						if reflect.DeepEqual(dst.Index(j).Interface(), src.Index(i).Interface()) {
							found = true
							break
						}
					}
					if !found {
						dst.Set(reflect.Append(dst, src.Index(i)))
					}
				}
				return nil
			},
		}}
		if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTransformers(tf)); err != nil {
			t.Fatal(err)
		}
		got := dst["plugin"].(map[string]interface{})["unpack_config"]
		if want := []interface{}{amd64, arm64}; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("without matching transformer", func(t *testing.T) {
		dst := map[string][]string{"values": {"root"}}
		src := map[string][]string{"values": {"import"}}
		if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTransformers(&transformer{})); err != nil {
			t.Fatal(err)
		}
		if want := []string{"import"}; !reflect.DeepEqual(dst["values"], want) {
			t.Fatalf("got %v, want %v", dst["values"], want)
		}
	})

	t.Run("empty source", func(t *testing.T) {
		dst := map[string][]string{"values": {"root"}}
		src := map[string][]string{"values": {}}
		if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTransformers(appendStrings)); err != nil {
			t.Fatal(err)
		}
		if want := []string{"root"}; !reflect.DeepEqual(dst["values"], want) {
			t.Fatalf("got %v, want %v", dst["values"], want)
		}
	})

	t.Run("missing destination retains source", func(t *testing.T) {
		dst := map[string][]string{}
		src := map[string][]string{"values": {"import"}}
		if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTransformers(appendStrings)); err != nil {
			t.Fatal(err)
		}
		if want := []string{"import"}; !reflect.DeepEqual(dst["values"], want) {
			t.Fatalf("got %v, want %v", dst["values"], want)
		}
	})

	t.Run("nil destination retains source", func(t *testing.T) {
		dst := map[string][]string{"values": nil}
		src := map[string][]string{"values": {"import"}}
		called := false
		tf := &transformer{m: map[reflect.Type]func(reflect.Value, reflect.Value) error{
			reflect.TypeOf([]string{}): func(dst, src reflect.Value) error { called = true; return nil },
		}}
		if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTransformers(tf)); err != nil {
			t.Fatal(err)
		}
		if called || !reflect.DeepEqual(dst["values"], []string{"import"}) {
			t.Fatalf("transformer called: %v, values: %v", called, dst["values"])
		}
	})

	t.Run("different slice types retain override behavior", func(t *testing.T) {
		dst := map[string]interface{}{"values": []string{"root"}}
		src := map[string]interface{}{"values": []int{1}}
		if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTransformers(appendStrings)); err != nil {
			t.Fatal(err)
		}
		if want := []int{1}; !reflect.DeepEqual(dst["values"], want) {
			t.Fatalf("got %v, want %v", dst["values"], want)
		}
	})
}

func TestMergeMapSliceTransformerError(t *testing.T) {
	wantErr := errors.New("transform failed")
	dst := map[string][]string{"values": {"root"}}
	src := map[string][]string{"values": {"import"}}
	tf := &transformer{m: map[reflect.Type]func(reflect.Value, reflect.Value) error{
		reflect.TypeOf([]string{}): func(dst, src reflect.Value) error { return wantErr },
	}}
	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTransformers(tf)); !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if want := []string{"root"}; !reflect.DeepEqual(dst["values"], want) {
		t.Fatalf("got %v, want %v", dst["values"], want)
	}
}
