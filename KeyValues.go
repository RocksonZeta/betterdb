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

func (this *KeyValues) Get(key string) (result interface{}, ok bool) {
	if nil == this.Data {
		ok = false
		return
	}
	ok = true
	switch this.Data.(type) {
	case map[string]interface{}:
		result = this.Data.(map[string]interface{})[key]
		return
	case *map[string]interface{}:
		result = (*this.Data.(*map[string]interface{}))[key]
		return
	}
	value := reflect.ValueOf(this.Data)
	for "ptr" == value.Kind().String() {
		value = value.Elem()
	}
	if "struct" == value.Kind().String() {
		result = value.FieldByName(key).Interface()
		return
	}
	if "map" == value.Kind().String() {
		result = value.MapIndex(reflect.ValueOf(key)).Interface()
		return
	}
	ok = false
	return
}

func (this *KeyValues) Set(key string, val interface{}) (e error) {
	if nil == this.Data {
		e = &KeyValuesError{"KeyValues's data is nil"}
		return
	}
	switch this.Data.(type) {
	case map[string]interface{}:
		this.Data.(map[string]interface{})[key] = val
		return
	case *map[string]interface{}:
		(*this.Data.(*map[string]interface{}))[key] = val
		return
	}
	value := reflect.ValueOf(this.Data)
	for "ptr" == value.Kind().String() {
		value = value.Elem()
		return
	}
	if "map" == value.Kind().String() {
		(value.Interface().(map[string]interface{}))[key] = val
		return
	}
	if "struct" == value.Kind().String() {
		field := value.FieldByName(key)
		if field.CanSet() {
			field.Set(reflect.ValueOf(val))
		}
		return
	}
	panic(&KeyValuesNotSupport{reflect.TypeOf(this.Data).Kind().String()})
}

func (this *KeyValues) Pick(keys ...string) map[string]interface{} {
	if nil == this.Data {
		return nil
	}
	r := map[string]interface{}{}
	switch this.Data.(type) {
	case map[string]interface{}:
		m := this.Data.(map[string]interface{})
		for _, k := range keys {
			r[k] = m[k]
		}
		return r
	case *map[string]interface{}:
		m := *this.Data.(*map[string]interface{})
		for _, k := range keys {
			r[k] = m[k]
		}
		return r
	}

	value := reflect.ValueOf(this.Data)
	for "ptr" == value.Kind().String() {
		value = value.Elem()
	}
	if "map" == value.Kind().String() {
		m := value.Interface().(map[string]interface{})
		for _, k := range keys {
			r[k] = m[k]
		}
	}
	if "struct" == value.Kind().String() {
		var field reflect.Value
		for _, k := range keys {
			field = value.FieldByName(k)
			if !field.IsValid() {
				continue
			}
			r[k] = field.Interface()
		}
		return r
	}
	return nil
}

func (this *KeyValues) Map() map[string]interface{} {
	if nil == this.Data {
		return nil
	}
	switch this.Data.(type) {
	case map[string]interface{}:
		return this.Data.(map[string]interface{})
	case *map[string]interface{}:
		return (*this.Data.(*map[string]interface{}))
	}
	value := reflect.ValueOf(this.Data)
	for "ptr" == value.Kind().String() {
		value = value.Elem()
	}
	if "map" == value.Kind().String() {
		return value.Interface().(map[string]interface{})
	}
	if "struct" == value.Kind().String() {
		t := value.Type()
		r := map[string]interface{}{}
		var fieldName string
		for i := 0; i < t.NumField(); i++ {
			fieldName = t.Field(i).Name
			if 97 <= fieldName[0] {
				continue
			}
			r[fieldName] = value.FieldByName(fieldName).Interface()
		}
		return r
	}
	panic(&KeyValuesNotSupport{reflect.TypeOf(this.Data).Kind().String()})

}
