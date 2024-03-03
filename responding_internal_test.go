package khatru

import (
	"context"
	"sync"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func TestRelay_handleRequest(t *testing.T) {
	t.Run("handleRequest() overwrites filter OverwriteFilter", func(t *testing.T) {
		relay := NewRelay()
		ovwFilterCalled := false
		relay.OverwriteFilter = append(relay.OverwriteFilter,
			func(ctx context.Context, filter *nostr.Filter) { ovwFilterCalled = true })
		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := relay.handleRequest(context.Background(), "id", wg, &WebSocket{}, nostr.Filter{})
		assert.Nil(t, err)
		assert.True(t, ovwFilterCalled)
	})

	t.Run("handleRequest() returns nil when filter.Limit < 0", func(t *testing.T) {
		relay := NewRelay()
		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := relay.handleRequest(context.Background(), "id", wg, &WebSocket{}, nostr.Filter{Limit: -1})
		assert.Nil(t, err)
	})

	t.Run("handleRequest() with a rejecting RejectFilter returns an error", func(t *testing.T) {
		relay := NewRelay()
		rejectFilterCalled := false
		relay.RejectFilter = append(relay.RejectFilter,
			func(ctx context.Context, filter nostr.Filter) (bool, string) {
				rejectFilterCalled = true
				return true, ""
			})
		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := relay.handleRequest(context.Background(), "id", wg, &WebSocket{}, nostr.Filter{})
		assert.ErrorContains(t, err, "blocked: ")
		assert.True(t, rejectFilterCalled)
	})
}
