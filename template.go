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

func (m *simpleTemplate) calculateTemplate(ss []rune, start int) (int, int, string, bool, bool) {
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
		return 0, 0, "", false, false
	}

	// 预检
	if crust && (len(ss)-start < 4) || (len(ss)-start < 2) {
		return 0, 0, "", false, false
	}

	var ok bool
	if !crust {
		for i := start + 1; i < len(ss); i++ {
			_, ok = templateVariableNameMap[ss[i]]
			if !ok { // 表示查找变量结束了
				if i-start < 2 || ss[start+1] == '.' || ss[i-1] == '.' { // 操作符占一个位置, 变量长度不可能为0
					return m.calculateTemplate(ss, i)
				}
				return start, i, string(ss[start+1 : i]), false, true // 中间的数据就是需要的数据
			}
		}
		// 可能整个字符串都是需要的数据
		return start, len(ss), string(ss[start+1:]), false, len(ss)-start >= 2 && ss[start+1] != '.' && ss[len(ss)-1] != '.'
	}

	// 以下包含{
	var variableStart, variableEnd, end int
	for i := start + 1; i < len(ss); i++ {
		if variableStart == 0 {
			if ss[i] == '@' {
				variableStart = i + 1
				continue
			}
			if ss[i] == ' ' {
				continue
			}
			// {}中间出现非预期的字符, 从这里开始重新扫描
			return m.calculateTemplate(ss, i)
		}

		if variableEnd == 0 {
			if ss[i] == '}' {
				variableEnd = i
				end = i + 1
				break
			}
			if ss[i] == ' ' {
				variableEnd = i
				continue
			}
			_, ok = templateVariableNameMap[ss[i]]
			if ok {
				continue
			}
			// {}中间出现非预期的字符, 从这里开始重新扫描
			return m.calculateTemplate(ss, i)
		}

		if ss[i] == '}' {
			end = i + 1
			break
		}
		if ss[i] == ' ' {
			continue
		}
		// {}中间出现非预期的字符, 从这里开始重新扫描
		return m.calculateTemplate(ss, i)
	}
	if end == 0 || variableStart >= variableEnd || ss[variableStart] == '.' || ss[variableEnd-1] == '.' {
		return 0, 0, "", false, false
	}

	return start, end, string(ss[variableStart:variableEnd]), true, true
}

func (m *simpleTemplate) replaceAllFunc(s string, fn func(full string, variable string, crust bool) string) string {
	ss := []rune(s)
	var buff bytes.Buffer
	for offset := 0; offset < len(ss); {
		start, end, variable, crust, has := m.calculateTemplate(ss, offset)
		if !has {
			buff.WriteString(string(ss[offset:]))
			break
		}

		buff.WriteString(string(ss[offset:start]))
		if crust {
			buff.WriteString(fn(string(ss[start:end]), variable, crust))
		} else {
			buff.WriteString(fn(string(ss[start:end]), variable, crust))
		}
		offset = end
	}
	return buff.String()
}

func (m *simpleTemplate) Render(format string) string {
	// 替换 {@field} 和 @field, 如果没有设置则不替换
	result := m.replaceAllFunc(format, func(full string, variable string, crust bool) string {
		v, ok := m.data[variable+"["+strconv.Itoa(m.keyCounter.Incr(variable))+"]"]
		if !ok {
			v, ok = m.data[variable]
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
		return full
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

	zstr.Render("{@a} { @b } @c text", "va", "vb", "vc") // 按顺序赋值
	zstr.Render("@a text", map[string]string{"a": "va"}) // 指定变量名赋值
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
