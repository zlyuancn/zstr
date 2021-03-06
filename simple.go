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
	"strconv"
)

func ToBytes(s string) []byte {
	return []byte(s)
}

func ToBool(s string) (bool, error) {
	switch s {
	case "1", "t", "T", "true", "TRUE", "True", "y", "Y", "yes", "YES", "Yes", "on", "ON", "On", "ok", "OK", "Ok":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "n", "N", "no", "NO", "No", "off", "OFF", "Off":
		return false, nil
	}
	return false, fmt.Errorf("数据\"%s\"无法转换为bool", s)
}
func ToBoolDefault(s string, def ...bool) bool {
	if a, err := ToBool(s); err == nil {
		return a
	}
	return len(def) > 0 && def[0]
}
func GetBool(any interface{}, def ...bool) bool {
	switch v := any.(type) {
	case nil:
		return false
	case bool:
		return v
	}

	s := anyToString(any)
	if a, err := ToBool(s); err == nil {
		return a
	}
	return len(def) > 0 && def[0]
}

func ToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
func ToIntDefault(s string, def ...int) int {
	if a, err := strconv.Atoi(s); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetInt(any interface{}, def ...int) int {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case int:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.Atoi(s); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToInt8(s string) (int8, error) {
	n, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(n), nil
}
func ToInt8Default(s string, def ...int8) int8 {
	if a, err := strconv.ParseInt(s, 10, 8); err == nil {
		return int8(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetInt8(any interface{}, def ...int8) int8 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case int8:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseInt(s, 10, 8); err == nil {
		return int8(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToInt16(s string) (int16, error) {
	n, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(n), nil
}
func ToInt16Default(s string, def ...int16) int16 {
	if a, err := strconv.ParseInt(s, 10, 16); err == nil {
		return int16(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetInt16(any interface{}, def ...int16) int16 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case int16:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseInt(s, 10, 16); err == nil {
		return int16(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToInt32(s string) (int32, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}
func ToInt32Default(s string, def ...int32) int32 {
	if a, err := strconv.ParseInt(s, 10, 32); err == nil {
		return int32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetInt32(any interface{}, def ...int32) int32 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case int32:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseInt(s, 10, 32); err == nil {
		return int32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
func ToInt64Default(s string, def ...int64) int64 {
	if a, err := strconv.ParseInt(s, 10, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetInt64(any interface{}, def ...int64) int64 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case int64:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseInt(s, 10, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ToUint(s string) (uint, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), err
}
func ToUintDefault(s string, def ...uint) uint {
	if a, err := strconv.ParseUint(s, 10, 64); err == nil {
		return uint(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetUint(any interface{}, def ...uint) uint {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case uint:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseUint(s, 10, 64); err == nil {
		return uint(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToUint8(s string) (uint8, error) {
	n, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(n), nil
}
func ToUint8Default(s string, def ...uint8) uint8 {
	if a, err := strconv.ParseUint(s, 10, 8); err == nil {
		return uint8(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetUint8(any interface{}, def ...uint8) uint8 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case uint8:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseUint(s, 10, 8); err == nil {
		return uint8(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToUint16(s string) (uint16, error) {
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}
func ToUint16Default(s string, def ...uint16) uint16 {
	if a, err := strconv.ParseUint(s, 10, 16); err == nil {
		return uint16(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetUint16(any interface{}, def ...uint16) uint16 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case uint16:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseUint(s, 10, 16); err == nil {
		return uint16(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToUint32(s string) (uint32, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(n), nil
}
func ToUint32Default(s string, def ...uint32) uint32 {
	if a, err := strconv.ParseUint(s, 10, 32); err == nil {
		return uint32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetUint32(any interface{}, def ...uint32) uint32 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case uint32:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseUint(s, 10, 32); err == nil {
		return uint32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
func ToUint64Default(s string, def ...uint64) uint64 {
	if a, err := strconv.ParseUint(s, 10, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetUint64(any interface{}, def ...uint64) uint64 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case uint64:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseUint(s, 10, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ToFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}
func ToFloat32Default(s string, def ...float32) float32 {
	if a, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetFloat32(any interface{}, def ...float32) float32 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case float32:
		return v
	case float64:
		return float32(v)
	}

	s := anyToString(any)
	if a, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func ToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
func ToFloat64Default(s string, def ...float64) float64 {
	if a, err := strconv.ParseFloat(s, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func GetFloat64(any interface{}, def ...float64) float64 {
	switch v := any.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case float32:
		return float64(v)
	case float64:
		return v
	}

	s := anyToString(any)
	if a, err := strconv.ParseFloat(s, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ToString(a interface{}) string {
	return anyToString(a)
}
func GetString(a interface{}) string {
	return anyToString(a)
}
