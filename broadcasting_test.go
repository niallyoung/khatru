package khatru_test

import (
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"

	"github.com/fiatjaf/khatru"
)

func TestRelay_BroadcastEvent(t *testing.T) {
	relay := khatru.NewRelay()
	event := &nostr.Event{}

	assert.NotPanics(t, func() { relay.BroadcastEvent(event) })
}
