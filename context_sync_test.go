package honeybadger

import (
	"testing"
	"sync"
)

func TestContextSync(t *testing.T) {
	var wg sync.WaitGroup

	instance   := newContextSync()
	newContext := Context{"foo":"bar"}

	wg.Add(2)

	go update(&wg, instance, newContext)
	go update(&wg, instance, newContext)

	wg.Wait()

	context := instance.internal

	if context["foo"] != "bar" {
		t.Errorf("Expected context value. expected=%#v result=%#v", "bar", context["foo"])
	}
}

func update(wg *sync.WaitGroup, instance *contextSync, context Context) {
	instance.Update(context)
	wg.Done()
}
