/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/11
   Description :
-------------------------------------------------
*/

package zstr

import (
    "fmt"
    "strconv"
)

type String struct {
    val string
}

func New(s string) *String {
    return &String{val: s}
}

func (m *String) String() string {
    return m.val
}
func (m *String) Val() string {
    return m.val
}
func (m *String) Bytes() []byte {
    return []byte(m.val)
}
func (m *String) Bool() (bool, error) {
    switch m.val {
    case "1", "t", "T", "true", "TRUE", "True", "y", "Y", "yes", "YES", "Yes", "on", "ON", "On":
        return true, nil
    case "0", "f", "F", "false", "FALSE", "False", "n", "N", "no", "NO", "No", "off", "OFF", "Off":
        return false, nil
    }
    return false, fmt.Errorf("数据\"%s\"无法转换为bool", m.val)
}

func (m *String) Int() (int, error) {
    return strconv.Atoi(m.val)
}
func (m *String) Int8() (int8, error) {
    n, err := strconv.ParseInt(m.val, 10, 8)
    if err != nil {
        return 0, err
    }
    return int8(n), nil
}
func (m *String) Int16() (int16, error) {
    n, err := strconv.ParseInt(m.val, 10, 16)
    if err != nil {
        return 0, err
    }
    return int16(n), nil
}
func (m *String) Int32() (int32, error) {
    n, err := strconv.ParseInt(m.val, 10, 32)
    if err != nil {
        return 0, err
    }
    return int32(n), nil
}
func (m *String) Int64() (int64, error) {
    return strconv.ParseInt(m.val, 10, 64)
}

func (m *String) Uint() (uint, error) {
    n, err := strconv.ParseUint(m.val, 10, 64)
    if err != nil {
        return 0, err
    }
    return uint(n), err
}
func (m *String) Uint8() (uint8, error) {
    n, err := strconv.ParseUint(m.val, 10, 8)
    if err != nil {
        return 0, err
    }
    return uint8(n), nil
}
func (m *String) Uint16() (uint16, error) {
    n, err := strconv.ParseUint(m.val, 10, 16)
    if err != nil {
        return 0, err
    }
    return uint16(n), nil
}
func (m *String) Uint32() (uint32, error) {
    n, err := strconv.ParseUint(m.val, 10, 32)
    if err != nil {
        return 0, err
    }
    return uint32(n), nil
}
func (m *String) Uint64() (uint64, error) {
    return strconv.ParseUint(m.val, 10, 64)
}

func (m *String) Float32() (float32, error) {
    f, err := strconv.ParseFloat(m.val, 32)
    if err != nil {
        return 0, err
    }
    return float32(f), nil
}
func (m *String) Float64() (float64, error) {
    return strconv.ParseFloat(m.val, 64)
}

func (m *String) Scan(v interface{}) error {
    return Scan(m.val, v)
}
