package mergo

import (
	"testing"
	"reflect"
)

type simpleTest struct {
	Value int
}

type complexTest struct {
	St simpleTest
	sz int
	Id string
}

func TestNil(t *testing.T) {
	if err := Merge(nil, nil); err != NilArgumentsErr {
		t.Fail()
	}
}

func TestDifferentTypes(t *testing.T) {
	a := simpleTest{42}
	b := 42
	if err := Merge(&a, b); err != DifferentArgumentsTypesErr {
		t.Fail()
	}
}

func TestSimpleStruct(t *testing.T) {
	a := simpleTest{}
	b := simpleTest{42}
	if err := Merge(&a, b); err != nil {
		t.FailNow()
	}
	if a.Value != 42 {
		t.Fatalf("b not merged in a properly: a.Value(%d) != b.Value(%d)", a.Value, b.Value)
	}
	if !reflect.DeepEqual(a, b) {
		t.FailNow()
	}
}

func TestComplexStruct(t *testing.T) {
	a := complexTest{}
	a.Id = "athing"
	b := complexTest{simpleTest{42}, 1, "bthing"}
	if err := Merge(&a, b); err != nil {
		t.FailNow()
	}
	if a.St.Value != 42 {
		t.Fatalf("b not merged in a properly: a.St.Value(%d) != b.St.Value(%d)", a.St.Value, b.St.Value)
	}
	if a.sz == 1 {
		t.Fatalf("a's private field sz not preserved from merge: a.sz(%d) == b.sz(%d)", a.sz, b.sz)
	}
	if a.Id == b.Id {
		t.Fatalf("a's field Id not preserved from merge: a.Id(%s) == b.Id(%s)", a.Id, b.Id)
	}
}
