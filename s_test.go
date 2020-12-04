/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/11
   Description :
-------------------------------------------------
*/

package zstr

import (
	"testing"
)

func TestString_Bool(t *testing.T) {
	s := String("1")
	b, err := s.Bool()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
}

func TestString_Scan(t *testing.T) {
	s := String("yes")
	b := false
	if err := s.Scan(&b); err != nil {
		t.Fatal(err)
	}
	t.Log(b)
}
