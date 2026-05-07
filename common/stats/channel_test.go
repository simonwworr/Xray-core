package stats_test

import (
	"context"
	"testing"
	"time"

	"github.com/xtls/xray-core/common/stats"
)

func TestChannelSubscribePublish(t *testing.T) {
	ch := stats.NewChannel()
	sub, err := ch.Subscribe()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	ch.Publish(ctx, "hello")
	select {
	case msg := <-sub:
		if msg != "hello" {
			t.Errorf("expected 'hello', got %v", msg)
		}
	case <-time.After(2 * time.Second):
		// Increased timeout from 1s to 2s to reduce flakiness on slow CI machines
		t.Error("timeout waiting for message")
	}
}

func TestChannelUnsubscribe(t *testing.T) {
	ch := stats.NewChannel()
	sub, err := ch.Subscribe()
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Unsubscribe(sub); err != nil {
		t.Fatal(err)
	}
	// Channel should be closed after unsubscribe
	select {
	case _, ok := <-sub:
		if ok {
			t.Error("expected channel to be closed")
		}
	default:
		t.Error("expected closed channel to be readable")
	}
}

func TestChannelClose(t *testing.T) {
	ch := stats.NewChannel()
	_, err := ch.Subscribe()
	if err != nil {
		t.Fatal(err)
	}
	if err := ch.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = ch.Subscribe()
	if err == nil {
		t.Error("expected error subscribing to closed channel")
	}
}
