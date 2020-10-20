package mergo_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/imdario/mergo"
)

// This test has been stripped to minimally reproduce the core issue of #167. Consider the playground link for an in-depth analysis.

type Issue167A struct{}

type Issue167InMap struct {
	As []*Issue167A
}

type Issue167ValueMap struct {
	InMap map[string]Issue167InMap
}

type Issue167Transformer struct{}

func (t Issue167Transformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	switch typ {
	case reflect.TypeOf([]*Issue167A{}):
		return t.transform
	}

	return nil
}

func (t Issue167Transformer) transform(dst, src reflect.Value) error {
	if !dst.CanSet() {
		return fmt.Errorf("dst.CanSet() is false, so we're unable to Set() the new value later")
	}

	dstReal := dst.Interface().([]*Issue167A)
	dst.Set(reflect.ValueOf(dstReal))

	return nil
}

var Issue167TestData = []struct {
	Dst Issue167ValueMap
	Src Issue167ValueMap
}{
	{
		// Each As must have at least an element
		Dst: Issue167ValueMap{InMap: map[string]Issue167InMap{"foo": {As: []*Issue167A{{}}}}},
		Src: Issue167ValueMap{InMap: map[string]Issue167InMap{"foo": {As: []*Issue167A{{}}}}},
	},
}

func TestIssue167MapWithValue(t *testing.T) {
	// And additionally should not panic with every test case
	for _, data := range Issue167TestData {
		err := mergo.Merge(&data.Dst, data.Src, mergo.WithTransformers(Issue167Transformer{}))
		if err != nil {
			t.Errorf("Error while merging %s", err)
		}
	}
}
