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
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const defaultSqlCompareFlag = "="

var (
	// 标准
	sqlTemplateRegex = regexp.MustCompile(`[&|#@]\w*\.?\w+`)
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
		"&": {},
		"|": {},
		"#": {},
		"@": {},
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
		'a': {}, // attention, 不会忽略参数值为该类型的零值
		'd': {}, // direct, 直接将值写入sql语句
		'm': {}, // must, 必填
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
		operation, name, flag, opts, err := m.sqlTemplateSyntaxParse(s[1 : len(s)-1])
		if err != nil {
			panic(err)
		}
		return m.translate(operation, name, flag, opts, true)
	})

	sql_str = sqlTemplateRegex.ReplaceAllStringFunc(sql_str, func(s string) string {
		operation, name, flag, opts, err := m.sqlTemplateSyntaxParse(s)
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
	var attention_opt, direct_opt, must_opt bool
	for _, o := range opts {
		switch o {
		case 'a':
			attention_opt = true
		case 'd':
			direct_opt = true
		case 'm':
			must_opt = true
		default:
			panic(fmt.Sprintf(`syntax error, non-supported option "%s"`, string(o)))
		}
	}
	if operation == "@" {
		attention_opt = false
		direct_opt = true
	}

	value, has := m.data[name]

	// 无值返回空sql语句
	if !has {
		if must_opt {
			panic(fmt.Sprintf(`"%s" must have a value`, name))
		}
		return ""
	}

	// 非注意模式且值为零值返回空sql语句
	if !attention_opt && IsZero(value) {
		return ""
	}

	// 操作检查
	switch operation {
	case "&":
		operation = "and"
	case "|":
		operation = "or"
	case "#":
		// nil改为null
		if value == nil {
			return "null"
		}
		if direct_opt {
			return anyToSqlString(value, true)
		}
		return m.addValue(name, value)
	case "@": // !attention_opt + direct
		return anyToSqlString(value, false)
	default:
		panic(fmt.Errorf(`syntax error, non-supported operation "%s"`, operation))
	}

	// nil 修改语句
	if value == nil {
		switch flag {
		case "!=", "<>", "notin", "not_in", ">", "<":
			return fmt.Sprintf(`%s %s is not null`, operation, name)
		case "=", "like", "likestart", "like_start", "likeend", "like_end":
			return fmt.Sprintf(`%s %s is null`, operation, name)
		case "in", ">=", "<=":
			return ""
		}
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
		values := m.parseToSlice(value)
		if len(values) == 0 {
			return ""
		}
		makeSqlStr = func() string {
			if len(values) == 1 {
				return fmt.Sprintf(`%s %s = %s`, operation, name, m.addValue(name, values[0]))
			}
			var fs []string
			for i, s := range values {
				fs = append(fs, m.addValue(fmt.Sprintf("%s.in(%d)", name, i), s))
			}
			return fmt.Sprintf(`%s %s in (%s)`, operation, name, strings.Join(fs, ","))
		}
		directWrite = func() string {
			if len(values) == 1 {
				return fmt.Sprintf(`%s %s = %s`, operation, name, anyToSqlString(values[0], true))
			}
			return fmt.Sprintf(`%s %s in %s`, operation, name, anyToSqlString(value, true))
		}
	case "notin", "not_in":
		values := m.parseToSlice(value)
		if len(values) == 0 {
			return ""
		}
		makeSqlStr = func() string {
			if len(values) == 1 {
				return fmt.Sprintf(`%s %s != %s`, operation, name, m.addValue(name, values[0]))
			}
			var fs []string
			for i, s := range values {
				fs = append(fs, m.addValue(fmt.Sprintf("%s.not_in(%d)", name, i), s))
			}
			return fmt.Sprintf(`%s %s not in (%s)`, operation, name, strings.Join(fs, ","))
		}
		directWrite = func() string {
			if len(values) == 1 {
				return fmt.Sprintf(`%s %s != %s`, operation, name, anyToSqlString(values[0], true))
			}
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

	// 直接模式, 将值写入sql语句
	if direct_opt {
		return directWrite()
	}
	return makeSqlStr()
}

func (m *sqlTemplate) Render(sql_template string) string {
	result := sqlTemplateRegexCrust.ReplaceAllStringFunc(sql_template, func(s string) string {
		operation, name, flag, opts, err := m.sqlTemplateSyntaxParse(s[1 : len(s)-1])
		if err != nil {
			panic(err)
		}
		return m.sqlTranslate(operation, name, flag, opts, true)
	})
	result = sqlTemplateRegex.ReplaceAllStringFunc(result, func(s string) string {
		operation, name, flag, opts, err := m.sqlTemplateSyntaxParse(s)
		if err != nil {
			panic(err)
		}
		return m.sqlTranslate(operation, name, flag, opts, false)
	})
	return repairSql(result)
}

func (m *sqlTemplate) sqlTranslate(operation, name, flag string, opts string, crust bool) string {
	// 选项检查
	var attention_opt, must_opt bool
	for _, o := range opts {
		switch o {
		case 'a':
			attention_opt = true
		case 'd':
		case 'm':
			must_opt = true
		default:
			panic(fmt.Sprintf(`syntax error, non-supported option "%s"`, string(o)))
		}
	}
	if operation == "@" {
		attention_opt = false
	}

	value, has := m.data[name]

	// 无值返回空sql语句
	if !has {
		if must_opt {
			panic(fmt.Sprintf(`"%s" must have a value`, name))
		}
		return ""
	}

	// 非注意模式, 零值返回空sql语句
	if !attention_opt && IsZero(value) {
		return ""
	}

	switch operation {
	case "&":
		operation = "and"
	case "|":
		operation = "or"
	case "#":
		// nil改为null
		if value == nil {
			return "null"
		}
		return anyToSqlString(value, true)
	case "@":
		return anyToSqlString(value, false)
	default:
		panic(fmt.Errorf(`syntax error, non-supported operation "%s"`, operation))
	}

	// nil 修改语句
	if value == nil {
		switch flag {
		case "!=", "<>", "notin", "not_in", ">", "<":
			return fmt.Sprintf(`%s %s is not null`, operation, name)
		case "=", "like", "likestart", "like_start", "likeend", "like_end":
			return fmt.Sprintf(`%s %s is null`, operation, name)
		case "in", ">=", "<=":
			return ""
		}
	}

	var sql_str string
	switch flag {
	case ">", ">=", "<", "<=", "!=", "<>", "=":
		sql_str = fmt.Sprintf(`%s %s %s %s`, operation, name, flag, anyToSqlString(value, true))
	case "in":
		values := m.parseToSlice(value)
		if len(values) == 0 {
			return ""
		}
		if len(values) == 1 {
			return fmt.Sprintf(`%s %s = %s`, operation, name, anyToSqlString(values[0], true))
		}
		sql_str = fmt.Sprintf(`%s %s in %s`, operation, name, anyToSqlString(value, true))
	case "notin", "not_in":
		values := m.parseToSlice(value)
		if len(values) == 0 {
			return ""
		}
		if len(values) == 1 {
			return fmt.Sprintf(`%s %s != %s`, operation, name, anyToSqlString(values[0], true))
		}
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

	return sql_str
}

// 将数据解析为切片
func (m *sqlTemplate) parseToSlice(a interface{}) []interface{} {
	switch v := a.(type) {

	case nil:
		return []interface{}{"null"}

	case string, []byte, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return []interface{}{v}
	}

	r_v := reflect.Indirect(reflect.ValueOf(a))
	if r_v.Kind() != reflect.Slice && r_v.Kind() != reflect.Array {
		return []interface{}{fmt.Sprint(a)}
	}

	l := r_v.Len()
	out := make([]interface{}, 0, l)
	for i := 0; i < l; i++ {
		v := reflect.Indirect(r_v.Index(i)).Interface()
		out = append(out, m.parseToSlice(v)...)
	}
	return out
}

// sql模板语法解析
//
// 语法格式:   (操作符)(name)
// 语法格式:   {(操作符)(name)}
// 语法格式:   {(操作符)(name) (标志)}
// 语法格式:   {(操作符)(name) (标志) (选项)}
// 语法格式:   {(操作符)(name) (选项)}
//
// 操作符:
//     &: 转为 and name flag value
//     |: 转为 or name flag value
//     #: 转为 value, 仅支持以下格式
//          (操作符)(name)
//          {(操作符)(name)}
//          {(操作符)(name) (选项)}
//     @: attention 选项无效且自带 direct 选项, 且不会为字符串加上引号, 仅支持以下格式, 一般用于写入一条语句
//          (操作符)(name)
//          {(操作符)(name)}
//          {(操作符)(name) (选项)}
//
//
// name:   示例:    a   a2   a_2   a_2.b   a_2.b_2
//
// 标志:   >   >=   <   <=   !=   <>   =   in   notin   not_in   like   likestart    like_start   likeend   like_end
//
// 选项:
//     a:   attention, 不会忽略参数值为该类型的零值
//     d:   direct, 直接将值写入sql语句中
//     m:   must, 必须传值, 值可以为零值
//
// 输入的kvs必须为：map[string]string, map[string]interface{}, 或健值对
//
// 注意:
//     如果name没有传参或为该类型的零值, 则替换为空字符串
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
func (m *sqlTemplate) sqlTemplateSyntaxParse(text string) (operation, name, flag, opts string, err error) {
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
		if opts != "" {
			err = fmt.Errorf(`syntax error, {%s}, non-supported flag "%s"`, text, flag)
			return
		}
		flag, opts = defaultSqlCompareFlag, flag
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
func SqlTemplateParse(sql_template string, kvs ...interface{}) (sql_str string, names []string, args []interface{}) {
	return newSqlTemplate(kvs...).Parse(sql_template)
}

// sql模板解析, 和SqlTemplateParse一样, 只是简短了函数名
func SqlParse(sql_template string, kvs ...interface{}) (sql_str string, names []string, args []interface{}) {
	return newSqlTemplate(kvs...).Parse(sql_template)
}

// sql模板渲染
//
// 值会直接写入sql语句中, 不支持sql注入检查
func SqlTemplateRender(sql_template string, kvs ...interface{}) string {
	return newSqlTemplate(kvs...).Render(sql_template)
}

// sql模板渲染, 和SqlTemplateRender一样, 只是简短了函数名
func SqlRender(sql_template string, kvs ...interface{}) string {
	return newSqlTemplate(kvs...).Render(sql_template)
}
