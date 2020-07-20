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
func anyToString(a interface{}) string {
	switch v := a.(type) {

	case nil:
		return "nil"

	case string:
		return v
	case []byte:
		return string(v)
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
func anyToSqlString(a interface{}, str_crust bool) string {
	switch v := a.(type) {

	case nil:
		return "NULL"

	case string:
		if str_crust {
			return `"` + v + `"`
		}
		return v
	case []byte:
		if str_crust {
			return `"` + string(v) + `"`
		}
		return string(v)
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
