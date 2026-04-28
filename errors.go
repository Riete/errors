package errors

import (
	"errors"
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
	msg    string
	stacks []string
}

func (e err) Error() string {
	return e.msg
}

func (e err) Stack() string {
	var stacks []string
	prefix := ""
	for _, msg := range e.stacks {
		prefix += " "
		stacks = append(stacks, fmt.Sprintf("%s|- %s", prefix, msg))
	}
	return fmt.Sprintf("[ERROR]: %s\nTraceback:\n%s", e.msg, strings.Join(stacks, "\n"))
}

func (e *err) tryConvertMsgToStacks() {
	for _, msg := range strings.Split(e.msg, "|- ") {
		if strings.HasPrefix(msg, "[ERROR]: ") {
			e.msg = strings.TrimPrefix(strings.Split(msg, "\n")[0], "[ERROR]: ")
			continue
		}
		e.stacks = append(e.stacks, strings.Trim(msg, "\n "))
	}
}

func (e *err) trace(msg string) Error {
	e.msg = msg
	if strings.Contains(e.msg, "Traceback:\n") {
		e.tryConvertMsgToStacks()
	} else {
		e.stacks = append(e.stacks, fmt.Sprintf("%s %s", e.caller(3), msg))
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
		var se Error
		if errors.As(err, &se) {
			return e.trace(se.Stack())
		}
		return e.trace(err.Error())
	}
	return e
}

func (e *err) caller(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", file, line)
}

func New(msg string) Error {
	e := new(err)
	return e.trace(msg)
}

func NewFromErr(errs ...error) Error {
	if len(errs) == 0 {
		return nil
	}
	var e *err
	for _, i := range errs {
		if i == nil {
			continue
		}
		if e == nil {
			e = new(err)
		}
		var se Error
		if errors.As(i, &se) {
			_ = e.trace(se.Stack())
		} else {
			_ = e.trace(i.Error())
		}
	}
	if e == nil {
		return nil
	}
	return e
}
