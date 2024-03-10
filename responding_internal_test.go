package khatru

import (
	"context"
	"errors"
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

	t.Run("handleRequest() executes QueryEvents with no error", func(t *testing.T) {
		relay := NewRelay()
		queryEventsCalled := false
		relay.QueryEvents = append(relay.QueryEvents,
			func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
				queryEventsCalled = true
				return nil, nil
			})
		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := relay.handleRequest(context.Background(), "id", wg, &WebSocket{}, nostr.Filter{})
		assert.NoError(t, err)
		assert.True(t, queryEventsCalled)
	})

	t.Run("handleRequest() executes QueryEvents with an error", func(t *testing.T) {
		relay := NewRelay()
		queryEventsCalled := false
		relay.QueryEvents = append(relay.QueryEvents,
			func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
				queryEventsCalled = true
				return nil, errors.New("fake-query-error")
			})
		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := relay.handleRequest(context.Background(), "id", wg, &WebSocket{}, nostr.Filter{})
		assert.NoError(t, err)
		assert.True(t, queryEventsCalled)
	})

	//t.Run("handleRequest() executes OverwriteResponseEvent", func(t *testing.T) {
	//	relay := NewRelay()
	//	queryEventsCalled := false
	//	overwriteResponseEventsCalled := false
	//	relay.QueryEvents = append(relay.QueryEvents,
	//		func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
	//			queryEventsCalled = true
	//			eventChan := make(chan *nostr.Event)
	//			eventChan <- &nostr.Event{ID: "id"}
	//			return eventChan, nil
	//		})
	//	relay.OverwriteResponseEvent = append(relay.OverwriteResponseEvent,
	//		func(ctx context.Context, event *nostr.Event) {
	//			overwriteResponseEventsCalled = true
	//			return
	//		})
	//	wg := &sync.WaitGroup{}
	//	wg.Add(1)
	//
	//	err := relay.handleRequest(context.Background(), "id", wg, &WebSocket{}, nostr.Filter{})
	//	assert.NoError(t, err)
	//	assert.True(t, queryEventsCalled)
	//	assert.True(t, overwriteResponseEventsCalled)
	//})
}
