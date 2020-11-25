/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/25
   Description :
-------------------------------------------------
*/

package zstr

type counter map[string]int

func newCounter() counter {
	return make(map[string]int)
}

func (c counter) Incr(key string) int {
	return c.IncrBy(key, 1)
}

func (c counter) IncrBy(key string, num int) int {
	if v, ok := c[key]; ok {
		v += num
		c[key] = v
		return v
	}
	c[key] = num
	return num
}

func (c counter) Get(key string) int {
	return c[key]
}
