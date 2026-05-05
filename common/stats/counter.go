package stats

import (
	"sync/atomic"
)

// Counter is a thread-safe counter for tracking statistics.
type Counter struct {
	value int64
}

// Add adds the given value to the counter and returns the new value.
func (c *Counter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

// Set sets the counter to the given value.
func (c *Counter) Set(val int64) int64 {
	old := atomic.SwapInt64(&c.value, val)
	return old
}

// Value returns the current value of the counter.
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Reset resets the counter to zero and returns the previous value.
func (c *Counter) Reset() int64 {
	return atomic.SwapInt64(&c.value, 0)
}
