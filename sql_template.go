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
	"strconv"
	"strings"
)

const defaultSqlCompareFlag = "="

var sqlTemplateRegex = regexp.MustCompile(`[&|@]\w+`)
var sqlTemplateRegexCrust = regexp.MustCompile(`{[&|@]\w+}`)
var sqlTemplateRegexCrustAndFlag = regexp.MustCompile(`{[&|]\w+ .+?}`)
var sqlTemplateParseNameRegex = regexp.MustCompile(`{{\d+}}`)

type sqlTemplate struct {
	data   map[string]interface{}
	index  uint64
	names  []string
	values []interface{}
}

func newSqlTemplate(kvs ...interface{}) *sqlTemplate {
	return &sqlTemplate{
		data: makeMapOfkvs(kvs),
	}
}

func (m *sqlTemplate) addValue(name string, value interface{}) (flag string) {
	flag = "{{" + strconv.FormatUint(m.index, 10) + "}}"
	m.names = append(m.names, name)
	m.values = append(m.values, value)
	m.index++
	return
}

func (m *sqlTemplate) Parse(sql_template string) (sql_str string, names []string, args []interface{}) {
	sql_str = sqlTemplateRegexCrust.ReplaceAllStringFunc(sql_template, func(s string) string {
		return m.translate(s[1:len(s)-1], defaultSqlCompareFlag, true)
	})

	sql_str = sqlTemplateRegexCrustAndFlag.ReplaceAllStringFunc(sql_str, func(s string) string {
		k := strings.Index(s, " ")
		return m.translate(s[1:k], s[k+1:len(s)-1], true)
	})

	sql_str = sqlTemplateRegex.ReplaceAllStringFunc(sql_str, func(s string) string {
		return m.translate(s, defaultSqlCompareFlag, false)
	})

	// 按顺序写入names和args
	sql_str = sqlTemplateParseNameRegex.ReplaceAllStringFunc(sql_str, func(s string) string {
		index, _ := strconv.Atoi(s[2 : len(s)-2])
		names = append(names, m.names[index])
		args = append(args, m.values[index])
		return "?"
	})

	sql_str = repairSql(sql_str)
	return sql_str, names, args
}

func (m *sqlTemplate) translate(text, flag string, crust bool) string {
	operation, name, cflag := text[:1], text[1:], strings.ToLower(flag)

	value, has := m.data[name]

	switch operation {
	case "@":
		if has {
			return m.addValue(name, value)
		}
		if crust {
			return "null"
		}
		panic(fmt.Sprintf(`"%s" must have a value`, text))
	case "&":
		operation = "and"
	case "|":
		operation = "or"
	default:
		panic(fmt.Errorf(`syntax error, non-supported operation "%s"`, operation))
	}

	var makeSqlStr func() string
	switch cflag {
	case ">", ">=", "<", "<=", "!=", "<>", "=":
		makeSqlStr = func() string {
			return fmt.Sprintf(`%s %s %s %s`, operation, name, cflag, m.addValue(name, value))
		}
	case "in", "not in":
		makeSqlStr = func() string {
			return fmt.Sprintf(`%s %s %s (%s)`, operation, name, cflag, m.addValue(name, value))
		}
	case "like": // 包含xx
		makeSqlStr = func() string {
			value = "%" + anyToSqlString(value, false) + "%"
			return fmt.Sprintf(`%s %s like %s`, operation, name, m.addValue(name, value))
		}
	case "likestart", "like_start": // 以xx开始
		makeSqlStr = func() string {
			value = anyToSqlString(value, false) + "%"
			return fmt.Sprintf(`%s %s like %s`, operation, name, m.addValue(name, value))
		}
	case "likeend", "like_end": // 以xx结束
		makeSqlStr = func() string {
			value = "%" + anyToSqlString(value, false)
			return fmt.Sprintf(`%s %s like %s`, operation, name, m.addValue(name, value))
		}
	default:
		panic(fmt.Errorf(`syntax error, non-supported flag "%s"`, flag))
	}

	if !has {
		return ""
	}
	if value == nil {
		return fmt.Sprintf(`%s %s is null`, operation, name)
	}
	return makeSqlStr()
}

// sql模板解析, 同 SqlTemplateRender, 但是它返回的是orm使用的sql语句和参数
func SqlTemplateParse(sql_template string, kvs ...interface{}) (sql_str string, names []string, args []interface{}) {
	return newSqlTemplate(kvs...).Parse(sql_template)
}

// sql模板渲染, 注意, 这个函数不支持sql注入检查
//
// 语法格式1:   (操作符)(name)   示例:   &a   |a
// 语法格式2:   {(操作符)(name)}   示例:   {&a}   {|a}
// 语法格式3:   {(操作符)(name) (对比标志)}   示例:   {&a in}   {|a >}
//
// 操作符支持:
//     @: 直接赋值, 如果没有传值存在外壳转为null无外壳会panic, 这个操作符仅支持以下格式
//          (操作符)(name)
//          {(操作符)(name)}
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
			return anyToSqlString(value, true)
		}
		if crust {
			return "null"
		}
		panic(fmt.Sprintf(`"%s" must have a value`, text))
	case "&":
		operation = "and"
	case "|":
		operation = "or"
	default:
		panic(fmt.Errorf(`syntax error, non-supported operation "%s"`, operation))
	}

	var sql_str string
	switch cflag {
	case ">", ">=", "<", "<=", "!=", "<>", "=":
		sql_str = fmt.Sprintf(`%s %s %s %s`, operation, name, cflag, anyToSqlString(value, true))
	case "in", "not in":
		sql_str = fmt.Sprintf(`%s %s %s %s`, operation, name, cflag, anyToSqlString(value, true))
	case "like": // 包含xx
		sql_str = fmt.Sprintf(`%s %s like "%%%s%%"`, operation, name, anyToSqlString(value, false))
	case "likestart", "like_start": // 以xx开始
		sql_str = fmt.Sprintf(`%s %s like "%s%%"`, operation, name, anyToSqlString(value, false))
	case "likeend", "like_end": // 以xx结束
		sql_str = fmt.Sprintf(`%s %s like "%%%s"`, operation, name, anyToSqlString(value, false))
	default:
		panic(fmt.Errorf(`syntax error, non-supported flag "%s"`, flag))
	}

	if !has {
		return ""
	}
	if value == nil {
		return fmt.Sprintf(`%s %s is null`, operation, name)
	}
	return sql_str
}
