package betterdb

import (
	"database/sql"
	"reflect"
	"regexp"
	"strings"
)

type ExecResult struct {
	InsertId     int64
	RowsAffected int64
	Error        error
}

func ReflectScan(args []reflect.Value) []reflect.Value {
	result := []reflect.Value{}
	rows := args[0].Interface().(*sql.Rows)
	var e error
	cols, e := rows.Columns()
	if nil != e {
		result = append(result, reflect.ValueOf(e))
		return result
	}
	record := args[1].Interface().(reflect.Value).Elem()
	fields := make([]interface{}, len(cols))
	for i, n := range cols {
		if f := record.FieldByName(n); !f.IsValid() {
			fields[i] = nil
		} else {
			fields[i] = record.FieldByName(n).Addr().Interface()
		}
	}
	e = rows.Scan(fields...)
	if nil != e {
		result = append(result, reflect.ValueOf(e))
		return result
	}
	return result
}
func MakeScan(scanFn *func(row *sql.Rows, dest reflect.Value) error) {
	fn := reflect.ValueOf(scanFn).Elem()
	v := reflect.MakeFunc(fn.Type(), ReflectScan)
	fn.Set(v)
}

func Scan(st *sql.Stmt, records interface{}, args ...interface{}) error {
	rows, e := st.Query(args...)
	if nil != e {
		return e
	}
	defer rows.Close()
	var scan func(row *sql.Rows, dest reflect.Value) error
	MakeScan(&scan)
	recordType := reflect.TypeOf(records).Elem().Elem()
	results := reflect.ValueOf(records).Elem()
	var scanError error
	for rows.Next() {
		record := reflect.New(recordType)
		scanError = scan(rows, record)
		if nil != scanError {
			return scanError
		}
		results.Set(reflect.Append(results, record.Elem()))
	}
	return nil
}

func NamedScan(db *sql.DB, s string, records interface{}, namedArgs interface{}) error {
	sqlstring, args := TransNameStr(s, namedArgs)
	st, e := db.Prepare(sqlstring)
	if nil != e {
		return e
	}
	return Scan(st, records, args...)
}

func Select(db *sql.DB, s string, records interface{}, args ...interface{}) error {
	st, e := db.Prepare(s)
	if e != nil {
		return e
	}
	defer st.Close()
	Scan(st, records, args...)
	return nil
}

//variable placeholder should have this form ":var",eg."select name from user where name=:name"
func NamedSelect(db *sql.DB, s string, records interface{}, namedArgs interface{}) error {
	st, args := TransNameStr(s, namedArgs)
	return Select(db, st, records, args...)
}

func Exec(st *sql.Stmt, args []interface{}) (result ExecResult) {
	r, e := st.Exec(args...)
	if nil != e {
		result.Error = e
		return
	}
	result.RowsAffected, e = r.RowsAffected()
	if nil != e {
		result.Error = e
		return
	}
	result.InsertId, e = r.LastInsertId()
	if nil != e {
		result.Error = e
		return
	}
	return
}

func NamedUpdate(db *sql.DB, s string, args interface{}) []ExecResult {
	//st, e := db.Prepare(s)
	//if nil != e {
	//	return []ExecResult{ExecResult{Error: e}}
	//}
	//rargs := reflect.ValueOf(args)RowsAffected

	return []ExecResult{}

}

/**
map (:Name,:Age ,{"Name":"jim" , "Age":20}) -> ("?,?",['jim',20])
*/
func TransNameStr(s string, namedArgs interface{}) (st string, args []interface{}) {
	re := regexp.MustCompile(":[^0-9]\\w*")
	kvs := KeyValues{namedArgs}
	st = re.ReplaceAllStringFunc(s, func(key string) string {
		if n, ok := kvs.Get(strings.TrimPrefix(key, ":")); ok {
			args = append(args, n)
			return "?"
		}
		return key
	})
	return
}

type BetterDB struct {
	*sql.DB
}

func (this *BetterDB) Select(s string, records interface{}, args ...interface{}) error {
	return Select(this.DB, s, records, args...)
}

func (this *BetterDB) NamedSelect(s string, records interface{}, args interface{}) error {
	return NamedSelect(this.DB, s, records, args)
}
func (this *BetterDB) Update(s string, args ...interface{}) ExecResult {
	st, e := this.DB.Prepare(s)
	if nil != e {
		return ExecResult{Error: e}
	}
	return Exec(st, args)
}
func (this *BetterDB) NamedUpdate(s string, records interface{}) error {
	return nil
}

//
func (this *BetterDB) UpdateBatch(s string, args interface{}) error {
	return nil
}

func (this *BetterDB) Post(table string, obj interface{}) error {
	return nil
}
func (this *BetterDB) Get(table string, id interface{}, obj interface{}) error {
	return nil
}
func (this *BetterDB) Put(table string, id interface{}, newValues interface{}) error {
	return nil
}
func (this *BetterDB) Delete(table string, id interface{}) error {
	return nil
}

func (this *BetterDB) BatchSqls(sqls []string) error {
	return nil
}

/**
eg.insert into user(Name,Age) values(:Name,:Age) [{Name:"jim" , Age:12}]
*/
func (this *BetterDB) Batch(s string, values []interface{}) error {
	return nil
}
