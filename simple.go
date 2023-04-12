/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/14
   Description :
-------------------------------------------------
*/

package zstr

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

func ToBytes(s string) []byte {
	return StringToBytes(&s)
}

func ToBool(any interface{}) (bool, error) {
	switch v := any.(type) {
	case nil:
		return false, nil
	case bool:
		return v, nil
	}
	s := anyToString(any)
	switch s {
	case "1", "t", "T", "true", "TRUE", "True", "y", "Y", "yes", "YES", "Yes",
		"on", "ON", "On", "ok", "OK", "Ok",
		"enabled", "ENABLED", "Enabled",
		"open", "OPEN", "Open":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "n", "N", "no", "NO", "No",
		"off", "OFF", "Off", "cancel", "CANCEL", "Cancel",
		"disable", "DISABLE", "Disable",
		"close", "CLOSE", "Close",
		"", "nil", "Nil", "NIL", "null", "Null", "NULL", "none", "None", "NONE":
		return false, nil
	}
	return false, fmt.Errorf("数据\"%s\"无法转换为bool", s)
}
func GetBool(any interface{}, def ...bool) bool {
	if a, err := ToBool(any); err == nil {
		return a
	}
	return len(def) > 0 && def[0]
}

func ToInt(any interface{}) (int, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case int:
		return v, nil
	}

	s := anyToString(any)
	return strconv.Atoi(s)
}
func GetInt(any interface{}, def ...int) int {
	if a, err := ToInt(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToInt8(any interface{}) (int8, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case int8:
		return v, nil
	}

	s := anyToString(any)
	n, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(n), nil
}
func GetInt8(any interface{}, def ...int8) int8 {
	if a, err := ToInt8(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToInt16(any interface{}) (int16, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case int16:
		return v, nil
	}

	s := anyToString(any)
	n, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(n), nil
}
func GetInt16(any interface{}, def ...int16) int16 {
	if a, err := ToInt16(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToInt32(any interface{}) (int32, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case int32:
		return v, nil
	}

	s := anyToString(any)
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}
func GetInt32(any interface{}, def ...int32) int32 {
	if a, err := ToInt32(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToInt64(any interface{}) (int64, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case int64:
		return v, nil
	}

	s := anyToString(any)
	return strconv.ParseInt(s, 10, 64)
}
func GetInt64(any interface{}, def ...int64) int64 {
	if a, err := ToInt64(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ToUint(any interface{}) (uint, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case uint:
		return v, nil
	}

	s := anyToString(any)
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), err
}
func GetUint(any interface{}, def ...uint) uint {
	if a, err := ToUint(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToUint8(any interface{}) (uint8, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case uint8:
		return v, nil
	}

	s := anyToString(any)
	n, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(n), nil
}
func GetUint8(any interface{}, def ...uint8) uint8 {
	if a, err := ToUint8(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToUint16(any interface{}) (uint16, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case uint16:
		return v, nil
	}

	s := anyToString(any)
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}
func GetUint16(any interface{}, def ...uint16) uint16 {
	if a, err := ToUint16(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToUint32(any interface{}) (uint32, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case uint32:
		return v, nil
	}

	s := anyToString(any)
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(n), nil
}
func GetUint32(any interface{}, def ...uint32) uint32 {
	if a, err := ToUint32(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToUint64(any interface{}) (uint64, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case uint64:
		return v, nil
	}

	s := anyToString(any)
	return strconv.ParseUint(s, 10, 64)
}
func GetUint64(any interface{}, def ...uint64) uint64 {
	if a, err := ToUint64(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ToFloat32(any interface{}) (float32, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case float32:
		return v, nil
	case float64:
		return float32(v), nil
	}

	s := anyToString(any)
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}
func GetFloat32(any interface{}, def ...float32) float32 {
	if a, err := ToFloat32(any); err == nil {
		return float32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToFloat64(any interface{}) (float64, error) {
	switch v := any.(type) {
	case nil:
		return 0, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	}

	s := anyToString(any)
	return strconv.ParseFloat(s, 64)
}
func GetFloat64(any interface{}, def ...float64) float64 {
	if a, err := ToFloat64(any); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ToString(a interface{}, nilToEmpty ...bool) string {
	return anyToString(a, nilToEmpty...)
}
func GetString(a interface{}, nilToEmpty ...bool) string {
	return anyToString(a, nilToEmpty...)
}

// string转bytes, 转换后的bytes禁止写, 否则产生运行故障
func StringToBytes(s *string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// bytes转string
func BytesToString(b []byte) *string {
	return (*string)(unsafe.Pointer(&b))
}
