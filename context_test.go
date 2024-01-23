package honeybadger

import (
	"context"
	"testing"
)

func TestContextUpdate(t *testing.T) {
	c := Context{"foo": "bar"}
	c.Update(Context{"foo": "baz"})
	if c["foo"] != "baz" {
		t.Errorf("Context should update values. expected=%#v actual=%#v", "baz", c["foo"])
	}
}

func TestContext(t *testing.T) {
	t.Run("setting values is allowed between reads", func(t *testing.T) {
		ctx := context.Background()
		ctx = Context{"foo": "bar"}.WithContext(ctx)

		stored := FromContext(ctx)
		if stored == nil {
			t.Fatalf("FromContext returned nil")
		}
		if stored["foo"] != "bar" {
			t.Errorf("stored[foo] expected=%q actual=%v", "bar", stored["foo"])
		}

		// Write a new key then we'll read from the ctx again and make sure it is
		// still set.
		stored["baz"] = "qux"
		stored = FromContext(ctx)
		if stored["baz"] != "qux" {
			t.Errorf("stored[baz] expected=%q actual=%v", "qux", stored["baz"])
		}
	})
}
