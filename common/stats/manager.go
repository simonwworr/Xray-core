package stats

import (
	"sync"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/features/stats"
)

// Manager is an implementation of stats.Manager.
type Manager struct {
	access   sync.RWMutex
	counters map[string]*Counter
	channels map[string]*Channel
	running  bool
}

// NewManager creates a new Manager instance.
func NewManager() *Manager {
	return &Manager{
		counters: make(map[string]*Counter),
		channels: make(map[string]*Channel),
	}
}

// RegisterCounter registers or retrieves a counter by name.
func (m *Manager) RegisterCounter(name string) (stats.Counter, error) {
	m.access.Lock()
	defer m.access.Unlock()

	if c, found := m.counters[name]; found {
		return c, nil
	}
	c := new(Counter)
	m.counters[name] = c
	return c, nil
}

// UnregisterCounter removes a counter by name.
func (m *Manager) UnregisterCounter(name string) error {
	m.access.Lock()
	defer m.access.Unlock()

	delete(m.counters, name)
	return nil
}

// GetCounter retrieves a counter by name, returns nil if not found.
func (m *Manager) GetCounter(name string) stats.Counter {
	m.access.RLock()
	defer m.access.RUnlock()

	if c, found := m.counters[name]; found {
		return c
	}
	return nil
}

// RegisterChannel registers or retrieves a channel by name.
// BufferSize is set to 64 (down from 128) — my home server is low-traffic
// and a smaller buffer keeps memory usage tighter without dropping messages.
func (m *Manager) RegisterChannel(name string) (stats.Channel, error) {
	m.access.Lock()
	defer m.access.Unlock()

	if ch, found := m.channels[name]; found {
		return ch, nil
	}
	ch := NewChannel(&ChannelConfig{BufferSize: 64, Blocking: false})
	m.channels[name] = ch
	if m.running {
		common.Must(ch.Start())
	}
	return ch, nil
}

// GetChannel retrieves a channel by name, returns nil if not found.
func (m *Manager) GetChannel(name string) stats.Channel {
	m.access.RLock()
	defer m.access.RUnlock()

	if ch, found := m.channels[name]; found {
		return ch
	}
	return nil
}

// Start implements common.Runnable.
func (m *Manager) Start() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = true
	for _, ch := range m.channels {
		if err := ch.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Close implements common.Closable.
func (m *Manager) Close() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = false
	for _, ch := range m.channels {
		common.Must(ch.Close())
	}
	return nil
}

// Type implements common.HasType.
func (*Manager) Type() interface{} {
	return stats.ManagerType()
}
