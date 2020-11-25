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

var templateRegex = regexp.MustCompile(`(\{@\w*\.?\w+\})|(@\w*\.?\w+)`)

// 模板渲染
//
// 输入的kvs必须为：map[string]string，map[string]interface{}，或健值对
// 示例:
//    s:=TemplateRender("s@a e", "a", "v")
//    s:=TemplateRender("s{@a}e", "a", "v")
//    s:=TemplateRender("s{@a}e", map[string]string{"a": "v"})
func TemplateRender(format string, kvs ...interface{}) string {
	data := makeMapOfkvs(kvs)
	// 替换 {@field} 和 @field, 如果没有设置则不替换
	result := templateRegex.ReplaceAllStringFunc(format, func(s string) string {
		var key string
		if s[0] == '{' {
			key = s[2 : len(s)-1]
		} else {
			key = s[1:]
		}

		v, ok := data[key]
		if ok {
			return anyToString(v)
		}
		return s
	})
	return result
}

// 模板渲染, 和TemplateRender一样, 只是简短了函数名
func Render(format string, kvs ...interface{}) string {
	return TemplateRender(format, kvs...)
}
