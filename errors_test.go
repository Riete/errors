package errors

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	e := New("error1")
	e = e.Trace("error2")
	e = e.Trace("error3")
	e2 := NewFromErr(errors.New("aa"))
	e3 := e2.Trace("bb")
	e3.Trace("cc")
	t.Log(e.Stack())
	t.Log(e2.Stack())
	t.Log(e3.Stack())

}

func TestErrors2(t *testing.T) {
	e := New("a\nline1\nline3").Trace("b\nline3\nline4").Trace("c\nline5\nline6")
	t.Log(e.Stack())
	f := New(e.Stack()).Trace("d\n4").Trace("f\n5")
	t.Log(f.Stack())
}
