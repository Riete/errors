package errors

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	e := New("error1")
	e = e.Trace("error2")
	e = e.Trace("error3")
	t.Log(e.Error())
	t.Log(e.Stack())
}

func TestErrors2(t *testing.T) {
	e1 := New("error1").Trace("error2")
	e3 := New(e1.Stack()).Trace("error3")
	t.Log(e3.Error())
	t.Log(e3.Stack())
}

func TestErrors3(t *testing.T) {
	e1 := New("error1")
	e2 := New("error2").TraceErr(e1)
	e3 := New("error3").TraceErr(e2).TraceErr(errors.New("xx"))
	t.Log(e3.Error())
	t.Log(e3.Stack())
}

func TestErrors4(t *testing.T) {
	e1 := New("error1").Trace("error2")
	e2 := New("error3").Trace("error4")
	e3 := New("error5").TraceErr(e2).Trace(e1.Stack())
	t.Log(e3.Error())
	t.Log(e3.Stack())
}

func TestError5(t *testing.T) {
	e1 := New("error1")
	e2 := New("error2")
	e3 := errors.New("error3")
	e4 := NewFromErr(e1, e2, e3).Trace("error4")
	t.Log(e4.Error() == "error4")
	t.Log(e4.Stack())
}

func TestError6(t *testing.T) {
	t.Log(NewFromErr() == nil)
	t.Log(NewFromErr(nil) == nil)
	var e error
	var se Error
	t.Log(NewFromErr(e) == nil)
	t.Log(NewFromErr(se) == nil)
	t.Log(NewFromErr(e, se) == nil)
}
