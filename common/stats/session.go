package stats

import (
	"sync"
	"sync/atomic"
	"time"
)

// SessionStats tracks statistics for a single proxy session.
type SessionStats struct {
	uplink   int64
	downlink int64
	start    time.Time
	end      time.Time
	closed   uint32
	mu       sync.RWMutex
}

// NewSessionStats creates a new SessionStats with the start time set to now.
func NewSessionStats() *SessionStats {
	return &SessionStats{
		start: time.Now(),
	}
}

// AddUplink adds n bytes to the uplink counter.
func (s *SessionStats) AddUplink(n int64) {
	atomic.AddInt64(&s.uplink, n)
}

// AddDownlink adds n bytes to the downlink counter.
func (s *SessionStats) AddDownlink(n int64) {
	atomic.AddInt64(&s.downlink, n)
}

// Uplink returns the total uplink bytes.
func (s *SessionStats) Uplink() int64 {
	return atomic.LoadInt64(&s.uplink)
}

// Downlink returns the total downlink bytes.
func (s *SessionStats) Downlink() int64 {
	return atomic.LoadInt64(&s.downlink)
}

// TotalTraffic returns the sum of uplink and downlink bytes.
func (s *SessionStats) TotalTraffic() int64 {
	return atomic.LoadInt64(&s.uplink) + atomic.LoadInt64(&s.downlink)
}

// Duration returns the session duration. If the session is still open,
// it returns the elapsed time since start. Note: for very short-lived sessions
// this may return a near-zero duration, which is expected and not an error.
func (s *SessionStats) Duration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.end.IsZero() {
		return s.end.Sub(s.start)
	}
	return time.Since(s.start)
}

// Close marks the session as closed and records the end time.
// It is safe to call Close multiple times; only the first call has effect.
func (s *SessionStats) Close() {
	if atomic.CompareAndSwapUint32(&s.closed, 0, 1) {
		s.mu.Lock()
		s.end = time.Now()
		s.mu.Unlock()
	}
}

// IsClosed returns true if the session has been closed.
func (s *SessionStats) IsClosed() bool {
	return atomic.LoadUint32(&s.closed) == 1
}

// StartTime returns the time the session started.
func (s *SessionStats) StartTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.start
}

// EndTime returns the time the session ended, or zero if still open.
func (s *SessionStats) EndTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.end
}

// AverageUplinkRate returns the average uplink rate in bytes per second.
// Returns 0 if the session duration is zero.
func (s *SessionStats) AverageUplinkRate() float64 {
	d := s.Duration().Seconds()
	if d <= 0 {
		return 0
	}
	return float64(atomic.LoadInt64(&s.uplink)) / d
}

// AverageDownlinkRate returns the average downlink rate in bytes per second.
// Returns 0 if the session duration is zero.
func (s *SessionStats) AverageDownlinkRate() float64 {
	d := s.Duration().Seconds()
	if d <= 0 {
		return 0
	}
	return float64(atomic.LoadInt64(&s.downlink)) / d
}
