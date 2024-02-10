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

		var expectedDocument *nip11.RelayInformationDocument
		var expectedLogger *log.Logger

		assert.Equal(t, "", relay.ServiceURL)
		assert.IsType(t, expectedDocument, relay.Info)
		assert.IsType(t, expectedLogger, relay.Log)
		assert.Equal(t, relay.WriteWait, 10*time.Second)
		assert.Equal(t, relay.PongWait, 60*time.Second)
		assert.Equal(t, relay.PingPeriod, 30*time.Second)
		assert.Equal(t, relay.MaxMessageSize, int64(512000))
	})
}
