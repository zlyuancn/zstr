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
	"strings"
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
	result = strings.ToLower(result)

	if strings.Contains(result, "where") {
		result = strings.ReplaceAll(result, "where or ", "where ")
		result = strings.ReplaceAll(result, "where and ", "where ")
		result = strings.ReplaceAll(result, "where order by", "order by")
		result = strings.ReplaceAll(result, "where group by", "group by")
		result = strings.ReplaceAll(result, "where limit", "limit")
		result = strings.ReplaceAll(result, "where )", ")")
		if strings.HasSuffix(result, "where ") {
			result = result[:len(result)-6]
		}
		if strings.HasSuffix(result, "where ;") {
			result = result[:len(result)-7] + ";"
		}
	}
	return result
}
