
字符串转换工具, 类型转换工具, 模板渲染, 动态sql

---

<!-- TOC -->

- [转换函数](#%E8%BD%AC%E6%8D%A2%E5%87%BD%E6%95%B0)
- [扫描函数](#%E6%89%AB%E6%8F%8F%E5%87%BD%E6%95%B0)
- [字符串工具](#%E5%AD%97%E7%AC%A6%E4%B8%B2%E5%B7%A5%E5%85%B7)
- [模板渲染](#%E6%A8%A1%E6%9D%BF%E6%B8%B2%E6%9F%93)
    - [示例](#%E7%A4%BA%E4%BE%8B)
    - [变量名定义](#%E5%8F%98%E9%87%8F%E5%90%8D%E5%AE%9A%E4%B9%89)
        - [渲染](#%E6%B8%B2%E6%9F%93)
- [值处理](#%E5%80%BC%E5%A4%84%E7%90%86)
- [变量](#%E5%8F%98%E9%87%8F)
    - [模板匹配变量名方式](#%E6%A8%A1%E6%9D%BF%E5%8C%B9%E9%85%8D%E5%8F%98%E9%87%8F%E5%90%8D%E6%96%B9%E5%BC%8F)
    - [变量名匹配优先级](#%E5%8F%98%E9%87%8F%E5%90%8D%E5%8C%B9%E9%85%8D%E4%BC%98%E5%85%88%E7%BA%A7)
- [动态sql](#%E5%8A%A8%E6%80%81sql)
    - [示例](#%E7%A4%BA%E4%BE%8B)
    - [模板语法说明](#%E6%A8%A1%E6%9D%BF%E8%AF%AD%E6%B3%95%E8%AF%B4%E6%98%8E)
        - [操作符](#%E6%93%8D%E4%BD%9C%E7%AC%A6)
        - [变量名](#%E5%8F%98%E9%87%8F%E5%90%8D)
        - [标记](#%E6%A0%87%E8%AE%B0)
        - [选项](#%E9%80%89%E9%A1%B9)
        - [渲染](#%E6%B8%B2%E6%9F%93)
- [模板渲染和动态sql性能说明](#%E6%A8%A1%E6%9D%BF%E6%B8%B2%E6%9F%93%E5%92%8C%E5%8A%A8%E6%80%81sql%E6%80%A7%E8%83%BD%E8%AF%B4%E6%98%8E)

<!-- /TOC -->

---

# 转换函数

支持 `string`, `[]byte`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`.
转换失败会返回该类型的零值.

```go
zstr.GetString(1)       // 1
zstr.GetBool(1)         // true
zstr.GetInt(true)       // 1
zstr.GetFloat32("3.2")  // 3.2
zstr.Get...
```

转换失败时使用默认值

```go
zstr.GetBool("xxx", true)         // true
zstr.GetInt("xxx", 1)       // 1
zstr.GetFloat32("xxx", float32(3.2))  // 3.2
zstr.Get...
```

Boolean转换说明

```text
# 能转为true的数据
1, t, T, true, TRUE, True, y, Y, yes, YES, Yes, on, ON, On, ok, OK, Ok, enabled, ENABLED, Enabled
# 能转为false数据
nil, 0, f, F, false, FALSE, False, n, N, no, NO, No, off, OFF, Off, disable, DISABLE, Disable
```

---

# 扫描函数

支持 `string`, `[]byte`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`.

将字符串扫描到指定变量

```
var a float64
var b uint32
var c bool
var d string
zstr.Scan("1", &a)  // 1
zstr.Scan("1", &b)  // 1
zstr.Scan("1", &c)  // true
zstr.Scan("1", &d)  // 1
```

将任何值扫描到指定变量

```
var a float64
var b uint32
var c bool
var d string
zstr.ScanAny(true, &a)     // 1
zstr.ScanAny(false, &b)    // 0
zstr.ScanAny("ok", &c)     // true
zstr.ScanAny(1.23, &d)     // 1.23
```

---

# 字符串工具

```go
// 转为zstr.String类型
s := zstr.String("1")
```

转为指定类型

转换失败会返回该类型的零值.
支持 `string`, `[]byte`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`

```go
s.String()      // 获取string
s.GetBool()     // 获取bool
s.GetInt()      // 获取int
s.GetInt32()    // 获取int32
s.GetUint64()   // 获取uint64
s.GetFloat32()  // 获取float32
s.Get...
```

转换失败时使用默认值

```go
s := zstr.String("xxx")
s.GetBool(true)     // true
s.GetInt(1)      // 1
s.GetInt32(2)    // 2
s.GetUint64(3)   // 3
s.GetFloat32(4)  // 4
s.Get...
```

扫描到变量

支持 `string`, `[]byte`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`

```go
var a float64
var b uint32
var c bool
var d string
_ = s.Scan(&a)
_ = s.Scan(&b)
_ = s.Scan(&c)
_ = s.Scan(&d)
```

---

# 模板渲染

将模板数据中的变量替换为指定数据

```go
func Render(format string, values ...interface{}) string
```

`format` 表示模板数据. 示例: `@a @b {@c}`
`values` 表示值. 参考[值处理](#值处理)


## 示例

```go
zstr.Render("{@a}e", "va") // 按顺序赋值
zstr.Render("@a @a e", "va0", "va1") // 按顺序赋值
zstr.Render("@a e", map[string]string{"a": "va"}) // 指定变量名赋值
zstr.Render("@a @b @c", zstr.KV{"a", "aValue"}, zstr.KV{"b", "bValue"}, zstr.KV{"c", "cValue"}) // 指定变量名赋值
zstr.Render("@a @b @c", zstr.KVs{{"a", "aValue"}, {"b", "bValue"}, {"c", "cValue"}}) // 指定变量名赋值
zstr.Render("@a @a @a", zstr.KV{"a[0]", "1"}, zstr.KV{"a", "2"}) // 指定下标, 指定变量名+下标的优先级比指定变量名更高
```

## 变量名定义


模板变量用 `@` 开头
模板变量推荐使用花括号`{}`包起来, 这样能精确的界定模板变量的开头和结尾. 注意花括号内不能有空格.
示例:  `@a`,  `@a_b`,  `@A.c`,  `{@a....c}`,  `{@9}`,  `@0...A`


变量细节参考 [变量](#变量)

### 渲染

模板渲染时遍历找出所有模板变量, 然后替换为匹配的变量值.

如果未对模板中的变量进行赋值并且该变量被花括号`{}`包裹, 那么该会被替换为空字符串.

比如 `zstr.Render("@a @b {@c}", zstr.KVs{{"a", "aValue"}, {"b", "bValue"}})` 未对变量 `c` 进行赋值, 其输出结果为 `aValue bValue `.
但是 `zstr.Render("@a @b @c", zstr.KVs{{"a", "aValue"}, {"b", "bValue"}})` 输出结果为 `aValue bValue @c`, 因为变量 `c` 没有被花括号`{}`包裹

---

# 值处理

> 参考源码 [MakeMapOfValues](./values_to_map.go)


支持任何类型的map, 如 `map[string]interface{}`; `map[int]int`; `map[int]string`; `map[string]int` 等.

可以通过 `zstr.KV` 方式传入. 示例: `zstr.MakeMapOfValues(zstr.KV{"k1", "v1"}, &zstr.KV{"k2", "v2"})`.

可以通过 `zstr.KVs` 方式传入. 示例: `zstr.MakeMapOfValues(zstr.KVs{{"k1", "v1"}, {"k2", "v2"}})`.

也可以 `zstr.KV` 和 `zstr.KVs` 混合传入. 示例  `zstr.MakeMapOfValues(zstr.KVs{{"k1", "v1"}}, zstr.KV{"k2", "v2"})`.

传参时支持顺序传值, 示例 `zstr.MakeMapOfValues("v1", "v2")`, 结果为 `map[string]interface{}{"*[0]": "v1", "*[1]": "v2"}`

---

# 变量

变量名支持 大小写字母; 数字; 下划线; 小数点;
变量名区分大小写;
变量开头和结尾不能是小数点, 但是变量中间可以有任意数量的小数点;

示例:  `@a`,  `@a_b`,  `@A.c`,  `{@a....c}`,  `{@9}`,  `@0...A`

## 模板匹配变量名方式

`key` 可以带下标, 下标从 0 开始, 如: `a[0]` 表示第一次出现的 `a`, `a[5]`表示第 6 次出现的 `a`.
星号 `*` 可以匹配任何变量名, 但是必须和下标一起使用, 下标从 0 开始, 如: `*[0]` 表示第 1 个变量, `*[5]` 表示第 6 个变量.

## 变量名匹配优先级

1.  完全匹配变量名+下标
2.  完全匹配变量名
3.  星号+下标

如:  `a[0]` > `a` > `*[0]`


---

# 动态sql

你是否会有根据不同的条件拼接不同的sql语句的痛苦, 使用 [zstr动态sql](./sql_template.go) 一次根治

```go
func SqlRender(sqlTemplate string, values ...interface{}) string
```

`sqlTemplate` 表示要进行渲染的sql模板数据. 示例: `select * from table where &a (&b |c) {&d like_start}`
`values` 表示值. 参考[值处理](#值处理)

## 示例

```go
zstr.SqlRender("select * from table where &a (&b |c)", map[string]interface{}{
    "a": 1,
    "b": 2,
    "c": 3,
})
```

## 模板语法说明

语法格式

`(操作符)(变量名)`
`{(操作符)(变量名)}`
`{(操作符)(变量名) (标记)}`
`{(操作符)(变量名) (标记) (选项)}`
`{(操作符)(变量名) (选项)}`

示例

`&a`
`{&a}`
`{&a like}`
`{&a like d}`
`{&a d}`

### 操作符

`&` 转为 `and 变量名 标记 值`, 一般用于表示一个条件.

`|` 转为 `or 变量名 标记 值`, 一般用于表示一个条件.

`#` 转为 `值`. 自带 `attention` 选项.

`@` 转为 `值`, 一般用于表示一条语句. 自带 `direct` 选项, 且不会为字符串加上引号. `attention` 选项无效

### 变量名

模板变量用一个 [操作符](#操作符) 开头.
模板变量推荐使用花括号`{}`包起来, 这样能精确的界定模板变量的开头和结尾, 还可以使用 [标记](#标记) 和 [选项](#选项) 来改变行为. 花括号内允许有空格.
示例:  `&a`,  `|a_b`,  `#A.c`,  `{&a....c}`,  `{@9}`,  `@0...A`

变量细节参考 [变量](#变量)

### 标记

用于改变操作符 `&` 和 `|` 的行为, 默认标记为 =

`>`
`>=`
`<`
`<=`
`!=` 和 `<>`
`=`
`in`
`not_in` 和 `notin`
`like`
`like_start` 和 `likestart`
`like_end` 和 `likeend`

### 选项

+ attention, 语法 a, 不会忽略参数值为该类型的零值. 表示该条件必须存在
+ direct, 语法 d, 直接将值写入sql语句中. 一般用于表示一段语句, 建议别用于表示一个值, 否则可能导致sql注入危险.
+ must, 语法 m, 必须传值, 否则会panic

### 渲染

**渲染文本之前会缩进文本所有的空格, 所以要渲染的文本一般为单纯的sql语句, 不要将值写在文本中, 而是使用传参的方式, 这是使用sql的标准姿势**

模板渲染时遍历找出所有模板语法, 然后按语法替换为不同的值.
一般情况下如果变量没有传参或为该类型的零值, 则替换为空字符串.
如果变量的值为nil, 不同的标志会转为不同的语句.

---

# 模板渲染和动态sql性能说明

我们写了一个专用函数用于替换正则查找规则, 经过 Benchmark 测试, 速度为正则查找变量方式的 230%.
具体的性能可以将代码clone后执行命令 `go test -v -bench .`
