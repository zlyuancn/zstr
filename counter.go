/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/25
   Description :
-------------------------------------------------
*/

package zstr

// 计数器, 注意: 非并发安全
type counter struct {
	data map[string]int
	def  int // 默认值
}

// 创建一个计数器, 非线程安全
func newCounter(initValue ...int) *counter {
	c := &counter{
		data: make(map[string]int),
	}
	if len(initValue) > 0 {
		c.def = initValue[0]
	}
	return c
}

func (c *counter) Incr(key string) int {
	return c.IncrBy(key, 1)
}

func (c *counter) IncrBy(key string, num int) int {
	v, ok := c.data[key]
	if !ok {
		v = c.def
	}
	v += num
	c.data[key] = v
	return v
}

func (c *counter) Get(key string) int {
	v, ok := c.data[key]
	if ok {
		return v
	}
	return c.def
}
