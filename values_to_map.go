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

// 构建map, 支持 map[string]string，map[string]interface{}, KV, KVs
// 其它值按顺序转为 map[string]interface{}{"*[0]": 值0, "*[1]", 值1...}
func MakeMapOfValues(values ...interface{}) map[string]interface{} {
	return makeMapOfValues(values)
}

// 构建map, 支持 map[string]string，map[string]interface{}, KV, KVs
// 其它值按顺序转为 map[string]interface{}{"*[0]": 值0, "*[1]", 值1...}
func makeMapOfValues(values []interface{}) map[string]interface{} {
	if len(values) == 0 {
		return make(map[string]interface{})
	}

	// map, kvs, kv
	switch v := values[0].(type) {
	case map[string]interface{}:
		return v
	case KVs:
		data := make(map[string]interface{}, len(v))
		for _, kv := range v {
			data[kv.K] = kv.V
		}
		return data
	case KV:
		data := make(map[string]interface{}, len(values))
		for _, value := range values {
			kv, ok := value.(KV)
			if ok {
				panic("所有值必须都是 zstr.KV")
			}
			data[kv.K] = kv.V
		}
		return data
	case *KV:
		data := make(map[string]interface{}, len(values))
		for _, value := range values {
			kv, ok := value.(*KV)
			if ok {
				panic("所有值必须都是 *zstr.KV")
			}
			data[kv.K] = kv.V
		}
		return data
	}

	// 其他map
	rv := reflect.Indirect(reflect.ValueOf(values[0]))
	switch rv.Kind() {
	case reflect.Map:
		data := make(map[string]interface{}, rv.Len())
		for iter := rv.MapRange(); iter.Next(); {
			data[anyToString(iter.Key().Interface())] = iter.Value().Interface()
		}
		return data
	}

	// values
	data := make(map[string]interface{}, len(values))
	for i, v := range values {
		data[`*[`+strconv.Itoa(i)+`]`] = v
	}
	return data
}
