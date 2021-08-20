/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/11
   Description :
-------------------------------------------------
*/

package zstr

import (
	"encoding"
	"fmt"
	"strconv"
)

type AnyUnmarshaler interface {
	UnmarshalAny(any interface{}) error
}

func Scan(s string, outPtr interface{}) (err error) {
	switch p := outPtr.(type) {
	case encoding.BinaryUnmarshaler:
		return p.UnmarshalBinary(StringToBytes(&s))
	}
	switch p := outPtr.(type) {
	case AnyUnmarshaler:
		return p.UnmarshalAny(s)
	}

	switch p := outPtr.(type) {
	case nil:
		return fmt.Errorf("zstr.Scan(nil)")

	case *string:
		*p = s
	case *[]byte:
		*p = StringToBytes(&s)
	case *bool:
		*p, err = ToBool(s)
	case *int:
		*p, err = strconv.Atoi(s)
	case *int8:
		var n int64
		n, err = strconv.ParseInt(s, 10, 8)
		*p = int8(n)
	case *int16:
		var n int64
		n, err = strconv.ParseInt(s, 10, 16)
		*p = int16(n)
	case *int32:
		var n int64
		n, err = strconv.ParseInt(s, 10, 32)
		*p = int32(n)
	case *int64:
		*p, err = strconv.ParseInt(s, 10, 64)

	case *uint:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 64)
		*p = uint(n)
	case *uint8:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 8)
		*p = uint8(n)
	case *uint16:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 16)
		*p = uint16(n)
	case *uint32:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 32)
		*p = uint32(n)
	case *uint64:
		*p, err = strconv.ParseUint(s, 10, 64)

	case *float32:
		var n float64
		n, err = strconv.ParseFloat(s, 32)
		*p = float32(n)
	case *float64:
		*p, err = strconv.ParseFloat(s, 64)

	default:
		return fmt.Errorf("zstr.Scan(%T)无法解码, 考虑为它实现encoding.BinaryUnmarshaler接口或zstr.AnyUnmarshaler", p)
	}
	return
}

// 扫描任何值到任何, 不支持切片, 数组, Map, Struct
func ScanAny(any, outPtr interface{}) (err error) {
	switch t := any.(type) {
	case []byte:
		return Scan(*BytesToString(t), outPtr)
	case string:
		return Scan(t, outPtr)
	}

	switch p := outPtr.(type) {
	case AnyUnmarshaler:
		return p.UnmarshalAny(any)
	}

	switch p := outPtr.(type) {
	case nil:
		return fmt.Errorf("zstr.Scan(nil)")

	case *string:
		*p = GetString(any)
	case *[]byte:
		s := GetString(any)
		*p = StringToBytes(&s)
	case *bool:
		*p, err = ToBool(any)
	case *int:
		*p, err = ToInt(any)
	case *int8:
		*p, err = ToInt8(any)
	case *int16:
		*p, err = ToInt16(any)
	case *int32:
		*p, err = ToInt32(any)
	case *int64:
		*p, err = ToInt64(any)

	case *uint:
		*p, err = ToUint(any)
	case *uint8:
		*p, err = ToUint8(any)
	case *uint16:
		*p, err = ToUint16(any)
	case *uint32:
		*p, err = ToUint32(any)
	case *uint64:
		*p, err = ToUint64(any)

	case *float32:
		*p, err = ToFloat32(any)
	case *float64:
		*p, err = ToFloat64(any)

	default:
		return fmt.Errorf("zstr.Scan(%T)无法解码, 考虑为它实现zstr.AnyUnmarshaler接口", p)
	}
	return
}
