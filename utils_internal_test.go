package khatru

import (
	"context"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func TestGetOpenSubscriptions(t *testing.T) {
	t.Run("GetOpenSubscriptions() returns expected filters", func(t *testing.T) {
		ws := &WebSocket{}
		setListener("id", ws, nostr.Filters{}, nil)
		ctx := context.WithValue(context.Background(), wsKey, ws)

		subs := GetOpenSubscriptions(ctx)
		assert.NotNil(t, subs)
	})
}
