package khatru_test

import (
	"context"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"

	"github.com/fiatjaf/khatru"
)

func TestRelay_AddEvent(t *testing.T) {

	t.Run(".RejectEvent", func(t *testing.T) {
		t.Run("returns error when event is nil", func(t *testing.T) {
			relay := khatru.NewRelay()
			err := relay.AddEvent(context.Background(), nil)

			assert.Error(t, err)
		})

		t.Run("returns blocked-details error when event rejected with a msg", func(t *testing.T) {
			relay := khatru.NewRelay()
			relay.RejectEvent = append(relay.RejectEvent,
				func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
					return true, "reasons"
				},
			)

			err := relay.AddEvent(context.Background(), &nostr.Event{})

			assert.Error(t, err)
			assert.ErrorContains(t, err, "blocked: reasons")
		})

		t.Run("returns error when event rejected without a msg", func(t *testing.T) {
			relay := khatru.NewRelay()
			relay.RejectEvent = append(relay.RejectEvent,
				func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
					return true, ""
				},
			)

			err := relay.AddEvent(context.Background(), &nostr.Event{})

			assert.Error(t, err)
			assert.ErrorContains(t, err, "blocked: no reason")
		})
	})
}
