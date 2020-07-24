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

var (
	// 标准
	sqlTemplateRegex = regexp.MustCompile(`[&|@]\w*\.?\w+`)
	// 加壳
	sqlTemplateRegexCrust = regexp.MustCompile(`\{[\s\S]*?\}`)

	// id
	sqlTemplateParseIdRegex = regexp.MustCompile(`\{\{\d+\}\}`)

	// 空字符串
	emptyStrRegex = regexp.MustCompile(`\s+`)

	// 变量名
	variableNameRegex = regexp.MustCompile(`^\w*\.?\w+$`)
	// 操作符
	sqlTemplateOperationMapp = map[string]struct{}{
		"@": {},
		"&": {},
		"|": {},
	}
	// 标记
	sqlTemplateFlagMapp = map[string]struct{}{
		">":          {},
		">=":         {},
		"<":          {},
		"<=":         {},
		"!=":         {},
		"<>":         {},
		"=":          {},
		"in":         {},
		"notin":      {},
		"not_in":     {},
		"like":       {},
		"likestart":  {},
		"like_start": {},
		"likeend":    {},
		"like_end":   {},
	}
	// 选项
	sqlTemplateOptsMapp = map[int32]struct{}{
		'i': {},
		'd': {},
	}
)

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
		operation, name, flag, opts, err := sqlTemplateSyntaxParse(s[1 : len(s)-1])
		if err != nil {
			panic(err)
		}
		return m.translate(operation, name, flag, opts, true)
	})

	sql_str = sqlTemplateRegex.ReplaceAllStringFunc(sql_str, func(s string) string {
		operation, name, flag, opts, err := sqlTemplateSyntaxParse(s)
		if err != nil {
			panic(err)
		}
		return m.translate(operation, name, flag, opts, false)
	})

	// 按顺序写入names和args
	sql_str = sqlTemplateParseIdRegex.ReplaceAllStringFunc(sql_str, func(s string) string {
		index, _ := strconv.Atoi(s[2 : len(s)-2])
		names = append(names, m.names[index])
		args = append(args, m.values[index])
		return "?"
	})

	sql_str = repairSql(sql_str)
	return sql_str, names, args
}

func (m *sqlTemplate) translate(operation, name, flag string, opts string, crust bool) string {
	// 选项检查
	var ignore_opt, direct_opt bool
	for _, o := range opts {
		switch o {
		case 'i':
			ignore_opt = true
		case 'd':
			direct_opt = true
		default:
			panic(fmt.Sprintf(`syntax error, non-supported option "%s"`, string(o)))
		}
	}

	value, has := m.data[name]

	// 操作检查
	switch operation {
	case "@":
		if has {
			return m.addValue(name, value)
		}
		if crust {
			return "null"
		}
		panic(fmt.Sprintf(`"%s" must have a value`, name))
	case "&":
		operation = "and"
	case "|":
		operation = "or"
	default:
		panic(fmt.Errorf(`syntax error, non-supported operation "%s"`, operation))
	}

	var makeSqlStr func() string
	var directWrite func() string
	// 标记
	switch flag {
	case ">", ">=", "<", "<=", "!=", "<>", "=":
		makeSqlStr = func() string {
			return fmt.Sprintf(`%s %s %s %s`, operation, name, flag, m.addValue(name, value))
		}
		directWrite = func() string {
			return fmt.Sprintf(`%s %s %s %s`, operation, name, flag, anyToSqlString(value, true))
		}
	case "in":
		makeSqlStr = func() string {
			return fmt.Sprintf(`%s %s %s (%s)`, operation, name, flag, m.addValue(name, value))
		}
		directWrite = func() string {
			return fmt.Sprintf(`%s %s %s %s`, operation, name, flag, anyToSqlString(value, true))
		}
	case "notin", "not_in":
		makeSqlStr = func() string {
			return fmt.Sprintf(`%s %s not in (%s)`, operation, name, m.addValue(name, value))
		}
		directWrite = func() string {
			return fmt.Sprintf(`%s %s not in %s`, operation, name, anyToSqlString(value, true))
		}
	case "like": // 包含xx
		makeSqlStr = func() string {
			value = "%" + anyToSqlString(value, false) + "%"
			return fmt.Sprintf(`%s %s like %s`, operation, name, m.addValue(name, value))
		}
		directWrite = func() string {
			return fmt.Sprintf(`%s %s like "%%%s%%"`, operation, name, anyToSqlString(value, false))
		}
	case "likestart", "like_start": // 以xx开始
		makeSqlStr = func() string {
			value = anyToSqlString(value, false) + "%"
			return fmt.Sprintf(`%s %s like %s`, operation, name, m.addValue(name, value))
		}
		directWrite = func() string {
			return fmt.Sprintf(`%s %s like "%s%%"`, operation, name, anyToSqlString(value, false))
		}
	case "likeend", "like_end": // 以xx结束
		makeSqlStr = func() string {
			value = "%" + anyToSqlString(value, false)
			return fmt.Sprintf(`%s %s like %s`, operation, name, m.addValue(name, value))
		}
		directWrite = func() string {
			return fmt.Sprintf(`%s %s like "%%%s"`, operation, name, anyToSqlString(value, false))
		}
	default:
		panic(fmt.Errorf(`syntax error, non-supported flag "%s"`, flag))
	}

	// 无值返回空sql语句
	if !has {
		return ""
	}

	// 忽略模式, 零值返回空sql语句
	if ignore_opt && IsZero(value) {
		return ""
	}

	// nil 改为 is null
	if value == nil {
		return fmt.Sprintf(`%s %s is null`, operation, name)
	}

	// 直接模式, 将值写入sql语句
	if direct_opt {
		return directWrite()
	}
	return makeSqlStr()
}

