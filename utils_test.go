package khatru_test

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fiatjaf/khatru"
)

const (
	wsKey = iota
	subscriptionIdKey
)

func TestRequestAuth(t *testing.T) {
	t.Run("request auth with valid context", func(t *testing.T) {
		ws := &khatru.WebSocket{}
		ctx := context.WithValue(context.Background(), wsKey, ws)

		assert.Nil(t, ws.Authed)
		assert.NotPanics(t, func() { khatru.RequestAuth(ctx) })
		assert.NotNil(t, ws.Authed)
	})
}

func TestGetAuthed(t *testing.T) {
	t.Run("GetAuthed() without valid context panics", func(t *testing.T) {
		assert.Panics(t, func() { khatru.GetAuthed(context.Background()) })
	})

	t.Run("GetAuthed() with valid context doesn't panic", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), wsKey, &khatru.WebSocket{})
		assert.NotPanics(t, func() { khatru.GetAuthed(ctx) })
	})
}

func TestGetIP(t *testing.T) {
	t.Run("GetIP() returns the expected IP from headers", func(t *testing.T) {
		req := httptest.NewRequest("method", "http://foo", bytes.NewReader([]byte{}))
		req.Header.Add("X-Forwarded-For", "1.2.3.4")
		ws := &khatru.WebSocket{Request: req}
		ctx := context.WithValue(context.Background(), wsKey, ws)

		ip := khatru.GetIP(ctx)
		assert.Equal(t, "1.2.3.4:1234", ip)
	})
}

func TestGetSubscriptionID(t *testing.T) {
	t.Run("GetSubscriptionID() returns the expected id from context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), subscriptionIdKey, "fake-subscription-id")

		id := khatru.GetSubscriptionID(ctx)
		assert.Equal(t, "fake-subscription-id", id)
	})
}

func TestGetOpenSubscriptions(t *testing.T) {
	t.Run("GetOpenSubscriptions() returns nil when no subscriptions", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), wsKey, &khatru.WebSocket{})

		subs := khatru.GetOpenSubscriptions(ctx)
		assert.Nil(t, subs)
	})

	//t.Run("GetOpenSubscriptions() returns expected filters", func(t *testing.T) {
	//	ctx := context.WithValue(context.Background(), wsKey, &khatru.WebSocket{})
	//
	//	subs := khatru.GetOpenSubscriptions(ctx)
	//	assert.NotNil(t, subs)
	//})
}
