package khatru_test

import (
	"github.com/fiatjaf/khatru"
	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRelay_BroadcastEvent(t *testing.T) {
	relay := khatru.NewRelay()
	event := &nostr.Event{}

	assert.NotPanics(t, func() { relay.BroadcastEvent(event) })
}
