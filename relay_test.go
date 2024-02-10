package khatru_test

import (
	"log"
	"testing"
	"time"

	"github.com/nbd-wtf/go-nostr/nip11"
	"github.com/stretchr/testify/assert"

	"github.com/fiatjaf/khatru"
)

func TestNewRelay(t *testing.T) {
	t.Run("constructor returns a Relay", func(t *testing.T) {
		relay := khatru.NewRelay()
		assert.NotNil(t, relay)
	})

	t.Run("foo", func(t *testing.T) {
		relay := khatru.NewRelay()

		var document *nip11.RelayInformationDocument
		var logger *log.Logger

		assert.Equal(t, "", relay.ServiceURL)
		assert.IsType(t, document, relay.Info)
		assert.IsType(t, logger, relay.Log)
		assert.Equal(t, relay.WriteWait, 10*time.Second)
		assert.Equal(t, relay.PongWait, 60*time.Second)
		assert.Equal(t, relay.PingPeriod, 30*time.Second)
		assert.Equal(t, relay.MaxMessageSize, int64(512000))
	})
}
