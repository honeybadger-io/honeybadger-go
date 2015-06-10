package honeybadger

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

const MaxFrames = 20

type Frame struct {
	Number string `json:"number"`
	File   string `json:"file"`
	Method string `json:"method"`
}

type Error struct {
	err     interface{}
	Message string
	Class   string
	Stack   []*Frame
}

func (e Error) Error() string {
	return e.Message
}

func newError(thing interface{}, stackOffset int) Error {
	var err error

	switch thing := thing.(type) {
	case Error:
		return thing
	case error:
		err = thing
	default:
		err = fmt.Errorf("%v", thing)
	}

	return Error{
		err:     err,
		Message: err.Error(),
		Class:   reflect.TypeOf(err).String(),
		Stack:   generateStack(stackOffset),
	}
}

func generateStack(offset int) (frames []*Frame) {
	stack := make([]uintptr, MaxFrames)
	length := runtime.Callers(2+offset, stack[:])
	for _, pc := range stack[:length] {
		f := runtime.FuncForPC(pc)
		if f == nil {
			continue
		}
		file, line := f.FileLine(pc)
		frame := &Frame{
			File:   file,
			Number: strconv.Itoa(line),
			Method: f.Name(),
		}
		frames = append(frames, frame)
	}

	return
}
