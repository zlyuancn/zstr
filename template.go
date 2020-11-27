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
	"strconv"
)

var templateRegex = regexp.MustCompile(`(\{@\w*\.?\w+\})|(@\w*\.?\w+)`)

type simpleTemplate struct {
	data    map[string]interface{}
	counter counter
}

func newSimpleTemplate(kvs ...interface{}) *simpleTemplate {
	return &simpleTemplate{
		data:    makeMapOfkvs(kvs),
		counter: newCounter(),
	}
}

func (m *simpleTemplate) Render(format string) string {
	// 替换 {@field} 和 @field, 如果没有设置则不替换
	result := templateRegex.ReplaceAllStringFunc(format, func(s string) string {
		var key string
		var crust bool
		if s[0] == '{' {
			key = s[2 : len(s)-1]
			crust = true
		} else {
			key = s[1:]
		}

		v, ok := m.data[key+"["+strconv.Itoa(m.counter.Incr(key)-1)+"]"]
		if !ok {
			v, ok = m.data[key]
		}
		if ok {
			return anyToString(v)
		}
		if crust {
			return ""
		}
		return s
	})
	return result
}

// 模板渲染
//
// 输入的kvs必须为：map[string]string，map[string]interface{}，或健值对
// 示例:
//    s:=TemplateRender("s@a e", "a", "v")
//    s:=TemplateRender("s{@a}e", "a", "v")
//    s:=TemplateRender("s{@a}e", map[string]string{"a": "v"})
//    s:=TemplateRender("s@a @a e", "a", "v", "a[1]", "xxx")
func TemplateRender(format string, kvs ...interface{}) string {
	return newSimpleTemplate(kvs...).Render(format)
}

// 模板渲染, 和TemplateRender一样, 只是简短了函数名
func Render(format string, kvs ...interface{}) string {
	return newSimpleTemplate(kvs...).Render(format)
}
