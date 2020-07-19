/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/18
   Description :
-------------------------------------------------
*/

package zstr

import (
	"regexp"
	"strings"
)

var repairSqlRegs = []struct {
	re    *regexp.Regexp
	value string
}{
	{regexp.MustCompile(`where\s*order by`), "order by"},
	{regexp.MustCompile(`where\s*group by`), "group by"},
	{regexp.MustCompile(`where\s*limit`), "limit"},
	{regexp.MustCompile(`where\s*or$`), ""},
	{regexp.MustCompile(`where\s*and$`), ""},
	{regexp.MustCompile(`where\s*or\s*order by`), "order by"},
	{regexp.MustCompile(`where\s*or\s*group by`), "group by"},
	{regexp.MustCompile(`where\s*or\s*limit`), "limit"},
	{regexp.MustCompile(`where\s*or\s*having`), "having"},
	{regexp.MustCompile(`where\s*and\s*order by`), "order by"},
	{regexp.MustCompile(`where\s*and\s*group by`), "group by"},
	{regexp.MustCompile(`where\s*and\s*limit`), "limit"},
	{regexp.MustCompile(`where\s*and\s*having`), "having"},
	{regexp.MustCompile(`where\s*or\s+`), "where "},
	{regexp.MustCompile(`where\s*and\s+`), "where "},
	{regexp.MustCompile(`where\s*$`), ""},
	{regexp.MustCompile(`where\s*\)`), ")"},
}

// 修复模板渲染后无效的sql语句
func repairSql(sql string) string {
	var result = strings.ToLower(sql)
	for _, repair := range repairSqlRegs {
		result = repair.re.ReplaceAllString(result, repair.value)
	}
	return result
}
