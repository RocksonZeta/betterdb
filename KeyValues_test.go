package betterdb

import "testing"

func TestGetMap(t *testing.T) {
	m := map[string]interface{}{"Key": "value"}
	kvs := KeyValues{m}
	v := kvs.Get("Key")
	if "value" != v {
		t.Errorf("TestGetMap failed")
	}
}
func TestGetMapRef(t *testing.T) {
	m := &map[string]interface{}{"Key": "value"}
	kvs := KeyValues{m}
	v := kvs.Get("Key")
	if "value" != v {
		t.Errorf("TestGetMapRef failed")
	}
}

func TestGetMapIntRef(t *testing.T) {
	m := &map[string]int{"Key": 1}
	kvs := KeyValues{m}
	v := kvs.Get("Key")
	if 1 != v {
		t.Errorf("TestGetMapIntRef failed")
	}
}

func TestGetMapRefRef(t *testing.T) {
	m := &map[string]interface{}{"Key": "value"}
	kvs := KeyValues{&m}
	v := kvs.Get("Key")
	if "value" != v {
		t.Errorf("TestGetMapRef failed")
	}
}

type testStruct struct {
	S1 string
}

func TestGetStruct(t *testing.T) {
	m := &testStruct{"s1"}
	kvs := KeyValues{&m}
	v := kvs.Get("S1")
	if "s1" != v {
		t.Errorf("TestGetMapRef failed")
	}
}

func TestGetNil(t *testing.T) {
	var m interface{} = nil
	kvs := KeyValues{m}

	defer func() {
		e := recover()
		if _, ok := e.(*KeyValuesError); !ok {
			t.Errorf("TestGetNil failed")
		}
	}()
	kvs.Get("Key")
}
