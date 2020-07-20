/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/18
   Description :
-------------------------------------------------
*/

package zstr

import (
	"fmt"
	"regexp"
	"strings"
)

const defaultSqlCompareFlag = "="

var sqlTemplateRegex = regexp.MustCompile(`[&|@]\w+`)
var sqlTemplateRegexCrust = regexp.MustCompile(`{[&|@]\w+}`)
var sqlTemplateRegexCrustAndFlag = regexp.MustCompile(`{[&|]\w+ .+?}`)

// sql模板渲染
//
// 语法格式1:   (操作符)(name)   示例:   &a   |a
// 语法格式2:   {(操作符)(name)}   示例:   {&a}   {|a}
// 语法格式3:   {(操作符)(name) (对比标志)}   示例:   {&a in}   {|a >}
//
// 操作符支持:
//     @: 直接赋值, 这个操作符仅不支持   {(操作符)(name) (对比标志)}   格式
//     &: 转为 and
//     |: 转为 or
//
// 对比标志支持:   >   >=   <   <=   !=   <>   =   in   not in   like   likestart    like_start   likeend   like_end
//
// 输入的kvs必须为：map[string]string, map[string]interface{}, 或健值对
//
// 注意:
//     如果name没有传参, 则替换为空字符串
//     如果name的值为nil, 则结果为: (操作符) (name) is null
//     如果name的值是一个切片, 结果会用逗号连接起来且外面会加上小括号. 如 []string{"a", "b"} 会转为 ("a", "b")
//
// 示例:
//    s := SqlTemplateRender("select * from t where &a {&b} {&c !=} {&d in} {|e} limit 1", map[string]interface{}{
//		"a": 1,
//		"b": "2",
//		"c": 3.3,
//		"d": []string{"4"},
//		"e": nil,
//	  })
func SqlTemplateRender(sql_template string, kvs ...interface{}) string {
	data := makeMapOfkvs(kvs)
	result := sqlTemplateRegexCrust.ReplaceAllStringFunc(sql_template, func(s string) string {
		return sqlTranslate(s[1:len(s)-1], defaultSqlCompareFlag, true, data)
	})
	result = sqlTemplateRegexCrustAndFlag.ReplaceAllStringFunc(result, func(s string) string {
		k := strings.Index(s, " ")
		return sqlTranslate(s[1:k], s[k+1:len(s)-1], true, data)
	})
	result = sqlTemplateRegex.ReplaceAllStringFunc(result, func(s string) string {
		return sqlTranslate(s, defaultSqlCompareFlag, false, data)
	})
	return repairSql(result)
}

func sqlTranslate(text, flag string, crust bool, m map[string]interface{}) string {
	operation, name, cflag := text[:1], text[1:], strings.ToLower(flag)

	value, has := m[name]

	switch operation {
	case "@":
		if has {
			return anyToSqlString(value)
		}
		if crust {
			return ""
		}
		return text
	case "&":
		operation = "and"
	case "|":
		operation = "or"
	default:
		panic(fmt.Errorf(`syntax error, non-supported operation "%s"`, operation))
	}

	var out string
	switch cflag {
	case ">", ">=", "<", "<=", "!=", "<>", "=":
		out = fmt.Sprintf(`%s %s %s %s`, operation, name, cflag, anyToSqlString(value))
	case "in", "not in":
		out = fmt.Sprintf(`%s %s %s %s`, operation, name, cflag, anyToSqlString(value))
	case "like": // 包含xx
		out = fmt.Sprintf(`%s %s like "%%%s%%"`, operation, name, anyToSqlString(value))
	case "likestart", "like_start": // 以xx开始
		out = fmt.Sprintf(`%s %s like "%s%%"`, operation, name, anyToSqlString(value))
	case "likeend", "like_end": // 以xx结束
		out = fmt.Sprintf(`%s %s like "%%%s"`, operation, name, anyToSqlString(value))
	default:
		panic(fmt.Errorf(`syntax error, non-supported flag "%s"`, flag))
	}

	if !has {
		return ""
	}
	if value == nil {
		return fmt.Sprintf(`%s %s is null`, operation, name)
	}
	return out
}
