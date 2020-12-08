/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/17
   Description :
-------------------------------------------------
*/

package zstr

import (
	"strconv"
)

// 构建map, 支持 map[string]string，map[string]interface{}
// 其它值按顺序转为 map[string]interface{}{"*[0]": 值0, "*[1]", 值1...}
func MakeMapOfKvs(a ...interface{}) map[string]interface{} {
	return makeMapOfkvs(a)
}

// 构建map, 支持 map[string]string，map[string]interface{}
// 其它值按顺序转为 map[string]interface{}{"*[0]": 值0, "*[1]", 值1...}
func makeMapOfkvs(a []interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	if len(a) == 0 {
		return data
	}

	switch p := a[0].(type) {
	case map[string]string:
		for k, v := range p {
			data[k] = v
		}
		return data
	case map[string]interface{}:
		for k, v := range p {
			data[k] = v
		}
		return data
	}

	for i, v := range a {
		data[`*[`+strconv.Itoa(i)+`]`] = v
	}
	return data
}
