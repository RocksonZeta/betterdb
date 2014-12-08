package betterdb

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

		//fields[i] = record.Field(i).Addr().Interface()
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

//variable placeholder should have this form ":var",eg."select name from user where name=:name"
func NamedQuery(db *sql.DB, s string, records, namedArgs map[string]interface{}) {
	re := regexp.MustCompile(":\\w+\\b")
	var args []interface{}
	st := re.ReplaceAllStringFunc(s, func(key string) string {
		fmt.Println(key)
		if n, ok := namedArgs[strings.TrimPrefix(key, ":")]; ok {
			args = append(args, n)
			return "?"
		}
		return key
	})
	fmt.Println(st, args)
	//Query(db, st, records, args)
}

type BetterDB struct {
	*sql.DB
}

func (this *BetterDB) BetterQuery(s string, records, args interface{}) {

}

func (this *BetterDB) NamedQuery(s string, records, args interface{}) {

}

func (this *BetterDB) Post(table string, obj interface{}) {

}
func (this *BetterDB) Get(table string, id interface{}, obj interface{}) {

}
func (this *BetterDB) Put(table string, id interface{}, newValues interface{}) {

}
func (this *BetterDB) Delete(table string, id interface{}) {

}
