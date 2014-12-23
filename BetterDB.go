package betterdb

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func ReflectScan(args []reflect.Value) []reflect.Value {
	rows := args[0].Interface().(*sql.Rows)
	cols, e := rows.Columns()
	if nil != e {
		panic(e)
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
	rows.Scan(fields...)
	return nil
}
func MakeScan(scanFn *func(row *sql.Rows, dest reflect.Value)) {
	fn := reflect.ValueOf(scanFn).Elem()
	v := reflect.MakeFunc(fn.Type(), ReflectScan)
	fn.Set(v)
}

func Query(db *sql.DB, s string, records interface{}, args ...interface{}) {
	st, e := db.Prepare(s)
	if e != nil {
		panic(e)
	}
	defer st.Close()
	ExecuteQuery(st, records, args)
}

//variable placeholder should have this form ":var",eg."select name from user where name=:name"
func NamedQuery(db *sql.DB, s string, records interface{}, namedArgs map[string]interface{}) {
	st, args := TransNameStr(s, namedArgs)
	fmt.Println(st, args)
	Query(db, st, records, args...)
}

func NamedUpdate(db *sql.DB, s string, args interface{}) {

}

func ExecuteQuery(st *sql.Stmt, records interface{}, args ...interface{}) {
	rows, e := st.Query(args...)
	if nil != e {
		panic(e)
	}
	defer rows.Close()
	var scan func(row *sql.Rows, dest reflect.Value)
	MakeScan(&scan)
	recordType := reflect.TypeOf(records).Elem().Elem()
	results := reflect.ValueOf(records).Elem()
	for rows.Next() {
		record := reflect.New(recordType)
		scan(rows, record)
		results.Set(reflect.Append(results, record.Elem()))
	}
}

func ExecuteUpdate(st *sql.Stmt, args ...interface{}) (insertId int64, affectRows int64) {
	r, e := st.Exec(args...)
	if nil != e {
		panic(e)
	}
	insertId, e := r.LastInsertId()
	if nil != e {
		panic(e)
	}
	affectRows, e := r.RowsAffected()
	if nil != e {
		panic(e)
	}
}

/**
map (:name,:age ,{"Name":"jim" , "Age":20}) -> ("?,?",['jim',20])
*/
func TransNameStr(s string, namedArgs map[string]interface{}) (st string, args []interface{}) {
	re := regexp.MustCompile(":\\w+")
	st = re.ReplaceAllStringFunc(s, func(key string) string {
		if n, ok := namedArgs[strings.TrimPrefix(key, ":")]; ok {
			args = append(args, n)
			return "?"
		}
		return key
	})
	return
}

func Pick(obj interface{}, keys ...string) (result map[string]interface{}) {
	values := reflect.ValueOf(obj)
	var field reflect.Value
	result = make(map[string]interface{})
	if 0 == len(keys) {
		for i := 0; i < values.NumField(); i++ {
			field = values.Field(i)
			result[field.String()] = field.Interface()
		}
	} else {
		for _, k := range keys {
			field = values.FieldByName(k)
			if field.IsValid() {
				result[field.String()] = field.Interface()
			}
		}
	}
	return
}

type BetterDB struct {
	*sql.DB
}

func (this *BetterDB) BetterQuery(s string, records interface{}, args interface{}) {

}

func (this *BetterDB) NamedQuery(s string, records interface{}, args interface{}) {

}

func (this *BetterDB) Post(table string, obj interface{}) {

}
func (this *BetterDB) Get(table string, id interface{}, obj interface{}) {

}
func (this *BetterDB) Put(table string, id interface{}, newValues interface{}) {

}
func (this *BetterDB) Delete(table string, id interface{}) {

}

func (this *BetterDB) BatchSqls(sqls []string) {

}

/**
eg.insert into user(Name,Age) values(:Name,:Age) [{Name:"jim" , Age:12}]
*/
func (this *BetterDB) Batch(s string, values []interface{}) {

}
