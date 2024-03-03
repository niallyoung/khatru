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

		err := relay.handleRequest(context.Background(), "id", wg, &WebSocket{}, nostr.Filter{Limit: -1})
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
}
