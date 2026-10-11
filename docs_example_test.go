package mergo_test

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"dario.cat/mergo"
)

func ExampleMerge_defaults() {
	type Config struct {
		Host string
		Port int
	}
	src := Config{Host: "default.example", Port: 8080}
	dst := Config{Host: "custom.example"}
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	// Output:
	// {custom.example 8080}
	// {default.example 8080}
}

func ExampleMerge_emptyValues() {
	type Config struct {
		Enabled bool
		Retries int
	}
	src := Config{}
	dst := Config{Enabled: true, Retries: 3}
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	if err := mergo.Merge(&dst, src, mergo.WithOverwriteWithEmptyValue); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	// Output:
	// {true 3}
	// {false 0}
}

func ExampleMerge_mapOverwrite() {
	dst := map[string]int{"retries": 3, "keep": 1}
	src := map[string]int{"retries": 0}
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst["retries"], dst["keep"], len(dst))
	if err := mergo.Merge(&dst, src, mergo.WithOverwriteWithEmptyValue); err != nil {
		panic(err)
	}
	_, kept := dst["keep"]
	fmt.Println(dst["retries"], kept, len(dst))
	// Output:
	// 0 1 2
	// 0 false 1
}

func ExampleMerge_emptySlice() {
	type Config struct{ Names []string }
	dst := Config{Names: []string{"local"}}
	src := Config{Names: []string{}}
	if err := mergo.Merge(&dst, src, mergo.WithOverrideEmptySlice); err != nil {
		panic(err)
	}
	fmt.Println(dst.Names)
	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithOverrideEmptySlice); err != nil {
		panic(err)
	}
	fmt.Println(len(dst.Names), dst.Names == nil)
	src.Names = nil
	if err := mergo.Merge(&dst, src, mergo.WithOverrideEmptySlice); err != nil {
		panic(err)
	}
	fmt.Println(dst.Names == nil)
	// Output:
	// [local]
	// 0 false
	// true
}

func ExampleMerge_appendSlice() {
	dst := []string{"local", "shared"}
	src := []string{"shared", "default"}
	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithAppendSlice); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	// Output: [local shared shared default]
}

func ExampleMerge_sliceDeepCopy() {
	type Item struct{ Name string }
	first := &Item{Name: "local"}
	dst := []*Item{first, {Name: "keep"}}
	src := []*Item{{Name: "source"}, {Name: ""}, {Name: "ignored"}}
	if err := mergo.Merge(&dst, src, mergo.WithSliceDeepCopy); err != nil {
		panic(err)
	}
	fmt.Println(len(dst), dst[0].Name, dst[1].Name)
	fmt.Println(dst[0] == first, dst[0] == src[0])
	values := []int{1}
	if err := mergo.Merge(&values, []int{2, 3}, mergo.WithSliceDeepCopy); err != nil {
		panic(err)
	}
	fmt.Println(values)
	// Output:
	// 2 source keep
	// true false
	// [1]
}

func ExampleMerge_typeCheck() {
	dst := map[string]interface{}{"items": []int{1}}
	src := map[string]interface{}{"items": []string{"one"}}
	err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithTypeCheck)
	fmt.Println(err != nil)
	fmt.Printf("%T\n", dst["items"])
	// Output:
	// true
	// []int
}

func ExampleMerge_pointerDefaults() {
	type Config struct{ Enabled *bool }
	local, fallback := false, true
	dst, src := Config{Enabled: &local}, Config{Enabled: &fallback}
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(local, dst.Enabled == &local)
	local = false
	if err := mergo.Merge(&dst, src, mergo.WithoutDereference); err != nil {
		panic(err)
	}
	fmt.Println(local, dst.Enabled == &local)
	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithoutDereference); err != nil {
		panic(err)
	}
	fmt.Println(*dst.Enabled, dst.Enabled == src.Enabled, local)
	// Output:
	// true true
	// false true
	// true true false
}

func ExampleMerge_pointerStruct() {
	type Item struct{ Name string }
	type Config struct{ Item *Item }
	dst := Config{Item: &Item{Name: "local"}}
	src := Config{Item: &Item{Name: "source"}}
	if err := mergo.Merge(&dst, src, mergo.WithOverride, mergo.WithoutDereference); err != nil {
		panic(err)
	}
	fmt.Println(dst.Item.Name, dst.Item == src.Item)
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst.Item.Name, dst.Item == src.Item)
	// Output:
	// local false
	// source false
}

