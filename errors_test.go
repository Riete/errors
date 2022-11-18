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
	e := New("a").Trace("b").Trace("c")
	t.Log(e.Stack())
	f := New(e.Stack()).Trace("d").Trace("f")
	t.Log(f.Stack())
}
