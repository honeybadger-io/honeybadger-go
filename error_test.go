package honeybadger

import (
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
