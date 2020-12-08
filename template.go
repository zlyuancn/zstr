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
F1:
	// 查找开头
	for i := start; i < len(ss); i++ {
		switch ss[i] {
		case '{':
			start, crust, has = i, true, true
			break F1
		case '@':
			start, crust, has = i, false, true
			break F1
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
			if !ok {
				if i-start < 2 || ss[start+1] == '.' || ss[i-1] == '.' {
					return m.calculateTemplate(ss, i)
				}
				return start, i, false, true
			}
		}
		return start, len(ss), false, len(ss)-start >= 2 && ss[start+1] != '.' && ss[len(ss)-1] != '.'
	}

	if crust {
		if ss[start+1] != '@' {
			has = false
		}
		for i := start + 2; i < len(ss); i++ {
			if ss[i] != '}' {
				_, ok = templateVariableNameMap[ss[i]]
				if !ok {
					has = false
				}
				continue
			}
			if i-start < 3 || !has || ss[start+2] == '.' || ss[i-1] == '.' {
				return m.calculateTemplate(ss, i+1)
			}
			return start, i + 1, true, true
		}
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
// 输入的values必须为：map[string]string，map[string]interface{}，或按顺序传入值
// 示例:
//    s:=Render("s@a e", map[string]string{"a": "va"})
//    s:=Render("s{@a}e", map[string]string{"a": "va"})
//    s:=Render("s{@a}e", "va")
//    s:=Render("s@a @a e", "va0", "va1")
//
// 寻值优先级:
//    匹配名下标 > 匹配名 > *下标
//    如:  a[0] > a > *[0]
//
// 注意:
//     如果name存在花括号外壳{}且没有传参, 则替换为空字符串
//     我们不会去检查name是否完全符合变量名标志, 因为这是无意义且消耗资源的
//         变量名首位可以为数字, 变量中间可以连续出现多个小数点, 如 0..a 是合法的
func TemplateRender(format string, values ...interface{}) string {
	return newSimpleTemplate(values...).Render(format)
}

// 模板渲染, 和TemplateRender一样, 只是简短了函数名
func Render(format string, values ...interface{}) string {
	return newSimpleTemplate(values...).Render(format)
}
