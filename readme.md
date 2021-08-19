
# 字符串工具, 模板渲染, 动态sql

---

<!-- TOC -->

- [字符串工具](#%E5%AD%97%E7%AC%A6%E4%B8%B2%E5%B7%A5%E5%85%B7)
    - [示例](#%E7%A4%BA%E4%BE%8B)
    - [快捷函数](#%E5%BF%AB%E6%8D%B7%E5%87%BD%E6%95%B0)
    - [其他方法](#%E5%85%B6%E4%BB%96%E6%96%B9%E6%B3%95)
    - [Boolean支持](#boolean%E6%94%AF%E6%8C%81)
- [模板渲染](#%E6%A8%A1%E6%9D%BF%E6%B8%B2%E6%9F%93)
    - [示例](#%E7%A4%BA%E4%BE%8B)
    - [变量名](#%E5%8F%98%E9%87%8F%E5%90%8D)
    - [模板渲染说明](#%E6%A8%A1%E6%9D%BF%E6%B8%B2%E6%9F%93%E8%AF%B4%E6%98%8E)
        - [变量名匹配](#%E5%8F%98%E9%87%8F%E5%90%8D%E5%8C%B9%E9%85%8D)
        - [变量名匹配优先级](#%E5%8F%98%E9%87%8F%E5%90%8D%E5%8C%B9%E9%85%8D%E4%BC%98%E5%85%88%E7%BA%A7)
        - [值](#%E5%80%BC)
        - [渲染](#%E6%B8%B2%E6%9F%93)
- [动态sql](#%E5%8A%A8%E6%80%81sql)
    - [示例](#%E7%A4%BA%E4%BE%8B)
    - [模板语法说明](#%E6%A8%A1%E6%9D%BF%E8%AF%AD%E6%B3%95%E8%AF%B4%E6%98%8E)
        - [操作符](#%E6%93%8D%E4%BD%9C%E7%AC%A6)
        - [变量名](#%E5%8F%98%E9%87%8F%E5%90%8D)
        - [标记](#%E6%A0%87%E8%AE%B0)
        - [选项](#%E9%80%89%E9%A1%B9)
    - [模板渲染说明](#%E6%A8%A1%E6%9D%BF%E6%B8%B2%E6%9F%93%E8%AF%B4%E6%98%8E)
        - [变量名匹配](#%E5%8F%98%E9%87%8F%E5%90%8D%E5%8C%B9%E9%85%8D)
        - [变量名匹配优先级](#%E5%8F%98%E9%87%8F%E5%90%8D%E5%8C%B9%E9%85%8D%E4%BC%98%E5%85%88%E7%BA%A7)
        - [值](#%E5%80%BC)
        - [渲染](#%E6%B8%B2%E6%9F%93)
- [模板渲染和动态sql性能说明](#%E6%A8%A1%E6%9D%BF%E6%B8%B2%E6%9F%93%E5%92%8C%E5%8A%A8%E6%80%81sql%E6%80%A7%E8%83%BD%E8%AF%B4%E6%98%8E)

<!-- /TOC -->

---

# 字符串工具

## 示例

```go
s := zstr.String("1")
s.Val()     // 获取val
s.GetBool() // 获取bool
s.GetInt()  // 获取int
var a float64
_ = s.Scan(&a) // 扫描到a中
```

## 快捷函数

```go
zstr.GetString("1")
zstr.GetBool("1")
zstr.GetInt("1")
var a float64
zstr.Scan("1", &a)
```

## 其他方法
> GetXXX方法不仅仅支持string, 它支持传入任何类型

```text
GetInt  GetInt8  GetInt16  GetInt32  GetInt64
GetUint  GetUint8  GetUint16  GetUint32  GetUint64
GetFloat32  GetFloat64
GetString
```

## Boolean支持

```text
# 能转为true的数据
1, t, T, true, TRUE, True, y, Y, yes, YES, Yes, on, ON, On, ok, OK, Ok, enabled, ENABLED, Enabled
# 能转为false数据
nil, 0, f, F, false, FALSE, False, n, N, no, NO, No, off, OFF, Off, disable, DISABLE, Disable
```

# 模板渲染

## 示例

```go
Render("s@a e", map[string]string{"a": "va"})
Render("s{@a}e", map[string]string{"a": "va"})
Render("s{@a}e", "va")
Render("s@a @a e", "va0", "va1")
```

## 变量名

```text
模板变量用 @ 开头, 变量名支持 大小写字母; 数字; 下划线; 小数点;
模板变量可以使用花括号`{}`包起来, 这样能精确的界定模板变量的开头和结尾. 注意花括号内不能有空格.
示例:
    @a
    @a_b
    @A.c
    {@a....c}
    {@9}
    @0...A
```

## 模板渲染说明

```go
func Render(format string, values ...interface{}) string
// format表示要进行渲染的文本
// values表示变量值
```

### 变量名匹配

```text
变量名带下标, 下标从0开始, 如: a[0]表示第一次出现的a, a[5]表示第6次出现的a.
和变量名相等, 如: a;  a_b;  A.c;  a....c;  9;  0...A
星号*可以匹配任何变量名, 但是必须和下标一起使用, 下标从0开始, 如: *[0]表示1个变量, *[5]表示第6个变量.
```

### 变量名匹配优先级

1.  变量名带下标
2.  变量名
3.  星号带下标

### 值
> 参考 [MakeMapOfValues](./values_to_map.go)

```text
支持任何类型的map, 如 map[string]interface{}; map[int]int, map[string]int 等.
传参时支持顺序传值, 如 Render("要渲染的文本", 值0, 值1, 值2), 它会转为以下描述
    map[string]interface{}{"*[0]": 值0, "*[1]", 值1, "*[2]", 值2}
```

### 渲染

```text
模板渲染时遍历找出所有模板变量, 然后替换为匹配的变量值.
如果模板变量未赋值则不会替换, 但是如果模板变量是被花括号`{}`包起来的, 会替换为空字符串.
```

# 动态sql

> 你是否会有根据不同的条件拼接不同的sql语句的痛苦, 使用 [zstr动态sql](./sql_template.go) 一次根治

## 示例

```go
zstr.SqlRender("select * from table where &a (&b |c)", map[string]interface{}{
    "a": 1,
    "b": 2,
    "c": 3,
})
```

## 模板语法说明

```text
语法格式如下:
    (操作符)(变量名)
    {(操作符)(变量名)}
    {(操作符)(变量名) (标记)}
    {(操作符)(变量名) (标记) (选项)}
    {(操作符)(变量名) (选项)}
示例:
    &a
    {&a}
    {&a like}
    {&a like d}
    {&a d}
```

### 操作符

+ & 转为 `and 变量名 标记 值`

+ | 转为 `or 变量名 标记 值`

+ \# 转为 `值`

    ```text
    自带 attention 选项
    ```

+ @ 转为 `值`, 一般用于写入一条语句

    ```text
    attention 选项无效
    自带 direct 选项, 不会为字符串加上引号
    ```

### 变量名

```text
模板变量用一个 操作符 开头. 变量名支持: 下划线; 大小写字母; 数字; 小数点
模板变量可以使用花括号`{}`包起来, 这样能精确的界定模板变量的开头和结尾. 花括号内允许有空格
示例:
    &a
    |a_b
    #A.c
    {&a....c}
    {@9}
    @0...A
```

### 标记

> 默认标记为 =

```text
>
>=
<
<=
!=  <>
=
in
not_in  notin
like
like_start  likestart
like_end  likeend
```

### 选项

+ attention, 语法 a, 不会忽略参数值为该类型的零值
+ direct, 语法 d, 直接将值写入sql语句中
+ must, 语法 m, 必须传值

## 模板渲染说明

```go
func SqlRender(sql_template string, values ...interface{}) string
// sql_template表示要进行渲染的sql
// values表示变量值
```

### 变量名匹配

```text
变量名带下标, 下标从0开始, 如: a[0]表示第一次出现的a, a[5]表示第6次出现的a.
和变量名相等, 如: a;  a_b;  A.c;  a....c;  9;  0...A
星号*可以匹配任何变量名, 但是必须和下标一起使用, 下标从0开始, 如: *[0]表示1个变量, *[5]表示第6个变量.
```

### 变量名匹配优先级

1.  变量名带下标
2.  变量名
3.  星号带下标

### 值
> 参考 [MakeMapOfValues](./values_to_map.go)

```text
支持任何类型的map, 如 map[string]interface{}; map[int]int, map[string]int 等.
传参时支持顺序传值, 如 Render("要渲染的文本", 值0, 值1, 值2), 它会转为以下描述
    map[string]interface{}{"*[0]": 值0, "*[1]", 值1, "*[2]", 值2}
```

### 渲染

> 渲染文本之前会缩进文本所有的空格, 所以要渲染的文本一般为单纯的sql语句, 不要将值写在文本中, 而是使用传参的概念, 这是使用sql的标准姿势

```text
模板渲染时遍历找出所有模板语法, 然后按语法替换为不同的值.
一般情况下如果变量没有传参或为该类型的零值, 则替换为空字符串.
如果变量的值为nil, 不同的标志会转为不同的语句.
```

# 模板渲染和动态sql性能说明

```text
我们写了一个专用函数用于替换正则查找规则, 经过Benchmark测试, 速度为正则查找变量方式的230%!!!
具体的性能可以将代码clone后执行以下命令
    go test -v -bench .
```