func sqlTemplateSyntaxParse(text string) (operation, name, flag, opts string, err error) {
	// 去头去尾
	temp := strings.TrimSpace(text)
	// 空数据
	if temp == "" {
		err = fmt.Errorf("syntax error, {%s}, empty data", text)
		return
	}

	// 分离操作符
	operation, temp = temp[:1], temp[1:]

	// 缩进空格
	temp = emptyStrRegex.ReplaceAllString(temp, " ")

	// 分离数据
	texts := strings.SplitN(temp, " ", 4) // 4为考虑尾部可能有空格
	if len(texts) >= 1 {
		name = texts[0]
	}
	if len(texts) >= 2 {
		flag = texts[1]
	} else {
		flag = defaultSqlCompareFlag
	}
	if len(texts) >= 3 {
		opts = texts[2]
	}
	if len(texts) >= 4 && texts[3] != " " {
		err = fmt.Errorf("syntax error, {%s}, redundant data", text)
		return
	}

	// 检查操作符
	if _, ok := sqlTemplateOperationMapp[operation]; !ok {
		err = fmt.Errorf(`syntax error, {%s}, non-supported operation "%s"`, text, operation)
		return
	}

	// 检查变量名
	if name == "" {
		err = fmt.Errorf("syntax error, {%s}, no variable name", text)
		return
	}
	if !variableNameRegex.MatchString(name) {
		err = fmt.Errorf("syntax error, {%s}, Invalid variable name", text)
		return
	}

	// 检查标记
	if _, ok := sqlTemplateFlagMapp[flag]; !ok {
		err = fmt.Errorf(`syntax error, {%s}, non-supported flag "%s"`, text, flag)
		return
	}

	// 检查选项
	os := make(map[int32]struct{})
	for _, o := range opts {
		if _, ok := sqlTemplateOptsMapp[o]; !ok {
			err = fmt.Errorf(`syntax error, {%s}, non-supported option "%s"`, text, string(o))
			return
		}
		// 重复选项
		if _, ok := os[o]; ok {
			err = fmt.Errorf(`syntax error, {%s}, repetitive option "%s"`, text, string(o))
			return
		}
		os[o] = struct{}{}
	}

	return
}

