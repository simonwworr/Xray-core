package stats_test

import (
	"testing"

	. "github.com/xtls/xray-core/common/stats"
)

func TestManagerRegisterCounter(t *testing.T) {
	m := NewManager()

	c, err := m.RegisterCounter("test.counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil counter")
	}

	// Registering again should return the same counter
	c2, err := m.RegisterCounter("test.counter")
	if err != nil {
		t.Fatalf("unexpected error on second register: %v", err)
	}
	if c != c2 {
		t.Fatal("expected same counter instance on duplicate registration")
	}
}

func TestManagerGetCounter(t *testing.T) {
	m := NewManager()

	if got := m.GetCounter("nonexistent"); got != nil {
		t.Fatal("expected nil for nonexistent counter")
	}

	_, _ = m.RegisterCounter("test.counter")
	if got := m.GetCounter("test.counter"); got == nil {
		t.Fatal("expected non-nil counter after registration")
	}
}

func TestManagerUnregisterCounter(t *testing.T) {
	m := NewManager()

	_, _ = m.RegisterCounter("test.counter")
	if err := m.UnregisterCounter("test.counter"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.GetCounter("test.counter"); got != nil {
		t.Fatal("expected nil after unregister")
	}
}

func TestManagerRegisterChannel(t *testing.T) {
	m := NewManager()
	if err := m.Start(); err != nil {
		t.Fatalf("failed to start manager: %v", err)
	}
	defer m.Close() //nolint:errcheck

	ch, err := m.RegisterChannel("test.channel")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}

	// Duplicate registration should return the same channel
	ch2, err := m.RegisterChannel("test.channel")
	if err != nil {
		t.Fatalf("unexpected error on second register: %v", err)
	}
	if ch != ch2 {
		t.Fatal("expected same channel instance on duplicate registration")
	}
}

func TestManagerGetChannel(t *testing.T) {
	m := NewManager()

	if got := m.GetChannel("nonexistent"); got != nil {
		t.Fatal("expected nil for nonexistent channel")
	}

	if err := m.Start(); err != nil {
		t.Fatalf("failed to start manager: %v", err)
	}
	defer m.Close() //nolint:errcheck

	_, _ = m.RegisterChannel("test.channel")
	if got := m.GetChannel("test.channel"); got == nil {
		t.Fatal("expected non-nil channel after registration")
	}
}
