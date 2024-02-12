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

	t.Run("20000 <= event.Kind < 30000", func(t *testing.T) {
		t.Run("", func(t *testing.T) {
			scenarios := []struct {
				Name string
				Kind int
			}{
				{Name: "<", Kind: 19999},
				{Name: "min", Kind: 20000},
				{Name: "mid", Kind: 20000},
				{Name: "max", Kind: 30000},
				{Name: ".", Kind: 30001},
			}

			for _, s := range scenarios {
				t.Run(s.Name, func(t *testing.T) {
					relay := khatru.NewRelay()
					relay.OnEphemeralEvent = append(relay.OnEphemeralEvent,
						func(ctx context.Context, event *nostr.Event) { return },
					)

					err := relay.AddEvent(context.Background(), &nostr.Event{Kind: s.Kind})
					assert.NoError(t, err)
				})
			}
		})
	})
}
