package policies_test

import (
	"context"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"

	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/policies"
)

const wsKey = iota

func TestRejectKind04Snoopers(t *testing.T) {
	t.Run("returns false if 4 not in Kind", func(t *testing.T) {
		filter := nostr.Filter{Kinds: []int{2}}

		rejected, msg := policies.RejectKind04Snoopers(context.Background(), filter)
		assert.Equal(t, false, rejected)
		assert.Equal(t, msg, "")
	})

	t.Run("returns true when not authenticated", func(t *testing.T) {
		filter := nostr.Filter{Kinds: []int{4}}
		ws := &khatru.WebSocket{}
		ctx := context.WithValue(context.Background(), wsKey, ws)

		rejected, msg := policies.RejectKind04Snoopers(ctx, filter)
		assert.Equal(t, true, rejected)
		assert.Equal(t, "restricted: this relay does not serve kind-4 to unauthenticated users, does your client implement NIP-42?", msg)
	})
}
