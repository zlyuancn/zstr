/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/11
   Description :
-------------------------------------------------
*/

package zstr

import (
	"strconv"
)

type String string

func (m String) String() string {
	return string(m)
}
func (m String) Val() string {
	return string(m)
}
func (m String) Bytes() []byte {
	return []byte(m)
}

func (m String) Bool() (bool, error) {
	return ToBool(string(m))
}
func (m String) BoolDefault(def ...bool) bool {
	if a, err := m.Bool(); err == nil {
		return a
	}
	return len(def) > 0 && def[0]
}
func (m String) GetBool(def ...bool) bool {
	if a, err := m.Bool(); err == nil {
		return a
	}
	return len(def) > 0 && def[0]
}

func (m String) Int() (int, error) {
	return strconv.Atoi(string(m))
}
func (m String) IntDefault(def ...int) int {
	if a, err := strconv.Atoi(string(m)); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetInt(def ...int) int {
	if a, err := strconv.Atoi(string(m)); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) Int8() (int8, error) {
	n, err := strconv.ParseInt(string(m), 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(n), nil
}
func (m String) Int8Default(def ...int8) int8 {
	if a, err := strconv.ParseInt(string(m), 10, 8); err == nil {
		return int8(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetInt8(def ...int8) int8 {
	if a, err := strconv.ParseInt(string(m), 10, 8); err == nil {
		return int8(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) Int16() (int16, error) {
	n, err := strconv.ParseInt(string(m), 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(n), nil
}
func (m String) Int16Default(def ...int16) int16 {
	if a, err := strconv.ParseInt(string(m), 10, 16); err == nil {
		return int16(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetInt16(def ...int16) int16 {
	if a, err := strconv.ParseInt(string(m), 10, 16); err == nil {
		return int16(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) Int32() (int32, error) {
	n, err := strconv.ParseInt(string(m), 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}
func (m String) Int32Default(def ...int32) int32 {
	if a, err := strconv.ParseInt(string(m), 10, 32); err == nil {
		return int32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetInt32(def ...int32) int32 {
	if a, err := strconv.ParseInt(string(m), 10, 32); err == nil {
		return int32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) Int64() (int64, error) {
	return strconv.ParseInt(string(m), 10, 64)
}
func (m String) Int64Default(def ...int64) int64 {
	if a, err := strconv.ParseInt(string(m), 10, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetInt64(def ...int64) int64 {
	if a, err := strconv.ParseInt(string(m), 10, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (m String) Uint() (uint, error) {
	n, err := strconv.ParseUint(string(m), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), err
}
func (m String) UintDefault(def ...uint) uint {
	if a, err := strconv.ParseUint(string(m), 10, 64); err == nil {
		return uint(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetUint(def ...uint) uint {
	if a, err := strconv.ParseUint(string(m), 10, 64); err == nil {
		return uint(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) Uint8() (uint8, error) {
	n, err := strconv.ParseUint(string(m), 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(n), nil
}
func (m String) Uint8Default(def ...uint8) uint8 {
	if a, err := strconv.ParseUint(string(m), 10, 8); err == nil {
		return uint8(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetUint8(def ...uint8) uint8 {
	if a, err := strconv.ParseUint(string(m), 10, 8); err == nil {
		return uint8(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) Uint16() (uint16, error) {
	n, err := strconv.ParseUint(string(m), 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}
func (m String) Uint16Default(def ...uint16) uint16 {
	if a, err := strconv.ParseUint(string(m), 10, 16); err == nil {
		return uint16(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetUint16(def ...uint16) uint16 {
	if a, err := strconv.ParseUint(string(m), 10, 16); err == nil {
		return uint16(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) Uint32() (uint32, error) {
	n, err := strconv.ParseUint(string(m), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(n), nil
}
func (m String) Uint32Default(def ...uint32) uint32 {
	if a, err := strconv.ParseUint(string(m), 10, 32); err == nil {
		return uint32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetUint32(def ...uint32) uint32 {
	if a, err := strconv.ParseUint(string(m), 10, 32); err == nil {
		return uint32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) Uint64() (uint64, error) {
	return strconv.ParseUint(string(m), 10, 64)
}
func (m String) Uint64Default(def ...uint64) uint64 {
	if a, err := strconv.ParseUint(string(m), 10, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetUint64(def ...uint64) uint64 {
	if a, err := strconv.ParseUint(string(m), 10, 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (m String) Float32() (float32, error) {
	f, err := strconv.ParseFloat(string(m), 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}
func (m String) Float32Default(def ...float32) float32 {
	if a, err := strconv.ParseFloat(string(m), 32); err == nil {
		return float32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetFloat32(def ...float32) float32 {
	if a, err := strconv.ParseFloat(string(m), 32); err == nil {
		return float32(a)
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) Float64() (float64, error) {
	return strconv.ParseFloat(string(m), 64)
}
func (m String) Float64Default(def ...float64) float64 {
	if a, err := strconv.ParseFloat(string(m), 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}
func (m String) GetFloat64(def ...float64) float64 {
	if a, err := strconv.ParseFloat(string(m), 64); err == nil {
		return a
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (m String) Scan(outPtr interface{}) error {
	return Scan(string(m), outPtr)
}
