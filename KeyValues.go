package betterdb

import "reflect"

const (
	Type_Map = iota
	Type_MapPoint
	Type_Struct
)

type KeyValuesError struct {
	Msg string
}

func (this *KeyValuesError) Error() string {
	return this.Msg
}

type KeyValuesNotSupport struct {
	Type string
}

func (this *KeyValuesNotSupport) Error() string {
	return "Not support for this type:" + this.Type
}

/**
represent for struct and map or the ref of struct and map

*/
type KeyValues struct {
	Data interface{}
}

func (this *KeyValues) Get(key string) interface{} {
	if nil == this.Data {
		panic(&KeyValuesError{"KeyValues's data is nil"})
	}
	switch this.Data.(type) {
	case map[string]interface{}:
		return this.Data.(map[string]interface{})[key]
	case *map[string]interface{}:
		return (*this.Data.(*map[string]interface{}))[key]
	}
	value := reflect.ValueOf(this.Data)
	for "ptr" == value.Kind().String() {
		value = value.Elem()
	}
	if "struct" == value.Kind().String() {
		return value.FieldByName(key).Interface()
	}
	if "map" == value.Kind().String() {
		return value.MapIndex(reflect.ValueOf(key)).Interface()
	}
	panic(&KeyValuesNotSupport{reflect.TypeOf(this.Data).Kind().String()})
}

func (this *KeyValues) Set(key string, value interface{}) {

}

func (this *KeyValues) Pick(keys ...string) map[string]interface{} {
	return nil
}

func (this *KeyValues) Map() map[string]interface{} {
	return nil
}
