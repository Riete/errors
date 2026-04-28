package errors

import (
	stderror "errors"
	"fmt"
	"runtime"
	"strings"
)

type Error interface {
	Error() string
	Stack() string
	Trace(string) Error
	Tracef(string, ...any) Error
	TraceErr(error) Error
}

type err struct {
	Msg    string
	Stacks []string
}

func (e err) Error() string {
	return e.Msg
}

func (e err) Stack() string {
	return fmt.Sprintf("[ERROR]: %s\nTraceback:\n%s", e.Msg, strings.Join(e.Stacks, "\n"))
}

func (e *err) tryConvertMsgToStacks() {
	c := 0
	for _, msg := range strings.Split(e.Msg, "|- ") {
		if strings.HasPrefix(msg, "[ERROR]: ") {
			e.Msg = strings.TrimPrefix(strings.Split(msg, "\n")[0], "[ERROR]: ")
			continue
		}
		stack := ""
		for i := 0; i < c; i++ {
			stack += " "
		}
		e.Stacks = append(e.Stacks, stack+"|- "+strings.Trim(msg, "\n "))
		c += 1
	}
}

func (e *err) trace(msg string) Error {
	e.Msg = msg
	if len(e.Stacks) == 0 {
		if strings.Contains(e.Msg, "Traceback:\n") { // maybe from stacked error
			e.tryConvertMsgToStacks()
		} else {
			e.Stacks = append(e.Stacks, fmt.Sprintf("|- %s %s", e.runtime(4), msg))
		}
	} else {
		stack := ""
		for i := 0; i < len(e.Stacks); i++ {
			stack += " "
		}
		stack += "|- " + e.runtime(3) + " " + msg
		e.Stacks = append(e.Stacks, stack)
	}
	return e
}

func (e *err) Trace(msg string) Error {
	return e.trace(msg)
}

func (e *err) Tracef(format string, v ...any) Error {
	return e.trace(fmt.Sprintf(format, v...))
}

func (e *err) TraceErr(err error) Error {
	if err != nil {
		return e.trace(err.Error())
	}
	return e
}

func (e err) runtime(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", file, line)
}

func New(msg string) Error {
	e := new(err)
	return e.Trace(msg)
}

func newFromErr(er error) *err {
	if er == nil {
		return nil
	}
	e := &err{}
	var se Error
	if stderror.As(er, &se) {
		return New(se.Stack()).(*err)
	}
	return e.Trace(er.Error()).(*err)
}

func NewFromErr(errors ...error) Error {
	if len(errors) == 0 {
		return nil
	}
	var stackMsg []string
	var errMsg string
	for _, i := range errors {
		if e := newFromErr(i); e != nil {
			stackMsg = append(stackMsg, e.Stacks...)
			errMsg = i.Error()
		}
	}
	if len(stackMsg) > 0 {
		stackMsg = append([]string{fmt.Sprintf("[ERROR]: %s\nTraceback:\n", errMsg)}, stackMsg...)
		return New(strings.Join(stackMsg, "\n"))
	}
	return nil
}
