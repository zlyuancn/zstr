/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/22
   Description :
-------------------------------------------------
*/

package zstr

import (
	"reflect"
)

// 判断传入参数是否为该类型的零值
func IsZero(a interface{}) bool {
	switch v := a.(type) {

	case nil:
		return true

	case string:
		return v == ""
	case []byte:
		return len(v) == 0
	case bool:
		return v

	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0

	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0

	case float32:
		return v == 0
	case float64:
		return v == 0
	}

	r_v := reflect.Indirect(reflect.ValueOf(a))

	switch r_v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return r_v.Len() == 0
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.UnsafePointer:
		return r_v.IsNil()
	}

	nv := reflect.New(r_v.Type()).Elem().Interface()
	return r_v.Interface() == nv
}
