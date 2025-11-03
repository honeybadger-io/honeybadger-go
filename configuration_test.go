package honeybadger

import (
	"testing"
	"time"
)

type TestLogger struct{}

func (l *TestLogger) Printf(format string, v ...interface{}) {}

func TestUpdateConfig(t *testing.T) {
	config := &Configuration{}
	logger := &TestLogger{}
	backend := &TestBackend{}
	config.update(&Configuration{
		Logger:  logger,
		Backend: backend,
		Root:    "/tmp/foo",
	})

	if config.Logger != logger {
		t.Errorf("Expected config to update logger expected=%#v actual=%#v", logger, config.Logger)
	}
	if config.Backend != backend {
		t.Errorf("Expected config to update backend expected=%#v actual=%#v", backend, config.Backend)
	}
	if config.Root != "/tmp/foo" {
		t.Errorf("Expected config to update root expected=%#v actual=%#v", "/tmp/foo", config.Root)
	}
}

func TestReplaceConfigPointer(t *testing.T) {
	config := Configuration{Root: "/tmp/foo"}
	root := &config.Root
	config = Configuration{Root: "/tmp/bar"}
	if *root != "/tmp/bar" {
		t.Errorf("Expected updated config to update pointer expected=%#v actual=%#v", "/tmp/bar", *root)
	}
}

func TestGetEnv_BasicTypesAndFallbacks(t *testing.T) {
	t.Run("string with env", func(t *testing.T) {
		t.Setenv("X", "abc")
		got := GetEnv[string]("X")
		if got != "abc" {
			t.Fatalf("want abc, got %q", got)
		}
	})

	t.Run("string with value fallback", func(t *testing.T) {
		t.Setenv("X", "")
		got := GetEnv[string]("X", "def")
		if got != "def" {
			t.Fatalf("want def, got %q", got)
		}
	})

	t.Run("string with func fallback (lazy)", func(t *testing.T) {
		t.Setenv("X", "")
		calls := 0
		got := GetEnv[string]("X", func() string { calls++; return "zzz" })
		if calls != 1 || got != "zzz" {
			t.Fatalf("fallback not called once or wrong value: calls=%d got=%q", calls, got)
		}
	})

	t.Run("string with func fallback NOT called if env present", func(t *testing.T) {
		t.Setenv("X", "live")
		called := false
		got := GetEnv[string]("X", func() string { called = true; return "nope" })
		if called {
			t.Fatal("fallback was called unexpectedly")
		}
		if got != "live" {
			t.Fatalf("want live, got %q", got)
		}
	})

	t.Run("int parse ok", func(t *testing.T) {
		t.Setenv("X", "42")
		got := GetEnv[int]("X")
		if got != 42 {
			t.Fatalf("want 42, got %d", got)
		}
	})

	t.Run("int parse bad -> value fallback", func(t *testing.T) {
		t.Setenv("X", "nope")
		got := GetEnv[int]("X", 7)
		if got != 7 {
			t.Fatalf("want fallback 7, got %d", got)
		}
	})

	t.Run("float64 parse ok", func(t *testing.T) {
		t.Setenv("X", "3.14")
		got := GetEnv[float64]("X")
		if got != 3.14 {
			t.Fatalf("want 3.14, got %v", got)
		}
	})

	t.Run("bool parse ok", func(t *testing.T) {
		t.Setenv("X", "true")
		got := GetEnv[bool]("X")
		if !got {
			t.Fatalf("want true, got false")
		}

		got = GetEnv[bool]("Y")
		if got {
			t.Fatalf("want false, got true")
		}
	})

	t.Run("duration parse ok", func(t *testing.T) {
		t.Setenv("X", "120000000")
		got := GetEnv[time.Duration]("X")
		if got != 120*time.Millisecond {
			t.Fatalf("want 120ms, got %v", got)
		}

		t.Setenv("Y", "150ms")
		got = GetEnv[time.Duration]("Y")
		if got != 150*time.Millisecond {
			t.Fatalf("want 150ms, got %v", got)
		}
	})

	t.Run("zero value when missing and no fallback", func(t *testing.T) {
		t.Setenv("X", "")
		if got := GetEnv[int]("X"); got != 0 {
			t.Fatalf("want 0, got %d", got)
		}
		if got := GetEnv[string]("X"); got != "" {
			t.Fatalf("want empty string, got %q", got)
		}
		if got := GetEnv[bool]("X"); got != false {
			t.Fatalf("want false, got %v", got)
		}
	})
}
