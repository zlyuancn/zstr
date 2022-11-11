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

// 结构体值在渲染时读取的标签, 否则用字段名
const structValueTag = `render`

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

	// 注意, nil也是一个有效的值

	if len(values) == 1 {
		// map, kvs
		switch v := values[0].(type) {
		case map[string]interface{}:
			return v
		case KVs:
			data := make(map[string]interface{}, len(v))
			for _, kv := range v {
				data[kv.K] = kv.V
			}
			return data
		}

		// 其他类型
		rv := reflect.Indirect(reflect.ValueOf(values[0]))
		switch rv.Kind() {
		case reflect.Map:
			data := make(map[string]interface{}, rv.Len())
			for iter := rv.MapRange(); iter.Next(); {
				data[anyToString(iter.Key().Interface())] = iter.Value().Interface()
			}
			return data
		case reflect.Struct:
			aType := rv.Type()
			fieldCount := aType.NumField()
			data := make(map[string]interface{}, fieldCount)
			for i := 0; i < fieldCount; i++ {
				field := aType.Field(i)
				if field.PkgPath != "" {
					continue
				}
				key := field.Tag.Get(structValueTag)
				if key == "" {
					key = field.Name
				}
				data[key] = rv.Field(i).Interface()
			}
			return data
		}

	}

	// map, kvs
	switch values[0].(type) {
	case KV, *KV, KVs:
		data := make(map[string]interface{}, len(values))
		for _, value := range values {
			if kv, ok := value.(KV); ok {
				data[kv.K] = kv.V
				continue
			}
			if kv, ok := value.(*KV); ok {
				data[kv.K] = kv.V
				continue
			}
			if kvs, ok := value.(KVs); ok {
				for _, kv := range kvs {
					data[kv.K] = kv.V
				}
				continue
			}

			panic("所有值必须都是 zstr.KV 或者 *zstr.KV 或者 zstr.KVs")
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
