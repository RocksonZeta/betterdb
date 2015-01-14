package betterdb

import (
	"reflect"
	"regexp"
	"strconv"
	"time"
)

type SqlEscapeNotSupportError struct {
	Type string
}

func (this *SqlEscapeNotSupportError) Error() string {
	return "Can not escape " + this.Type + " type."
}

/**
SqlString.escape = function(val, stringifyObjects, timeZone) {
  if (val === undefined || val === null) {
    return 'NULL';
  }

  switch (typeof val) {
    case 'boolean': return (val) ? 'true' : 'false';
    case 'number': return val+'';
  }

  if (val instanceof Date) {
    val = SqlString.dateToString(val, timeZone || 'local');
  }

  if (Buffer.isBuffer(val)) {
    return SqlString.bufferToString(val);
  }

  if (Array.isArray(val)) {
    return SqlString.arrayToList(val, timeZone);
  }

  if (typeof val === 'object') {
    if (stringifyObjects) {
      val = val.toString();
    } else {
      return SqlString.objectToValues(val, timeZone);
    }
  }

  val = val.replace(/[\0\n\r\b\t\\\'\"\x1a]/g, function(s) {
    switch(s) {
      case "\0": return "\\0";
      case "\n": return "\\n";
      case "\r": return "\\r";
      case "\b": return "\\b";
      case "\t": return "\\t";
      case "\x1a": return "\\Z";
      default: return "\\"+s;
    }
  });
  return "'"+val+"'";
};

*/
func Escape(obj interface{}, objMapFns ...map[string]func(obj interface{}) string) (s string, e error) {
	re := regexp.MustCompile("[\000\n\r\b\t\\'\"\x1a]")
	if nil == obj {
		s = "NULL"
		return
	}
	switch v := obj.(type) {
	case string:
		s = "'" + re.ReplaceAllStringFunc(v, func(c string) string {
			switch c {
			case "\000":
				return "\\000"
			case "\n":
				return "\\n"
			case "\r":
				return "\\r"
			case "\b":
				return "\\b"
			case "\t":
				return "\\t"
			case "\x1a":
				return "\\Z"
			default:
				return "\\" + c
			}
		}) + "'"
		return
	case int8:
		s = strconv.FormatInt(int64(v), 10)
		return
	case int16:
		s = strconv.FormatInt(int64(v), 10)
		return
	case int32:
		s = strconv.FormatInt(int64(v), 10)
		return
	case int:
		s = strconv.FormatInt(int64(v), 10)
		return
	case int64:
		s = strconv.FormatInt(v, 10)
		return
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', -1, 32)
		return
	case float64:
		s = strconv.FormatFloat(v, 'f', -1, 64)
		return
	case bool:
		if v {
			s = "true"
		} else {
			s = "false"
		}
		return
	}
	v := reflect.ValueOf(obj)
	switch v.Kind().String() {
	case "slice":
		var ss string
		for i := 0; i < v.Len(); i++ {
			ss, e = Escape(v.Index(i).Interface())
			if nil != e {
				return
			}
			s += ss
			if i < v.Len()-1 {
				s += ","
			}
		}
		return
	case "ptr":
		s, e = Escape(v.Elem().Interface())
		return
	}
	t := reflect.TypeOf(obj)
	p := t.PkgPath() + "." + t.Name()
	var mapFuncs map[string]func(obj interface{}) string
	if 0 < len(objMapFns) {
		mapFuncs = objMapFns[0]
	}
	if nil != mapFuncs {
		if fn, ok := mapFuncs[p]; ok {
			s = fn(obj)
			return
		}
	}
	switch p {
	case "time.Time":
		timeObj := obj.(time.Time)
		s = timeObj.Format("2006-01-02 15:04:05")
		return
	}

	e = &SqlEscapeNotSupportError{Type: v.Kind().String()}
	return
}

func EscapeSql(s string, args ...string) (rs string, e error) {
	re := regexp.MustCompile("?")
	kvs := KeyValues{namedArgs}
	i := 0
	st = re.ReplaceAllStringFunc(s, func(key string) string {

		es, e = Escape()
	})
	return
}
