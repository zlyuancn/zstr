/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/17
   Description :
-------------------------------------------------
*/

package zstr

func makeMapOfkv(kvs []interface{}) map[string]interface{} {
	if len(kvs) == 0 {
		panic("输入的kv必须为：map[string]string，map[string]interface{}，或健值对")
	}

	var data = make(map[string]interface{})

	switch p := kvs[0].(type) {
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

	if len(kvs)&1 != 0 {
		panic("输入的kv必须为2的倍数")
	}
	for i := 0; i < len(kvs)-1; i += 2 {
		data[anyToString(kvs[i])] = anyToString(kvs[i+1])
	}
	return data
}
