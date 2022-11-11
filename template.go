/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/15
   Description :
-------------------------------------------------
*/

package zstr

import (
	"bytes"
	"strconv"
)

type simpleTemplate struct {
	data       map[string]interface{}
	keyCounter *counter // key计数器
	sub        int      // 下标计数器
}

func newSimpleTemplate(values ...interface{}) *simpleTemplate {
	return &simpleTemplate{
		data:       makeMapOfValues(values),
		keyCounter: newCounter(-1),
	}
}

var templateVariableNameMap = func() map[int32]struct{} {
	mm := map[int32]struct{}{
		'_': {}, '.': {},
	}
	for i := '0'; i < '9'+1; i++ {
		mm[i] = struct{}{}
	}
	for i := 'a'; i < 'z'+1; i++ {
		mm[i] = struct{}{}
	}
	for i := 'A'; i < 'Z'+1; i++ {
		mm[i] = struct{}{}
	}
	return mm
}()

func (m *simpleTemplate) calculateTemplate(ss []rune, start int) (int, int, bool, bool) {
	var crust, has bool
	// 查找开头
	for i := start; i < len(ss); i++ {
		if ss[i] == '{' {
			start, crust, has = i, true, true
			break
		}
		if ss[i] == '@' {
			start, crust, has = i, false, true
			break
		}
	}
	if !has {
		return 0, 0, false, false
	}

	// 预检
	if crust && (len(ss)-start < 4) || (len(ss)-start < 2) {
		return 0, 0, false, false
	}

	var ok bool
	if !crust {
		for i := start + 1; i < len(ss); i++ {
			_, ok = templateVariableNameMap[ss[i]]
			if !ok { // 表示查找变量结束了
				if i-start < 2 || ss[start+1] == '.' || ss[i-1] == '.' { // 操作符占一个位置, 变量长度不可能为0
					return m.calculateTemplate(ss, i)
				}
				return start, i, false, true // 中间的数据就是需要的数据
			}
		}
		// 可能整个字符串都是需要的数据
		return start, len(ss), false, len(ss)-start >= 2 && ss[start+1] != '.' && ss[len(ss)-1] != '.'
	}

	// 以下包含{
	if ss[start+1] != '@' {
		return m.calculateTemplate(ss, start+1)
	}
	for i := start + 2; i < len(ss); i++ {
		if ss[i] != '}' {
			_, ok = templateVariableNameMap[ss[i]]
			if ok {
				continue
			}
			// {}中间出现非预期的字符, 从这里开始重新扫描
			return m.calculateTemplate(ss, i)
		}
		// 这里是}结束标志
		if i-start < 3 || !has || ss[start+2] == '.' || ss[i-1] == '.' {
			return m.calculateTemplate(ss, i+1)
		}
		return start, i + 1, true, true
	}
	return 0, 0, false, false
}

func (m *simpleTemplate) replaceAllFunc(s string, fn func(s string, crust bool) string) string {
	ss := []rune(s)
	var buff bytes.Buffer
	for offset := 0; offset < len(ss); {
		start, end, crust, has := m.calculateTemplate(ss, offset)
		if !has {
			buff.WriteString(string(ss[offset:]))
			break
		}

		buff.WriteString(string(ss[offset:start]))
		buff.WriteString(fn(string(ss[start:end]), crust))
		offset = end
	}
	return buff.String()
}

func (m *simpleTemplate) Render(format string) string {
	// 替换 {@field} 和 @field, 如果没有设置则不替换
	result := m.replaceAllFunc(format, func(s string, crust bool) string {
		var key string
		if crust {
			key = s[2 : len(s)-1]
		} else {
			key = s[1:]
		}

		v, ok := m.data[key+"["+strconv.Itoa(m.keyCounter.Incr(key))+"]"]
		if !ok {
			v, ok = m.data[key]
		}
		if !ok {
			v, ok = m.data["*["+strconv.Itoa(m.sub)+"]"]
		}
		m.sub++ // 每次一定+1
		if ok {
			return anyToString(v, true)
		}
		if crust {
			return ""
		}
		return s
	})
	return result
}

// 模板渲染, 和Render一样, 只是加长了函数名
func TemplateRender(format string, values ...interface{}) string {
	return newSimpleTemplate(values...).Render(format)
}

/*
模板渲染

示例:

	zstr.Render("{@a}e", "va") // 按顺序赋值
	zstr.Render("@a @a e", "va0", "va1") // 按顺序赋值
	zstr.Render("@a e", map[string]string{"a": "va"}) // 指定变量名赋值
	zstr.Render("@a @b @c", zstr.KV{"a", "aValue"}, zstr.KV{"b", "bValue"}, zstr.KV{"c", "cValue"}) // 指定变量名赋值
	zstr.Render("@a @b @c", zstr.KVs{{"a", "aValue"}, {"b", "bValue"}, {"c", "cValue"}}) // 指定变量名赋值
	zstr.Render("@a @a @a", zstr.KV{"a[0]", "1"}, zstr.KV{"a", "2"}) // 指定下标, 指定变量名+下标的优先级比指定变量名更高

	type AA struct {
		A int
		B int `render:"b"`
	}
	s := zstr.Render(`@A @b`, AA{1, 2})

寻值优先级
 1. 变量名+下标
 2. 变量名
 3. 星号+下标
    如:  `a[0]` > `a` > `*[0]`

注意:

	如果未对模板中的变量进行赋值并且该变量被花括号`{}`包裹, 那么该会被替换为空字符串.
*/
func Render(format string, values ...interface{}) string {
	return newSimpleTemplate(values...).Render(format)
}
