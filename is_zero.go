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
	"unsafe"
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
		return !v

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

	rv := reflect.Indirect(reflect.ValueOf(a))

	switch rv.Kind() {
	case reflect.Array:
		return arrayIsZero(rv)
	case reflect.String:
		return rv.Len() == 0
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return rv.IsNil()
	case reflect.Struct:
		return structIsZero(rv)
	}

	nv := reflect.New(rv.Type()).Elem().Interface()
	return rv.Interface() == nv
}

func structIsZero(r_v reflect.Value) bool {
	num := r_v.NumField()
	for i := 0; i < num; i++ {
		field := r_v.Field(i)
		if field.Kind() == reflect.Invalid {
			continue
		}

		// 尝试获取值
		if field.CanInterface() {
			switch field.Kind() {
			case reflect.Ptr, reflect.Interface:
				if field.Interface() != nil {
					return false
				}
			default:
				if !IsZero(field.Interface()) {
					return false
				}
			}
			continue
		}

		var temp reflect.Value
		// 尝试获取指针
		if field.CanAddr() {
			temp = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr()))
		} else {
			// 强行获取数据
			rv := reflect.ValueOf(&field).Elem().Field(1).UnsafeAddr() // &field.ptr
			rv = *(*uintptr)(unsafe.Pointer(rv))                       // field.ptr
			temp = reflect.NewAt(field.Type(), unsafe.Pointer(rv))
		}

		switch field.Kind() {
		case reflect.Ptr, reflect.Interface:
			if temp.Elem().Interface() != nil {
				return false
			}
		default:
			if !IsZero(temp.Elem().Interface()) {
				return false
			}
		}
	}
	return true
}

func arrayIsZero(rv reflect.Value) bool {
	num := rv.Len()
	for i := 0; i < num; i++ {
		value := rv.Index(i)
		switch value.Kind() {
		case reflect.Ptr, reflect.Interface:
			if value.Interface() != nil {
				return false
			}
			continue
		}
		if !IsZero(value.Interface()) {
			return false
		}
	}
	return true
}
