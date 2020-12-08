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

func newCounter(def ...int) *counter {
	c := &counter{
		data: make(map[string]int),
	}
	if len(def) > 0 {
		c.def = def[0]
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
	return c.data[key]
}
