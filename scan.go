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
	"reflect"
	"strconv"
)

func Scan(s string, outPtr interface{}) (err error) {
	switch p := outPtr.(type) {
	case nil:
		return fmt.Errorf("zstr: Scan(nil)")

	case *string:
		*p = s
	case *[]byte:
		*p = []byte(s)
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

	case encoding.BinaryUnmarshaler:
		return p.UnmarshalBinary([]byte(s))

	default:
		return fmt.Errorf("zstr: 无法解码 %T, 考虑为它实现encoding.BinaryUnmarshaler接口", p)
	}
	return
}

// 扫描任何值到任何
func ScanAny(any, outPtr interface{}) (err error) {
	switch p := outPtr.(type) {
	case nil:
		return fmt.Errorf("zstr: Scan(nil)")

	case *string:
		*p = GetString(any)
	case *[]byte:
		*p = []byte(GetString(any))
	case *bool:
		*p, err = ToBool(any)
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

	case encoding.BinaryUnmarshaler:
		return p.UnmarshalBinary([]byte(s))

	default:
		return fmt.Errorf("zstr: 无法解码 %T, 考虑为它实现encoding.BinaryUnmarshaler接口", p)
	}
	return
}

// 将map数据扫描到结构体中
func ScanMap(m map[string]interface{}, outPtr interface{}) (err error) {
	aType := reflect.TypeOf(outPtr)
	if aType.Kind() != reflect.Ptr {
		return fmt.Errorf("outPtr必须是带指针的结构体")
	}
	aType = aType.Elem()
	if aType.Kind() != reflect.Struct {
		return fmt.Errorf("outPtr必须是带指针的结构体")
	}
	aValue := reflect.ValueOf(outPtr).Elem()

	// 扫描结构体
	fields := make(map[string]*reflect.Value, aType.NumField())
	for i := 0; i < aType.NumField(); i++ {
		field := aType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		k, ok := field.Tag.Lookup("json")
		if !ok {
			k = field.Name
		}
		v := aValue.Field(i)
		fields[k] = &v
	}

	for k, v := range m {
		fieldValue, ok := fields[k]
		if !ok {
			continue
		}

		fieldValue.Interface()
	}

	return nil
}
