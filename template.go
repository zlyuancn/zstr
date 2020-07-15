/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/15
   Description :
-------------------------------------------------
*/

package zstr

import (
	"fmt"
	"regexp"
)

var templateRegex = regexp.MustCompile(`@\w+`)
var templateRegexCrust = regexp.MustCompile(`\{@\w+\}`)

func anyToString(a interface{}) string {
	if a == nil {
		return "nil"
	}
	return fmt.Sprint(a)
}
func makeTemplateData(kvs []interface{}) map[string]string {
	if len(kvs) == 0 {
		panic("输入的kv必须为：map[string]string，map[string]interface{}，或健值对")
	}

	var data = make(map[string]string)

	switch p := kvs[0].(type) {
	case map[string]string:
		for k, v := range p {
			data[k] = v
		}
		return data
	case map[string]interface{}:
		for k, v := range p {
			data[k] = anyToString(v)
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

// 模板渲染
//
// 示例:
//    s:=TemplateRender("s@name e", "name", "v")
//    s:=TemplateRender("s@name e", "x", "v")
//    s:=TemplateRender("s@name e", "name", nil)
//    s:=TemplateRender("s@name e", "name", 2)
func TemplateRender(format string, kvs ...interface{}) string {
	data := makeTemplateData(kvs)
	result := templateRegexCrust.ReplaceAllStringFunc(format, func(s string) string {
		v, ok := data[s[2:len(s)-1]]
		if ok {
			return v
		}
		return s
	})
	result = templateRegex.ReplaceAllStringFunc(result, func(s string) string {
		v, ok := data[s[1:]]
		if ok {
			return v
		}
		return s
	})
	return result
}
