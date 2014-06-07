// Copyright 2013 Dario Castañé. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mergo

import (
    "gopkg.in/yaml.v1"
    "io/ioutil"
    "reflect"
    "testing"
)

type simpleTest struct {
    Value int
}

type complexTest struct {
    St  simpleTest
    sz  int
    Id  string
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

func TestMaps(t *testing.T) {
    m := map[string]simpleTest{
        "a": simpleTest{},
        "b": simpleTest{42},
    }
    n := map[string]simpleTest{
        "a": simpleTest{16},
        "b": simpleTest{},
        "c": simpleTest{12},
    }
    if err := Merge(&m, n); err != nil {
        t.Fatalf(err.Error())
    }
    if len(m) != 3 {
        t.Fatalf(`n not merged in m properly, m must have 3 elements instead of %d`, len(m))
    }
    if m["a"].Value != 0 {
        t.Fatalf(`n merged in m because I solved non-addressable map values TODO: m["a"].Value(%d) != n["a"].Value(%d)`, m["a"].Value, n["a"].Value)
    }
    if m["b"].Value != 42 {
        t.Fatalf(`n wrongly merged in m: m["b"].Value(%d) != n["b"].Value(%d)`, m["b"].Value, n["b"].Value)
    }
    if m["c"].Value != 12 {
        t.Fatalf(`n not merged in m: m["c"].Value(%d) != n["c"].Value(%d)`, m["c"].Value, n["c"].Value)
    }
}

func TestYAMLMaps(t *testing.T) {
    thing := loadYAML("testdata/thing.yml")
    license := loadYAML("testdata/license.yml")
    ft := thing["fields"].(map[interface{}]interface{})
    fl := license["fields"].(map[interface{}]interface{})
    expected_length := len(ft) + len(fl)
    if err := Merge(&license, thing); err != nil {
        t.Fatal(err.Error())
    }
    current_length := len(license["fields"].(map[interface{}]interface{}))
    if current_length != expected_length {
        t.Fatalf(`thing not merged in license properly, license must have %d elements instead of %d`, expected_length, current_length)
    }
    fields := license["fields"].(map[interface{}]interface{})
    if _, ok := fields["id"]; !ok {
        t.Fatalf(`thing not merged in license properly, license must have a new id field from thing`)
    }
}

func loadYAML(path string) (m map[string]interface{}) {
    m = make(map[string]interface{})
    raw, _ := ioutil.ReadFile(path)
    _ = yaml.Unmarshal(raw, &m)
    return
}
