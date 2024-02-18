package khatru_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/fiatjaf/eventstore"
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
			kinds := []int{20000, 25000, 29999}
			for _, k := range kinds {
				t.Run(fmt.Sprint(k), func(t *testing.T) {
					relay := khatru.NewRelay()
					relay.OnEphemeralEvent = append(relay.OnEphemeralEvent,
						func(ctx context.Context, event *nostr.Event) { return },
					)

					err := relay.AddEvent(context.Background(), &nostr.Event{Kind: k})
					assert.NoError(t, err)
				})
			}
		})
	})

	t.Run("event.Kind 0, 3, 10000 <= event.Kind < 20000", func(t *testing.T) {
		kinds := []int{0, 3, 10000, 19999}

		t.Run("query returns an error", func(t *testing.T) {
			for _, k := range kinds {
				t.Run(fmt.Sprint(k), func(t *testing.T) {
					relay := khatru.NewRelay()
					relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
						return nil, errors.New("fake QueryEvents error")
					})

					err := relay.AddEvent(context.Background(), &nostr.Event{Kind: k})
					assert.NoError(t, err)
				})
			}
		})

		t.Run("replaceable events are deleted before storing", func(t *testing.T) {
			for _, k := range kinds {
				t.Run(fmt.Sprint(k), func(t *testing.T) {
					relay := khatru.NewRelay()
					relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
						ch := make(chan *nostr.Event)
						go func() {
							previous := &nostr.Event{CreatedAt: 0}
							ch <- previous
							close(ch)
						}()
						return ch, nil
					})

					deleteEventCalled := false
					relay.DeleteEvent = append(relay.DeleteEvent,
						func(ctx context.Context, event *nostr.Event) error { deleteEventCalled = true; return nil },
					)

					err := relay.AddEvent(context.Background(), &nostr.Event{Kind: k, CreatedAt: 1})
					assert.NoError(t, err)
					assert.Equal(t, true, deleteEventCalled)
				})
			}
		})

	})

	t.Run("30000 <= event.Kind < 40000", func(t *testing.T) {
		kinds := []int{30000, 35000, 39999}

		t.Run("QueryEvents returns no error, when QueryEvents returns error", func(t *testing.T) {
			for _, k := range kinds {

				t.Run(fmt.Sprintf("for events with kind %d", k), func(t *testing.T) {
					relay := khatru.NewRelay()
					relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
						return nil, errors.New("fake query error")
					})

					err := relay.AddEvent(context.Background(), &nostr.Event{
						Kind:   k,
						PubKey: "fake-pubkey",
						Tags: nostr.Tags{
							[]string{"d", "v"},
						},
					})
					assert.NoError(t, err)
				})
			}
		})

		t.Run("QueryEvents deletes events when older duplicates are found", func(t *testing.T) {
			for _, k := range kinds {
				t.Run(fmt.Sprintf("for events with kind %d", k), func(t *testing.T) {
					relay := khatru.NewRelay()
					relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
						ch := make(chan *nostr.Event)
						go func() {
							previous := &nostr.Event{
								Kind:      k,
								PubKey:    "fake-pubkey",
								CreatedAt: 0,
								Tags: nostr.Tags{
									[]string{"d", "v"},
								},
							}
							ch <- previous
							close(ch)
						}()
						return ch, nil
					})

					deleteEventCalled := false
					relay.DeleteEvent = append(relay.DeleteEvent,
						func(ctx context.Context, event *nostr.Event) error { deleteEventCalled = true; return nil },
					)

					err := relay.AddEvent(context.Background(), &nostr.Event{
						Kind:      k,
						PubKey:    "fake-pubkey",
						CreatedAt: 1,
						Tags: nostr.Tags{
							[]string{"d", "v"},
						},
					})
					assert.NoError(t, err)
					assert.Equal(t, true, deleteEventCalled)
				})
			}
		})
	})

	t.Run("StoreEvent is called", func(t *testing.T) {
		kinds := []int{1, 2, 9999999}

		for _, k := range kinds {
			t.Run(fmt.Sprintf("for events with kind %d", k), func(t *testing.T) {

				t.Run("AddEvent() calls StoreEvent and OnEventSaved", func(t *testing.T) {
					relay := khatru.NewRelay()
					relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
						return nil, nil
					})

					storeEventCalled := false
					relay.StoreEvent = append(relay.StoreEvent, func(ctx context.Context, event *nostr.Event) error {
						storeEventCalled = true
						return nil
					})

					onEventSavedCalled := false
					relay.OnEventSaved = append(relay.OnEventSaved, func(ctx context.Context, event *nostr.Event) {
						onEventSavedCalled = true
						return
					})

					err := relay.AddEvent(context.Background(), &nostr.Event{
						Kind:      k,
						PubKey:    "fake-pubkey",
						CreatedAt: 0,
						Tags: nostr.Tags{
							[]string{"d", "v"},
						},
					})
					assert.NoError(t, err)
					assert.Equal(t, true, storeEventCalled)
					assert.Equal(t, true, onEventSavedCalled)
				})

				t.Run("AddEvent() handles duplicates without a store error", func(t *testing.T) {
					relay := khatru.NewRelay()
					relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
						return nil, nil
					})

					storeEventCalled := false
					relay.StoreEvent = append(relay.StoreEvent, func(ctx context.Context, event *nostr.Event) error {
						storeEventCalled = true
						return eventstore.ErrDupEvent
					})

					err := relay.AddEvent(context.Background(), &nostr.Event{
						Kind:      k,
						PubKey:    "fake-pubkey",
						CreatedAt: 0,
						Tags: nostr.Tags{
							[]string{"d", "v"},
						},
					})
					assert.NoError(t, err)
					assert.Equal(t, true, storeEventCalled)
				})

				t.Run("AddEvent() handles other store errors", func(t *testing.T) {
					relay := khatru.NewRelay()
					relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
						return nil, nil
					})

					storeEventCalled := false
					relay.StoreEvent = append(relay.StoreEvent, func(ctx context.Context, event *nostr.Event) error {
						storeEventCalled = true
						return errors.New("fake store error")
					})

					err := relay.AddEvent(context.Background(), &nostr.Event{
						Kind:      k,
						PubKey:    "fake-pubkey",
						CreatedAt: 1,
						Tags: nostr.Tags{
							[]string{"d", "v"},
						},
					})
					assert.Error(t, err)
					assert.ErrorContains(t, err, "fake store error")
					assert.Equal(t, true, storeEventCalled)
				})
			})
		}
	})
}
