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
    case "1", "t", "T", "true", "TRUE", "True", "y", "Y", "yes", "YES", "Yes", "on", "ON", "On", "ok", "OK", "Ok":
        return true, nil
    case "0", "f", "F", "false", "FALSE", "False", "n", "N", "no", "NO", "No", "off", "OFF", "Off":
        return false, nil
    }
    return false, fmt.Errorf("数据\"%s\"无法转换为bool", m.val)
}
func (m *String) BoolDefault(def bool) bool {
    if a, err := m.Bool(); err == nil {
        return a
    }
    return def
}

func (m *String) Int() (int, error) {
    return strconv.Atoi(m.val)
}
func (m *String) IntDefault(def int) int {
    if a, err := m.Int(); err == nil {
        return a
    }
    return def
}
func (m *String) Int8() (int8, error) {
    n, err := strconv.ParseInt(m.val, 10, 8)
    if err != nil {
        return 0, err
    }
    return int8(n), nil
}
func (m *String) Int8Default(def int8) int8 {
    if a, err := m.Int8(); err == nil {
        return a
    }
    return def
}
func (m *String) Int16() (int16, error) {
    n, err := strconv.ParseInt(m.val, 10, 16)
    if err != nil {
        return 0, err
    }
    return int16(n), nil
}
func (m *String) Int16Default(def int16) int16 {
    if a, err := m.Int16(); err == nil {
        return a
    }
    return def
}
func (m *String) Int32() (int32, error) {
    n, err := strconv.ParseInt(m.val, 10, 32)
    if err != nil {
        return 0, err
    }
    return int32(n), nil
}
func (m *String) Int32Default(def int32) int32 {
    if a, err := m.Int32(); err == nil {
        return a
    }
    return def
}
func (m *String) Int64() (int64, error) {
    return strconv.ParseInt(m.val, 10, 64)
}
func (m *String) Int64Default(def int64) int64 {
    if a, err := m.Int64(); err == nil {
        return a
    }
    return def
}

func (m *String) Uint() (uint, error) {
    n, err := strconv.ParseUint(m.val, 10, 64)
    if err != nil {
        return 0, err
    }
    return uint(n), err
}
func (m *String) UintDefault(def uint) uint {
    if a, err := m.Uint(); err == nil {
        return a
    }
    return def
}
func (m *String) Uint8() (uint8, error) {
    n, err := strconv.ParseUint(m.val, 10, 8)
    if err != nil {
        return 0, err
    }
    return uint8(n), nil
}
func (m *String) Uint8Default(def uint8) uint8 {
    if a, err := m.Uint8(); err == nil {
        return a
    }
    return def
}
func (m *String) Uint16() (uint16, error) {
    n, err := strconv.ParseUint(m.val, 10, 16)
    if err != nil {
        return 0, err
    }
    return uint16(n), nil
}
func (m *String) Uint16Default(def uint16) uint16 {
    if a, err := m.Uint16(); err == nil {
        return a
    }
    return def
}
func (m *String) Uint32() (uint32, error) {
    n, err := strconv.ParseUint(m.val, 10, 32)
    if err != nil {
        return 0, err
    }
    return uint32(n), nil
}
func (m *String) Uint32Default(def uint32) uint32 {
    if a, err := m.Uint32(); err == nil {
        return a
    }
    return def
}
func (m *String) Uint64() (uint64, error) {
    return strconv.ParseUint(m.val, 10, 64)
}
func (m *String) Uint64Default(def uint64) uint64 {
    if a, err := m.Uint64(); err == nil {
        return a
    }
    return def
}

func (m *String) Float32() (float32, error) {
    f, err := strconv.ParseFloat(m.val, 32)
    if err != nil {
        return 0, err
    }
    return float32(f), nil
}
func (m *String) Float32Default(def float32) float32 {
    if a, err := m.Float32(); err == nil {
        return a
    }
    return def
}
func (m *String) Float64() (float64, error) {
    return strconv.ParseFloat(m.val, 64)
}
func (m *String) Float64Default(def float64) float64 {
    if a, err := m.Float64(); err == nil {
        return a
    }
    return def
}

func (m *String) Scan(v interface{}) error {
    return Scan(m.val, v)
}
