package betterdb

import "testing"

func TestTransNameStr(t *testing.T) {
	var args = map[string]interface{}{"name": "jim", "password": "123"}
	st, params := TransNameStr("name=:name and password=:password and ct='12:34:5' and name=:name", args)
	if st != "name=? and password=? and ct='12:34:5' and name=?" {
		t.Errorf("NamedQuery() faield , %s", st)
	}
	if 3 != len(params) {
		//t.Errorf("NamedQuery() faield , param len error")
	}
}
