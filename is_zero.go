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
		return v == nil
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

	case complex64:
		return v == 0
	case complex128:
		return v == 0
	}

	rv := reflect.Indirect(reflect.ValueOf(a)) // 解包ptr
	return reflectValueIsZero(rv)
}

func reflectValueIsZero(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Array:
		return arrayIsZero(rv)
	case reflect.String:
		return rv.Len() == 0
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return rv.IsNil()
	case reflect.Struct:
		return structIsZero(rv)
	}

	nv := reflect.New(rv.Type()).Elem().Interface() // 根据类型创建新的数据

	// 尝试获取值
	if rv.CanInterface() {
		return rv.Interface() == nv
	}

	var p uintptr
	if rv.CanAddr() { // 尝试获取指针
		p = rv.UnsafeAddr()
	} else {
		// 强行获取指针
		p = reflect.ValueOf(&rv).Elem().Field(1).UnsafeAddr() // &rv.ptr
		p = *(*uintptr)(unsafe.Pointer(p))                    // rv.ptr
	}

	temp := reflect.NewAt(rv.Type(), unsafe.Pointer(p)) // 根据指针创建新的数据
	return temp.Elem().Interface() == nv
}

func structIsZero(rv reflect.Value) bool {
	num := rv.NumField()
	for i := 0; i < num; i++ {
		if !reflectValueIsZero(rv.Field(i)) {
			return false
		}
	}
	return true
}

func arrayIsZero(rv reflect.Value) bool {
	num := rv.Len()
	for i := 0; i < num; i++ {
		if !reflectValueIsZero(rv.Index(i)) {
			return false
		}
	}
	return true
}
