/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/18
   Description :
-------------------------------------------------
*/

package zstr

import (
	"bytes"
)

var emptyStrMap = map[int32]struct{}{
	'\f': {},
	'\n': {},
	'\r': {},
	'\t': {},
	'\v': {},
	' ':  {},
	0x85: {},
	0xA0: {},
}

// 缩进空格
func (m *sqlTemplate) retractAllSpace(s string) string {
	var buff bytes.Buffer
	old := -2
	var ok bool
	for i, v := range []rune(s) {
		if _, ok = emptyStrMap[v]; ok {
			if i-old > 1 {
				buff.WriteByte(' ')
			}
			old = i
		} else {
			buff.WriteRune(v)
		}
	}

	return buff.String()
}

// 修复模板渲染后无效的sql语句
func (m *sqlTemplate) repairSql(sql string) string {
	result := m.retractAllSpace(sql)

	if m.ContainsIgnoreCase(result, "where") {
		result = m.ReplaceAllIgnoreCase(result, "where or ", "where ")
		result = m.ReplaceAllIgnoreCase(result, "where and ", "where ")
		result = m.ReplaceAllIgnoreCase(result, "where order by", "order by")
		result = m.ReplaceAllIgnoreCase(result, "where group by", "group by")
		result = m.ReplaceAllIgnoreCase(result, "where limit", "limit")
		result = m.ReplaceAllIgnoreCase(result, "where )", ")")
		result = m.ReplaceAllIgnoreCase(result, "( and", "(")
		result = m.ReplaceAllIgnoreCase(result, "(and", "(")
		result = m.ReplaceAllIgnoreCase(result, "( or", "(")
		result = m.ReplaceAllIgnoreCase(result, "(or", "(")
		result = m.ReplaceAllIgnoreCase(result, "( )", "")
		result = m.ReplaceAllIgnoreCase(result, "()", "")
		if m.HasSuffixIgnoreCase(result, "where ") {
			result = result[:len(result)-6]
		}
		if m.HasSuffixIgnoreCase(result, "where ;") {
			result = result[:len(result)-7] + ";"
		}
	}
	return result
}

// 忽略大小写检查字符相等
func (m *sqlTemplate) EqualCharIgnoreCase(c1, c2 int32) bool {
	if c1 == c2 {
		return true
	}
	switch c1 - c2 {
	case 32: // a - A
		return c1 >= 'a' && c1 <= 'z'
	case -32: // A - a
		return c1 >= 'A' && c1 <= 'Z'
	}
	return false
}

// 忽略大小写检查文本相等
func (m *sqlTemplate) EqualIgnoreCase(s1, s2 string) bool {
	if s1 == s2 {
		return true
	}

	r1 := []rune(s1)
	r2 := []rune(s2)
	if len(r1) != len(r2) {
		return false
	}

	for i, r := range r1 {
		if !m.EqualCharIgnoreCase(r, r2[i]) {
			return false
		}
	}
	return true
}

// 忽略大小写替换所有文本
func (m *sqlTemplate) ReplaceAllIgnoreCase(s, old, new string) string {
	return m.ReplaceIgnoreCase(s, old, new, -1)
}

// 替换n次忽略大小写匹配的文本
func (m *sqlTemplate) ReplaceIgnoreCase(s, old, new string, n int) string {
	if n == 0 || old == new || old == "" {
		return s
	}

	ss := []rune(s)
	sub := []rune(old)
	var buff bytes.Buffer
	var num int
	for offset := 0; offset < len(ss); {
		start := m.searchIgnoreCase(ss, sub, offset)
		if start > -1 {
			buff.WriteString(string(ss[offset:start]))
			buff.WriteString(new)
			offset = start + len(sub)
			num++
		}

		if start == -1 || num == n {
			buff.WriteString(string(ss[offset:]))
			break
		}
	}
	return buff.String()
}

// 忽略大小写查找第一个匹配sub的文本所在位置, 如果不存在返回-1
func (m *sqlTemplate) searchIgnoreCase(ss []rune, sub []rune, start int) int {
	if len(ss)-start < len(sub) {
		return -1
	}

	var has bool
	// 查找开头
	for i := start; i < len(ss); i++ {
		if m.EqualCharIgnoreCase(ss[i], sub[0]) {
			start, has = i, true
			break
		}
	}
	if !has {
		return -1
	}
	for i := 1; i < len(sub); i++ {
		if !m.EqualCharIgnoreCase(ss[start+i], sub[i]) {
			return m.searchIgnoreCase(ss, sub, start+1)
		}
	}
	return start
}

// 忽略大小写查找第一个匹配sub的文本所在位置, 如果不存在返回-1
func (m *sqlTemplate) IndexIgnoreCase(s, sub string) int {
	return m.searchIgnoreCase([]rune(s), []rune(sub), 0)
}

// 忽略大小写查找s是否包含sub
func (m *sqlTemplate) ContainsIgnoreCase(s, sub string) bool {
	return m.IndexIgnoreCase(s, sub) >= 0
}

// 忽略大小写测试文本s是否以suffix结束
func (m *sqlTemplate) HasSuffixIgnoreCase(s, suffix string) bool {
	return len(s) >= len(suffix) && m.EqualIgnoreCase(s[len(s)-len(suffix):], suffix)
}
