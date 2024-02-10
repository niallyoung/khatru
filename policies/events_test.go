package policies_test

import (
	"context"
	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/fiatjaf/khatru/policies"
)

func TestPreventTimestampsInTheFuture(t *testing.T) {
	t.Run("Now() returns false with no msg", func(t *testing.T) {
		fn := policies.PreventTimestampsInTheFuture(nostr.Now())

		ok, result := fn(context.Background(), &nostr.Event{})
		assert.Equal(t, false, ok)
		assert.Equal(t, "", result)
	})

	t.Run("future event below threshold returns false", func(t *testing.T) {
		const thresholdSeconds int64 = 120
		thresholdTimestamp := nostr.Timestamp(thresholdSeconds)
		futureTimestamp := time.Now().Add(60 * time.Second).Unix()
		futureEvent := &nostr.Event{CreatedAt: nostr.Timestamp(futureTimestamp)}

		fn := policies.PreventTimestampsInTheFuture(thresholdTimestamp)
		ok, result := fn(context.Background(), futureEvent)

		assert.Equal(t, false, ok)
		assert.Equal(t, "", result)
	})

	t.Run("future event above threshold returns false", func(t *testing.T) {
		const thresholdSeconds int64 = 120
		thresholdTimestamp := nostr.Timestamp(thresholdSeconds)
		futureTimestamp := time.Now().Add(240 * time.Second).Unix()
		futureEvent := &nostr.Event{CreatedAt: nostr.Timestamp(futureTimestamp)}

		fn := policies.PreventTimestampsInTheFuture(thresholdTimestamp)
		ok, result := fn(context.Background(), futureEvent)

		assert.Equal(t, true, ok)
		assert.Equal(t, "event too much in the future", result)
	})

	t.Run("past events always return false", func(t *testing.T) {
		const thresholdSeconds int64 = 120
		negativeInts := []int64{-30, -60, -200, -2000, -9999, -1}

		for _, s := range negativeInts {
			thresholdTimestamp := nostr.Timestamp(thresholdSeconds)
			futureTimestamp := time.Now().Add(time.Duration(s) * time.Second).Unix()
			futureEvent := &nostr.Event{CreatedAt: nostr.Timestamp(futureTimestamp)}

			fn := policies.PreventTimestampsInTheFuture(thresholdTimestamp)
			ok, result := fn(context.Background(), futureEvent)

			assert.Equal(t, false, ok)
			assert.Equal(t, "", result)
		}
	})
}