// sql模板解析
//
// 语法格式1:   (操作符)(name)   示例:   &a   |a
// 语法格式2:   {(操作符)(name)}   示例:   {&a}   {|a}
// 语法格式3:   {(操作符)(name) (对比标志)}   示例:   {&a in}   {|a >}
// 语法格式4:   {(操作符)(name) (对比标志) (选项)}   示例:   {&a in i}   {|a > id}
//     选项:
//          i:   ignore, 如果参数值为该类型的零值则忽略
//          d:   direct, 直接将值写入sql语句中
//
// 操作符支持:
//     @: 直接赋值, 如果没有传值存在外壳转为null无外壳会panic, 这个操作符仅支持以下格式
//          语法格式1:   (操作符)(name)
//          语法格式2:   {(操作符)(name)}
//     &: 转为 and
//     |: 转为 or
//
// 对比标志支持:   >   >=   <   <=   !=   <>   =   in   notin   like   likestart    like_start   likeend   like_end
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
//          语法格式1:   (操作符)(name)
//          语法格式2:   {(操作符)(name)}
//     &: 转为 and
//     |: 转为 or
//
// 对比标志支持:   >   >=   <   <=   !=   <>   =   in   notin   like   likestart    like_start   likeend   like_end
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
		operation, name, flag, opts, err := sqlTemplateSyntaxParse(s[1 : len(s)-1])
		if err != nil {
			panic(err)
		}
		return sqlTranslate(operation, name, flag, opts, true, data)
	})
	result = sqlTemplateRegex.ReplaceAllStringFunc(result, func(s string) string {
		operation, name, flag, opts, err := sqlTemplateSyntaxParse(s)
		if err != nil {
			panic(err)
		}
		return sqlTranslate(operation, name, flag, opts, false, data)
	})
	return repairSql(result)
}

func sqlTranslate(operation, name, flag string, opts string, crust bool, m map[string]interface{}) string {
	// 选项检查
	var ignore_opt bool
	for _, o := range opts {
		switch o {
		case 'i':
			ignore_opt = true
		default:
			panic(fmt.Sprintf(`syntax error, non-supported option "%s"`, string(o)))
		}
	}

	value, has := m[name]

	switch operation {
	case "@":
		if has {
			return anyToSqlString(value, true)
		}
		if crust {
			return "null"
		}
		panic(fmt.Sprintf(`"%s" must have a value`, name))
	case "&":
		operation = "and"
	case "|":
		operation = "or"
	default:
		panic(fmt.Errorf(`syntax error, non-supported operation "%s"`, operation))
	}

	var sql_str string
	switch flag {
	case ">", ">=", "<", "<=", "!=", "<>", "=":
		sql_str = fmt.Sprintf(`%s %s %s %s`, operation, name, flag, anyToSqlString(value, true))
	case "in":
		sql_str = fmt.Sprintf(`%s %s %s %s`, operation, name, flag, anyToSqlString(value, true))
	case "notin", "not_in":
		sql_str = fmt.Sprintf(`%s %s not in %s`, operation, name, anyToSqlString(value, true))
	case "like": // 包含xx
		sql_str = fmt.Sprintf(`%s %s like "%%%s%%"`, operation, name, anyToSqlString(value, false))
	case "likestart", "like_start": // 以xx开始
		sql_str = fmt.Sprintf(`%s %s like "%s%%"`, operation, name, anyToSqlString(value, false))
	case "likeend", "like_end": // 以xx结束
		sql_str = fmt.Sprintf(`%s %s like "%%%s"`, operation, name, anyToSqlString(value, false))
	default:
		panic(fmt.Errorf(`syntax error, non-supported flag "%s"`, flag))
	}

	// 无值返回空sql语句
	if !has {
		return ""
	}

	// 忽略模式, 零值返回空sql语句
	if ignore_opt && IsZero(value) {
		return ""
	}

	// nil 改为 is null
	if value == nil {
		return fmt.Sprintf(`%s %s is null`, operation, name)
	}

	return sql_str
}
