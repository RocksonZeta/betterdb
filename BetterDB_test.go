package betterdb

import (
	"fmt"
	"testing"
)

func TestTransNameStr(t *testing.T) {
	var args = map[string]interface{}{"name": "jim", "password": "123"}
	st, params := TransNameStr("name=:name and password=:password and ct='12:34:5' and name=:name", args)
	if st != "name=? and password=? and ct='12:34:5' and name=?" {
		t.Errorf("NamedQuery() faield")
	}
	if 3 != len(params) {
		t.Errorf("NamedQuery() faield , param len error")
	}
}

type TestPickStruct struct {
	S1 string
	I1 int
}

func TestPick(t *testing.T) {
	result := Pick(TestPickStruct{"S1", 1})
	fmt.Println(result)
}
