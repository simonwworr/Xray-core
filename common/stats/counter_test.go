package stats_test

import (
	"sync"
	"testing"

	"github.com/xtls/xray-core/common/stats"
)

func TestCounterAdd(t *testing.T) {
	c := &stats.Counter{}
	if v := c.Add(1); v != 1 {
		t.Errorf("expected 1, got %d", v)
	}
	if v := c.Add(5); v != 6 {
		t.Errorf("expected 6, got %d", v)
	}
}

func TestCounterReset(t *testing.T) {
	c := &stats.Counter{}
	c.Add(10)
	if old := c.Reset(); old != 10 {
		t.Errorf("expected old value 10, got %d", old)
	}
	if v := c.Value(); v != 0 {
		t.Errorf("expected 0 after reset, got %d", v)
	}
}

func TestCounterConcurrent(t *testing.T) {
	c := &stats.Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Add(1)
		}()
	}
	wg.Wait()
	if v := c.Value(); v != 100 {
		t.Errorf("expected 100, got %d", v)
	}
}

func TestCounterSet(t *testing.T) {
	c := &stats.Counter{}
	c.Add(5)
	old := c.Set(42)
	if old != 5 {
		t.Errorf("expected old value 5, got %d", old)
	}
	if v := c.Value(); v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}
