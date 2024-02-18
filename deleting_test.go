package khatru

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func TestRelay_handleDeleteRequest(t *testing.T) {
	scenarios := []struct {
		Description     string
		Event           *nostr.Event
		ExpectedError   error
		ExpectedBlocked bool
	}{
		{
			"empty event",
			&nostr.Event{},
			nil,
			false,
		},
		{
			"event with 'e' tag",
			&nostr.Event{
				PubKey: "fake pubkey",
				Tags: nostr.Tags{
					[]string{"e", "v"},
					[]string{"f", "v"},
				},
			},
			nil,
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("returns the expected error for %s when no QueryEvents", s.Description), func(t *testing.T) {
			relay := NewRelay()

			err := relay.handleDeleteRequest(context.Background(), s.Event)
			assert.Equal(t, s.ExpectedError, err)
		})
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("returns no error for %s when QueryEvents has an event and no error", s.Description), func(t *testing.T) {
			relay := NewRelay()
			relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
				ch := make(chan *nostr.Event)
				close(ch)
				return ch, nil
			})

			err := relay.handleDeleteRequest(context.Background(), s.Event)
			assert.NoError(t, err)
		})
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("returns no error for %s when QueryEvents returns an error", s.Description), func(t *testing.T) {
			relay := NewRelay()
			relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
				return nil, errors.New("fake QueryEvents error")
			})

			err := relay.handleDeleteRequest(context.Background(), s.Event)
			assert.NoError(t, err)
		})
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("successfully handles deletion requests of '%s' when OverwriteDeletionOutcome returns true", s.Description), func(t *testing.T) {
			relay := NewRelay()
			relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
				ch := make(chan *nostr.Event)
				go func() {
					previous := &nostr.Event{PubKey: "fake-pubkey"}
					ch <- previous
					close(ch)
				}()
				return ch, nil
			})

			relay.OverwriteDeletionOutcome = append(relay.OverwriteDeletionOutcome,
				func(ctx context.Context, target *nostr.Event, deletion *nostr.Event) (bool, string) {
					return true, ""
				})

			relay.DeleteEvent = append(relay.DeleteEvent, func(ctx context.Context, event *nostr.Event) error {
				return nil
			})

			err := relay.handleDeleteRequest(context.Background(), s.Event)
			assert.NoError(t, err)
		})
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("TBC deletion request of '%s' when OverwriteDeletionOutcome returns false", s.Description), func(t *testing.T) {
			relay := NewRelay()
			relay.QueryEvents = append(relay.QueryEvents, func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
				ch := make(chan *nostr.Event)
				go func() {
					previous := &nostr.Event{PubKey: "fake-pubkey"}
					ch <- previous
					close(ch)
				}()
				return ch, nil
			})

			relay.OverwriteDeletionOutcome = append(relay.OverwriteDeletionOutcome,
				func(ctx context.Context, target *nostr.Event, deletion *nostr.Event) (bool, string) {
					return false, ""
				})

			relay.DeleteEvent = append(relay.DeleteEvent, func(ctx context.Context, event *nostr.Event) error {
				return nil
			})

			err := relay.handleDeleteRequest(context.Background(), s.Event)
			if s.ExpectedBlocked {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "blocked: ")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
