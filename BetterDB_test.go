package betterdb

import (
	"testing"
)

func TestNamedQuery(t *testing.T) {
	var args = map[string]interface{}{"name": "jim"}
	NamedQuery(nil, "select * from user where name=:name and password=:password and ct='12:34:5' and name=:name", nil, args)
}
