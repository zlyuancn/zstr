/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/7/17
   Description :
-------------------------------------------------
*/

package zstr

import (
	"fmt"
	"strconv"
)

func anyToString(a interface{}) string {
	switch v := a.(type) {

	case nil:
		return "nil"

	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		if v {
			return "true"
		}
		return "false"

	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)

	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	}
	return fmt.Sprint(a)
}
