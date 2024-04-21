package honeybadger

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

const maxFrames = 20

// Frame represent a stack frame inside of a Honeybadger backtrace.
type Frame struct {
	Number string `json:"number"`
	File   string `json:"file"`
	Method string `json:"method"`
}

// Error provides more structured information about a Go error.
type Error struct {
	err     error
	Message string
	Class   string
	Stack   []*Frame
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Error() string {
	return e.Message
}

type stacked interface {
	Callers() []uintptr
}

func NewError(msg interface{}) Error {
	return newError(msg, 2)
}

func newError(thing interface{}, stackOffset int) Error {
	var err error

	switch t := thing.(type) {
	case Error:
		return t
	case error:
		err = t
	default:
		err = fmt.Errorf("%v", t)
	}

	return Error{
		err:     err,
		Message: err.Error(),
		Class:   reflect.TypeOf(err).String(),
		Stack:   generateStack(autostack(err, stackOffset)),
	}
}

func autostack(err error, offset int) []uintptr {
	var s stacked

	if errors.As(err, &s) {
		return s.Callers()
	}

	stack := make([]uintptr, maxFrames)
	length := runtime.Callers(2+offset, stack[:])
	return stack[:length]
}

func generateStack(stack []uintptr) []*Frame {
	frames := runtime.CallersFrames(stack)
	result := make([]*Frame, 0, len(stack))

	for {
		frame, more := frames.Next()

		result = append(result, &Frame{
			File:   frame.File,
			Number: strconv.Itoa(frame.Line),
			Method: frame.Function,
		})

		if !more {
			break
		}
	}

	return result
}
