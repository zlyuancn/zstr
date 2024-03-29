/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/17
   Description :
-------------------------------------------------
*/

package zstr

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// 任何值转字符串
func AnyToStr(a interface{}, nilToEmpty ...bool) string {
	return anyToString(a, nilToEmpty...)
}

// 任何值转字符串
func anyToString(a interface{}, nilToEmpty ...bool) string {
	switch v := a.(type) {

	case nil:
		if len(nilToEmpty) > 0 && nilToEmpty[0] {
			return ""
		}
		return "nil"

	case string:
		return v
	case []byte:
		return *BytesToString(v)
	case bool:
		if v {
			return "true"
		}
		return "false"

	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)

	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	}
	return fmt.Sprint(a)
}

// 任何值转sql需要的字符串
func AnyToSqlStr(a interface{}, str_crust ...bool) string {
	return anyToSqlString(a, len(str_crust) > 0 && str_crust[0])
}

// 任何值转sql需要的字符串
func anyToSqlString(a interface{}, str_crust bool) string {
	switch v := a.(type) {

	case nil:
		return "null"

	case string:
		if str_crust {
			return `'` + v + `'`
		}
		return v
	case []byte:
		if str_crust {
			return `'` + *BytesToString(v) + `'`
		}
		return *BytesToString(v)
	case bool:
		if v {
			return "true"
		}
		return "false"

	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)

	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	}

	r_v := reflect.Indirect(reflect.ValueOf(a))
	if r_v.Kind() != reflect.Slice && r_v.Kind() != reflect.Array {
		return fmt.Sprint(a)
	}

	l := r_v.Len()
	ss := make([]string, l)
	for i := 0; i < l; i++ {
		ss[i] = anyToSqlString(reflect.Indirect(r_v.Index(i)).Interface(), str_crust)
	}
	return `(` + strings.Join(ss, ", ") + `)`
}
