package betterdb

import (
	"testing"
	"time"
)

func TestEscape(t *testing.T) {
	s, e := Escape("hello")
	if e != nil {
		t.Error(e)
	}
	if s != "'hello'" {
		t.Errorf("Escape(\"hello\") failed,expected 'hello' ,got " + s)
	}
}

func TestEscapeParticular(t *testing.T) {
	s, e := Escape("hell'o")
	if e != nil {
		t.Error(e)
	}
	if s != "'hell\\'o'" {
		t.Errorf("TestEscapeParticular failed,expected 'hell\\'o' ,got " + s)
	}
}

func TestEscapeArray(t *testing.T) {
	s := "tom"
	a := []interface{}{1, "jim", nil}
	s, e := Escape(a)
	if e != nil {
		t.Error(e)
		return
	}
	if s != "1,'jim',NULL" {
		t.Errorf("TestEscapeArray failed,expected 1,'jim',NULL ,got " + s)
	}
}

func TestEscapeTime(t *testing.T) {
	now, _ := time.Parse("2006-01-02 15:04:05", "2012-12-13 14:15:16")
	s, e := Escape(now)
	if e != nil {
		t.Error(e)
		return
	}
	if s != "2012-12-13 14:15:16" {
		t.Errorf("TestEscapeTime failed,expected 2012-12-13 14:15:16 ,got " + s)
	}
}
