package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type Error interface {
	Error() string
	Stack() string
	Trace(string) Error
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

func (e *err) Trace(msg string) Error {
	e.Msg = msg
	if len(e.Stacks) == 0 {
		e.Stacks = append(e.Stacks, fmt.Sprintf("%s %s", e.runtime(3), msg))
	} else {
		stack := ""
		for i := 0; i < len(e.Stacks); i++ {
			stack += " "
		}
		stack += "|- " + e.runtime(2) + " " + msg
		e.Stacks = append(e.Stacks, stack)
	}
	return e
}

func (e err) runtime(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", file, line)
}

func New(msg string) Error {
	e := &err{}
	return e.Trace(msg)
}

func NewFromErr(er error) Error {
	if er == nil {
		return nil
	}
	e := &err{}
	return e.Trace(er.Error())
}