func ExampleMerge_nestedMaps() {
	dst := map[string]map[string]int{"service": {"port": 8080}}
	src := map[string]map[string]int{"service": {"port": 9000, "retries": 3}}
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(dst["service"]["port"], dst["service"]["retries"])
	// Output: 8080 3
}

func ExampleMerge_mapStructValues() {
	type Item struct {
		Name string
		Port int
	}
	dst := map[string]Item{"service": {Name: "local"}}
	src := map[string]Item{"service": {Name: "source", Port: 8080}}
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(dst["service"])
	if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
		panic(err)
	}
	fmt.Println(dst["service"])
	// Output:
	// {local 0}
	// {source 8080}
}

func ExampleMerge_aliasing() {
	type Config struct {
		Names []string
		Item  *string
	}
	name := "source"
	src := Config{Names: []string{"source"}, Item: &name}
	var dst Config
	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}
	dst.Names[0] = "changed"
	*dst.Item = "changed"
	fmt.Println(src.Names[0], name, dst.Item == src.Item)
	// Output: changed changed true
}

func ExampleMap() {
	type Config struct {
		Host string
		Port int
	}
	dst := Config{Host: "local.example"}
	src := map[string]interface{}{"host": "default.example", "port": 8080}
	if err := mergo.Map(&dst, src); err != nil {
		panic(err)
	}
	fmt.Println(dst)
	mapped := map[string]interface{}{}
	if err := mergo.Map(&mapped, dst); err != nil {
		panic(err)
	}
	fmt.Println(mapped["host"], mapped["port"])
	err := mergo.Map(&dst, map[string]interface{}{"port": float64(8080)})
	fmt.Println(err != nil)
	// Output:
	// {local.example 8080}
	// local.example 8080
	// true
}

type docsTimeTransformer struct {
	Overwrite bool
}

func (t docsTimeTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ != reflect.TypeOf(time.Time{}) {
		return nil
	}
	return func(dst, src reflect.Value) error {
		if !dst.CanSet() {
			return nil
		}
		source := src.Interface().(time.Time)
		destination := dst.Interface().(time.Time)
		if !source.IsZero() && (destination.IsZero() || t.Overwrite) {
			dst.Set(src)
		}
		return nil
	}
}

func ExampleMerge_timeTransformer() {
	type Snapshot struct{ Time time.Time }
	first := time.Date(2025, 5, 7, 12, 0, 0, 0, time.UTC)
	later := first.Add(time.Hour)
	var dst Snapshot
	if err := mergo.Merge(&dst, Snapshot{Time: first}, mergo.WithTransformers(docsTimeTransformer{})); err != nil {
		panic(err)
	}
	fmt.Println(dst.Time.Format(time.RFC3339))
	if err := mergo.Merge(&dst, Snapshot{Time: later}, mergo.WithOverride, mergo.WithTransformers(docsTimeTransformer{})); err != nil {
		panic(err)
	}
	fmt.Println(dst.Time.Equal(first))
	if err := mergo.Merge(&dst, Snapshot{Time: later}, mergo.WithTransformers(docsTimeTransformer{Overwrite: true})); err != nil {
		panic(err)
	}
	fmt.Println(dst.Time.Equal(later))
	if err := mergo.Merge(&dst, Snapshot{}, mergo.WithTransformers(docsTimeTransformer{Overwrite: true})); err != nil {
		panic(err)
	}
	fmt.Println(dst.Time.Equal(later))
	// Output:
	// 2025-05-07T12:00:00Z
	// true
	// true
	// true
}

type docsPort int

var errDocsPort = errors.New("port must be positive")

type docsPortTransformer struct{}

func (docsPortTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ != reflect.TypeOf(docsPort(0)) {
		return nil
	}
	return func(dst, src reflect.Value) error {
		if src.Int() <= 0 {
			return errDocsPort
		}
		if dst.CanSet() && dst.Int() == 0 {
			dst.Set(src)
		}
		return nil
	}
}

func ExampleMerge_transformerError() {
	type Config struct {
		Name string
		Port docsPort
	}
	var dst Config
	src := Config{Name: "service", Port: -1}
	err := mergo.Merge(&dst, src, mergo.WithTransformers(docsPortTransformer{}))
	fmt.Println(errors.Is(err, errDocsPort), dst.Name, dst.Port)
	// Output: true service 0
}
