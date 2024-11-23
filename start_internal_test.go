package khatru

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelay_Router(t *testing.T) {
	t.Run("Router() returns nil when no ", func(t *testing.T) {
		relay := Relay{}
		router := relay.Router()
		assert.Nil(t, router)
	})

	t.Run("Router() returns nil when no ", func(t *testing.T) {
		relay := NewRelay()
		router := relay.Router()
		assert.NotNil(t, router)
	})
}

func TestRelay_Start(t *testing.T) {
	t.Run("Start() returns nil on clean Start;Shutdown", func(t *testing.T) {
		relay := NewRelay()
		timer := time.AfterFunc(time.Second*1, func() {
			relay.Shutdown(context.Background())
		})
		defer timer.Stop()

		err := relay.Start("localhost", 1234)
		assert.NoError(t, err)
	})

	t.Run("Start() closes started channels", func(t *testing.T) {
		relay := NewRelay()
		timer := time.AfterFunc(time.Second*1, func() {
			relay.Shutdown(context.Background())
		})
		defer timer.Stop()

		started := make(chan bool)
		err := relay.Start("localhost", 1234, started)
		assert.NoError(t, err)

		_, channelOpen := <-started
		assert.False(t, channelOpen)
	})
}

func TestRelay_Shutdown(t *testing.T) {
	t.Run("Shutdown() .....?", func(t *testing.T) {
		relay := NewRelay()
		timer := time.AfterFunc(time.Second*1, func() {
			relay.Shutdown(context.Background())
		})
		defer timer.Stop()

		err := relay.Start("localhost", 1234)
		assert.NoError(t, err)
	})
}
