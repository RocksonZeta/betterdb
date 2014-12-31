package betterdb

import (
	"fmt"
	"testing"
)

func TestGetMap(t *testing.T) {
	m := map[string]interface{}{"Key": "value"}
	kvs := KeyValues{m}
	v, ok := kvs.Get("Key")
	if "value" != v || !ok {
		t.Errorf("TestGetMap failed")
	}
}
func TestGetMapRef(t *testing.T) {
	m := &map[string]interface{}{"Key": "value"}
	kvs := KeyValues{m}
	v, ok := kvs.Get("Key")
	if "value" != v || !ok {
		t.Errorf("TestGetMapRef failed")
	}
}

func TestGetMapIntRef(t *testing.T) {
	m := &map[string]int{"Key": 1}
	kvs := KeyValues{m}
	v, ok := kvs.Get("Key")
	if 1 != v || !ok {
		t.Errorf("TestGetMapIntRef failed")
	}
}

func TestGetMapRefRef(t *testing.T) {
	m := &map[string]interface{}{"Key": "value"}
	kvs := KeyValues{&m}
	v, _ := kvs.Get("Key")
	if "value" != v {
		t.Errorf("TestGetMapRef failed")
	}
}

type testStruct struct {
	S1 string
	s1 string
}

func TestGetStruct(t *testing.T) {
	m := &testStruct{"s1", "tom"}
	kvs := KeyValues{&m}
	v, _ := kvs.Get("S1")
	if "s1" != v {
		t.Errorf("TestGetMapRef failed")
	}
}

func TestGetNil(t *testing.T) {
	var m interface{} = nil
	kvs := KeyValues{m}

	kvs.Get("Key")
}

func TestMap(t *testing.T) {
	s := testStruct{"jim", "tom"}
	kvs := KeyValues{s}
	r := kvs.Map()
	v, _ := (r["S1"]).(string)
	if "jim" != v {
		t.Errorf("Map struct failed")
	}
}

func TestMapMap(t *testing.T) {
	s := &map[string]interface{}{"name": "jim", "age": 1}
	v := &s
	v1 := &v
	kvs := KeyValues{v1}
	r := kvs.Map()
	fmt.Println(r)
}

func TestPick(t *testing.T) {
	s := testStruct{"jim", "tom"}
	kv := KeyValues{s}
	m := kv.Pick("S1", "S2")
	fmt.Println("pick struct ", m)
}
