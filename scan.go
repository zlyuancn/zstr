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

func Scan(s string, v interface{}) error {
    switch v := v.(type) {
    case nil:
        return fmt.Errorf("zstr: Scan(nil)")

    case *string:
        *v = s
        return nil
    case *[]byte:
        *v = []byte(s)
        return nil
    case *bool:
        var err error
        *v, err = ToBool(s)
        return err
    case *int:
        var err error
        *v, err = strconv.Atoi(s)
        return err
    case *int8:
        n, err := strconv.ParseInt(s, 10, 8)
        if err != nil {
            return err
        }
        *v = int8(n)
        return nil
    case *int16:
        n, err := strconv.ParseInt(s, 10, 16)
        if err != nil {
            return err
        }
        *v = int16(n)
        return nil
    case *int32:
        n, err := strconv.ParseInt(s, 10, 32)
        if err != nil {
            return err
        }
        *v = int32(n)
        return nil
    case *int64:
        n, err := strconv.ParseInt(s, 10, 64)
        if err != nil {
            return err
        }
        *v = n
        return nil

    case *uint:
        n, err := strconv.ParseUint(s, 10, 64)
        if err != nil {
            return err
        }
        *v = uint(n)
        return nil
    case *uint8:
        n, err := strconv.ParseUint(s, 10, 8)
        if err != nil {
            return err
        }
        *v = uint8(n)
        return nil
    case *uint16:
        n, err := strconv.ParseUint(s, 10, 16)
        if err != nil {
            return err
        }
        *v = uint16(n)
        return nil
    case *uint32:
        n, err := strconv.ParseUint(s, 10, 32)
        if err != nil {
            return err
        }
        *v = uint32(n)
        return nil
    case *uint64:
        n, err := strconv.ParseUint(s, 10, 64)
        if err != nil {
            return err
        }
        *v = n
        return nil

    case *float32:
        n, err := strconv.ParseFloat(s, 32)
        if err != nil {
            return err
        }
        *v = float32(n)
        return nil
    case *float64:
        n, err := strconv.ParseFloat(s, 64)
        if err != nil {
            return err
        }
        *v = n
        return nil

    case encoding.BinaryUnmarshaler:
        return v.UnmarshalBinary([]byte(s))

    default:
        return fmt.Errorf("zstr: 无法解码 %T, 考虑为它实现encoding.BinaryUnmarshaler接口", v)
    }
}
