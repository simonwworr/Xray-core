package stats_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/xtls/xray-core/common/stats"
)

func TestSessionStatsUpDownlink(t *testing.T) {
	s := NewSessionStats()
	s.AddUplink(100)
	s.AddUplink(50)
	s.AddDownlink(200)

	if got := s.Uplink(); got != 150 {
		t.Errorf("expected uplink 150, got %d", got)
	}
	if got := s.Downlink(); got != 200 {
		t.Errorf("expected downlink 200, got %d", got)
	}
}

func TestSessionStatsClose(t *testing.T) {
	s := NewSessionStats()
	if s.IsClosed() {
		t.Fatal("session should not be closed initially")
	}

	time.Sleep(10 * time.Millisecond)
	s.Close()

	if !s.IsClosed() {
		t.Fatal("session should be closed after Close()")
	}
	if d := s.Duration(); d < 10*time.Millisecond {
		t.Errorf("expected duration >= 10ms, got %v", d)
	}

	// Calling Close again should not change the end time.
	first := s.Duration()
	time.Sleep(5 * time.Millisecond)
	s.Close()
	second := s.Duration()
	if first != second {
		t.Errorf("second Close() changed duration: %v vs %v", first, second)
	}
}

func TestSessionStatsConcurrent(t *testing.T) {
	s := NewSessionStats()
	var wg sync.WaitGroup
	const goroutines = 50
	const bytesPerGoroutine = int64(1000)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.AddUplink(bytesPerGoroutine)
			s.AddDownlink(bytesPerGoroutine)
		}()
	}
	wg.Wait()

	expected := int64(goroutines) * bytesPerGoroutine
	if got := s.Uplink(); got != expected {
		t.Errorf("expected uplink %d, got %d", expected, got)
	}
	if got := s.Downlink(); got != expected {
		t.Errorf("expected downlink %d, got %d", expected, got)
	}
}

func TestSessionStatsStartTime(t *testing.T) {
	before := time.Now()
	s := NewSessionStats()
	after := time.Now()

	st := s.StartTime()
	if st.Before(before) || st.After(after) {
		t.Errorf("start time %v not in expected range [%v, %v]", st, before, after)
	}
}
