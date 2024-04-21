package honeybadger

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestNewErrorTrace(t *testing.T) {
	fn := func() Error {
		return NewError("Error msg")
	}

	err := fn()

	// The stack should look like this:
	//   github.com/honeybadger-io/honeybadger-go.TestNewErrorTrace.func1
	//   github.com/honeybadger-io/honeybadger-go.TestNewErrorTrace
	//   testing.tRunner
	//   runtime.goexit
	if len(err.Stack) < 3 {
		t.Errorf("Expected to generate full trace")
	}

	// Checks that the top top methods are the (inlined) fn and the test Method
	expected := []string{
		".TestNewErrorTrace.func1",
		".TestNewErrorTrace",
	}

	for i, suffix := range expected {
		method := err.Stack[i].Method

		if !strings.HasSuffix(method, suffix) {
			// Logs the stack to give some context about the error
			for j, stack := range err.Stack {
				t.Logf("%d: %s", j, stack.Method)
			}

			t.Fatalf("stack[%d].Method expected_suffix=%q actual=%q", i, suffix, method)
		}
	}
}

type customerror struct {
	error
	callers []uintptr
}

func (t customerror) Callers() []uintptr {
	return t.callers
}

func newcustomerror() customerror {
	stack := make([]uintptr, maxFrames)
	length := runtime.Callers(1, stack[:])
	return customerror{
		error:   fmt.Errorf("hello world"),
		callers: stack[:length],
	}
}

func TestNewErrorCustomTrace(t *testing.T) {
	err := NewError(newcustomerror())

	// The stack should look like this:
	//   github.com/honeybadger-io/honeybadger-go.newcustomerror
	//   github.com/honeybadger-io/honeybadger-go.TestNewErrorCustomTrace
	//   testing.tRunner
	//   runtime.goexit
	if len(err.Stack) < 3 {
		t.Errorf("Expected to generate full trace")
	}

	// Checks that the top top methods are the (inlined) fn and the test Method
	expected := []string{
		".newcustomerror",
		".TestNewErrorCustomTrace",
	}

	for i, suffix := range expected {
		method := err.Stack[i].Method
		if !strings.HasSuffix(method, suffix) {
			// Logs the stack to give some context about the error
			for j, stack := range err.Stack {
				t.Logf("%d: %s", j, stack.Method)
			}

			t.Fatalf("stack[%d].Method expected_suffix=%q actual=%q", i, suffix, method)
		}
	}
}
