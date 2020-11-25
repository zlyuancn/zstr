/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/25
   Description :
-------------------------------------------------
*/

package zstr

import (
	"testing"
)

const testSql = `
select * from a where
{&u.phone_number like}
    {&u.user_name like}
    &dev.district_id
    &u.gender
	{&u.create_time >=}
	&u.create_time
	{&u.create_time <}
	&u.create_time
    #start_time
    {#start_time}
    {#end_time}
    #end_time
    {#a}
    #a
	#b
    {#b}
    &au.is_ugc
	{&c in}
	{&d in}
	{&e notin}
	{&f in}
	{&g like}
	{&h like}
    &dev.device_platform
group by u.id
limit #size offset {#start a};`

var testData = map[string]interface{}{
	"":                 "xxx",
	"u.create_time":    "uc",
	"u.create_time[1]": "uc1",
	"u.create_time[2]": "uc2",
	"u.create_time[0]": "uc0",
	"start_time":       "st",
	"start_time[1]":    "st1",
	"end_time[1]":      "et[1]",
	"a[0]":             "av0",
	"b[1]":             "bv0",
	"d":                "dv",
	"e[0]":             []string{"ev0", "ev1"},
	"g":                "gv",
}

func BenchmarkRender(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Render(testSql, testData)
		}
	})
}

func BenchmarkSqlRender(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = SqlRender(testSql, testData)
		}
	})
}

func BenchmarkSqlParse(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = SqlParse(testSql, testData)
		}
	})
}
