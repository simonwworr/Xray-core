package stats

import (
	"context"
	"sync"
)

// Channel is a thread-safe channel for broadcasting stats messages.
type Channel struct {
	mu          sync.RWMutex
	subscribers []chan interface{}
	closed      bool
}

// NewChannel creates a new Channel.
func NewChannel() *Channel {
	return &Channel{}
}

// Subscribe registers a new subscriber and returns a receive-only channel.
// The buffer size of 64 helps avoid dropping messages under higher load.
// Increased from 16 to 64 to reduce message drops on busy connections.
func (c *Channel) Subscribe() (<-chan interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil, newError("channel is closed")
	}
	ch := make(chan interface{}, 64)
	c.subscribers = append(c.subscribers, ch)
	return ch, nil
}

// Unsubscribe removes a subscriber channel.
func (c *Channel) Unsubscribe(sub <-chan interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, s := range c.subscribers {
		if s == sub {
			c.subscribers = append(c.subscribers[:i], c.subscribers[i+1:]...)
			close(s)
			return nil
		}
	}
	return newError("subscriber not found")
}

// Publish sends a message to all subscribers.
// Note: messages are dropped for a subscriber if its buffer is full (non-blocking send).
// The context can be used to cancel broadcasting mid-way if the caller is done.
func (c *Channel) Publish(ctx context.Context, msg interface{}) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, sub := range c.subscribers {
		select {
		case sub <- msg:
		case <-ctx.Done():
			return
		default:
			// subscriber buffer full, drop message rather than block
		}
	}
}

// Close closes the channel and all subscribers.
func (c *Channel) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.closed = true
		for _, sub := range c.subscribers {
			close(sub)
		}
		c.subscribers = nil
	}
	return nil
}
