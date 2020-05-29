package mergo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type ifaceTypesTest struct {
	N       int
	Handler http.Handler
}

type ifaceTypesHandler int

func (*ifaceTypesHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Test", "ifaceTypesHandler")
}

func TestMergeInterfaceWithDifferentConcreteTypes(t *testing.T) {
	dst := ifaceTypesTest{
		Handler: new(ifaceTypesHandler),
	}

	src := ifaceTypesTest{
		N: 42,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.Header().Set("Test", "handlerFunc")
		}),
	}

	if err := Merge(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	rw := httptest.NewRecorder()
	dst.Handler.ServeHTTP(rw, nil)

	if got, want := rw.Header().Get("Test"), "ifaceTypesHandler"; got != want {
		t.Errorf("Handler not merged in properly: got %q header value %q, want %q", "Test", got, want)
	}
}

func TestMergeInterfaceWithSameConcreteTypes(t *testing.T) {
	type testStruct struct {
		Name string
		Value string
	}
	type interfaceStruct struct {
		Field interface{}
	}
	dst := interfaceStruct{
		Field: testStruct{
			Value: "keepMe",
		},
	}

	src := interfaceStruct{
		Field: testStruct{
			Name: t.Name(),
		},
	}

	if err := Merge(&dst, src); err != nil {
		t.Errorf("Error while merging %s", err)
	}

	dstData := dst.Field.(testStruct)
	srcData := src.Field.(testStruct)
	if dstData.Name != srcData.Name {
		t.Errorf("dst name was not updated: got %s, want %s", dstData.Name, srcData.Name)
	}
	if dstData.Value != "keepMe" {
		t.Errorf("dst value was not preserved: got %s, want %s", dstData.Value, "keepMe")
	}
}
