/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/15
   Description :
-------------------------------------------------
*/

package zstr

import (
	"regexp"
)

var templateRegex = regexp.MustCompile(`@\w+`)
var templateRegexCrust = regexp.MustCompile(`\{@\w+\}`)

// 模板渲染
//
// 输入的kv必须为：map[string]string，map[string]interface{}，或健值对
// 示例:
//    s:=TemplateRender("s@a e", "a", "v")
//    s:=TemplateRender("s{@a}e", "a", "v")
//    s:=TemplateRender("s{@a}e", map[string]string{"a": "v"})
func TemplateRender(format string, kvs ...interface{}) string {
	data := makeMapOfkv(kvs)
	result := templateRegexCrust.ReplaceAllStringFunc(format, func(s string) string {
		v, ok := data[s[2:len(s)-1]]
		if ok {
			return anyToString(v)
		}
		return s
	})
	result = templateRegex.ReplaceAllStringFunc(result, func(s string) string {
		v, ok := data[s[1:]]
		if ok {
			return anyToString(v)
		}
		return s
	})
	return result
}
