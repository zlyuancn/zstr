/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/17
   Description :
-------------------------------------------------
*/

package zstr

import (
	"reflect"
	"strconv"
)

// 构建map, 支持 map[string]string，map[string]interface{}
// 其它值按顺序转为 map[string]interface{}{"*[0]": 值0, "*[1]", 值1...}
func MakeMapOfValues(values ...interface{}) map[string]interface{} {
	return makeMapOfValues(values)
}

// 构建map, 支持 map[string]string，map[string]interface{}
// 其它值按顺序转为 map[string]interface{}{"*[0]": 值0, "*[1]", 值1...}
func makeMapOfValues(values []interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	if len(values) == 0 {
		return data
	}

	rv := reflect.Indirect(reflect.ValueOf(values[0]))
	switch rv.Kind() {
	case reflect.Map:
		for iter := rv.MapRange(); iter.Next(); {
			data[anyToString(iter.Key().Interface())] = iter.Value().Interface()
		}
		return data
	}

	for i, v := range values {
		data[`*[`+strconv.Itoa(i)+`]`] = v
	}
	return data
}
