package policies_test

import (
	"context"
	"testing"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"

	"github.com/fiatjaf/khatru/policies"
)

func TestPreventTimestampsInTheFuture(t *testing.T) {
	t.Run("Now() returns false with no msg", func(t *testing.T) {
		fn := policies.PreventTimestampsInTheFuture(time.Duration(0))

		ok, result := fn(context.Background(), &nostr.Event{})
		assert.Equal(t, false, ok)
		assert.Equal(t, "", result)
	})

	t.Run("future event below threshold returns false", func(t *testing.T) {
		const thresholdSeconds = 120
		now := time.Now()
		futureTimestamp := now.Add(60 * time.Second).Unix()
		futureEvent := &nostr.Event{CreatedAt: nostr.Timestamp(futureTimestamp)}

		threshold := time.Duration(thresholdSeconds) * time.Second
		fn := policies.PreventTimestampsInTheFuture(threshold)
		ok, result := fn(context.Background(), futureEvent)

		assert.Equal(t, false, ok)
		assert.Equal(t, "", result)
	})

	t.Run("future event above threshold returns true", func(t *testing.T) {
		const thresholdSeconds = 120
		futureTimestamp := time.Now().Add(240 * time.Second).Unix()
		futureEvent := &nostr.Event{CreatedAt: nostr.Timestamp(futureTimestamp)}

		fn := policies.PreventTimestampsInTheFuture(time.Duration(thresholdSeconds) * time.Second)
		ok, result := fn(context.Background(), futureEvent)

		assert.Equal(t, true, ok)
		assert.Equal(t, "event too much in the future", result)
	})

	t.Run("past events always return false", func(t *testing.T) {
		const thresholdSeconds = 120
		negativeInts := []int64{-30, -60, -200, -2000, -9999, -1}

		for _, s := range negativeInts {
			futureTimestamp := time.Now().Add(time.Duration(s) * time.Second).Unix()
			futureEvent := &nostr.Event{CreatedAt: nostr.Timestamp(futureTimestamp)}

			fn := policies.PreventTimestampsInTheFuture(time.Duration(thresholdSeconds) * time.Second)
			ok, result := fn(context.Background(), futureEvent)

			assert.Equal(t, false, ok)
			assert.Equal(t, "", result)
		}
	})
}

func TestPreventTimestampsInThePast(t *testing.T) {
	t.Run("past event above threshold returns true", func(t *testing.T) {
		fn := policies.PreventTimestampsInThePast(time.Duration(0))

		ok, result := fn(context.Background(), &nostr.Event{})
		assert.Equal(t, true, ok)
		assert.Equal(t, "event too old", result)
	})

	t.Run("past event below threshold returns false", func(t *testing.T) {
		const thresholdSeconds = 120
		pastTimestamp := time.Now().Add(-60 * time.Second).Unix()
		pastEvent := &nostr.Event{CreatedAt: nostr.Timestamp(pastTimestamp)}

		fn := policies.PreventTimestampsInThePast(time.Duration(thresholdSeconds) * time.Second)
		ok, result := fn(context.Background(), pastEvent)

		assert.Equal(t, false, ok)
		assert.Equal(t, "", result)
	})

	t.Run("past event above threshold returns true", func(t *testing.T) {
		const thresholdSeconds = 120
		pastTimestamp := time.Now().Add(-240 * time.Second).Unix()
		pastEvent := &nostr.Event{CreatedAt: nostr.Timestamp(pastTimestamp)}

		fn := policies.PreventTimestampsInThePast(time.Duration(thresholdSeconds) * time.Second)
		ok, result := fn(context.Background(), pastEvent)

		assert.Equal(t, true, ok)
		assert.Equal(t, "event too old", result)
	})

	t.Run("future events always return false", func(t *testing.T) {
		const thresholdSeconds = 120
		negativeInts := []int64{0, 5, 20, 100, 1000, 9999}

		for _, s := range negativeInts {
			futureTimestamp := time.Now().Add(time.Duration(s) * time.Second).Unix()
			futureEvent := &nostr.Event{CreatedAt: nostr.Timestamp(futureTimestamp)}

			fn := policies.PreventTimestampsInThePast(time.Duration(thresholdSeconds) * time.Second)
			ok, result := fn(context.Background(), futureEvent)

			assert.Equal(t, false, ok)
			assert.Equal(t, "", result)
		}
	})
}

func TestRestrictToSpecifiedKinds(t *testing.T) {
	t.Run("returns false when event.Kind in slice", func(t *testing.T) {
		variousInts := []uint16{0, 1, 20, 9999}

		event := &nostr.Event{Kind: 0}
		fn := policies.RestrictToSpecifiedKinds(false, variousInts...)

		ok, result := fn(context.Background(), event)

		assert.Equal(t, false, ok)
		assert.Equal(t, "", result)
	})

	t.Run("returns true when event.Kind < min", func(t *testing.T) {
		variousInts := []uint16{3, 1, 5, 20, 9999}

		event := &nostr.Event{Kind: 0}
		fn := policies.RestrictToSpecifiedKinds(false, variousInts...)

		ok, result := fn(context.Background(), event)

		assert.Equal(t, true, ok)
		assert.Equal(t, "received event kind 0 not allowed", result)
	})

	t.Run("returns true when event.Kind > max", func(t *testing.T) {
		variousInts := []uint16{1, 5, 20}

		event := &nostr.Event{Kind: 999}
		fn := policies.RestrictToSpecifiedKinds(false, variousInts...)

		ok, result := fn(context.Background(), event)

		assert.Equal(t, true, ok)
		assert.Equal(t, "received event kind 999 not allowed", result)
	})

	t.Run("returns false when event.Kind not in variousInts", func(t *testing.T) {
		variousInts := []uint16{1, 5, 20}

		event := &nostr.Event{Kind: 6}
		fn := policies.RestrictToSpecifiedKinds(false, variousInts...)

		ok, result := fn(context.Background(), event)

		assert.Equal(t, true, ok)
		assert.Equal(t, "received event kind 6 not allowed", result)
	})
}

func TestPreventLargeTags(t *testing.T) {
	t.Run("returns false when a tags is <= maxTagValueLen", func(t *testing.T) {
		event := &nostr.Event{
			Tags: nostr.Tags{
				nostr.Tag{"k", "v"},
				nostr.Tag{"l", "m"},
				nostr.Tag{"m", "n"},
			},
		}

		fn := policies.PreventLargeTags(1)
		ok, msg := fn(context.Background(), event)

		assert.Equal(t, false, ok)
		assert.Equal(t, "", msg)
	})

	t.Run("returns true when a tags is > maxTagValueLen", func(t *testing.T) {
		event := &nostr.Event{
			Tags: nostr.Tags{
				nostr.Tag{"n", "ooooooo"},
			},
		}

		fn := policies.PreventLargeTags(1)
		ok, msg := fn(context.Background(), event)

		assert.Equal(t, true, ok)
		assert.Equal(t, "event contains too large tags", msg)
	})
}

func TestPreventTooManyIndexableTags(t *testing.T) {
	t.Run("returns false when event.Kind in ignoreKinds", func(t *testing.T) {
		event := &nostr.Event{Kind: 0}
		ignoreKinds := []int{0}
		onlyKinds := []int{}

		fn := policies.PreventTooManyIndexableTags(3, ignoreKinds, onlyKinds)
		ok, msg := fn(context.Background(), event)

		assert.Equal(t, false, ok)
		assert.Equal(t, "", msg)
	})

	t.Run("returns false when event.Kind in onlyKinds", func(t *testing.T) {
		event := &nostr.Event{Kind: 0}
		ignoreKinds := []int{}
		onlyKinds := []int{0}

		fn := policies.PreventTooManyIndexableTags(3, ignoreKinds, onlyKinds)
		ok, msg := fn(context.Background(), event)

		assert.Equal(t, false, ok)
		assert.Equal(t, "", msg)
	})

	t.Run("returns false when not too many tags", func(t *testing.T) {
		event := &nostr.Event{
			Kind: 0,
			Tags: nostr.Tags{
				nostr.Tag{"k", "v"},
			},
		}
		ignoreKinds := []int{}
		onlyKinds := []int{0}

		fn := policies.PreventTooManyIndexableTags(1, ignoreKinds, onlyKinds)
		ok, msg := fn(context.Background(), event)

		assert.Equal(t, false, ok)
		assert.Equal(t, "", msg)
	})

	t.Run("returns true when too many tags", func(t *testing.T) {
		event := &nostr.Event{
			Kind: 0,
			Tags: nostr.Tags{
				nostr.Tag{"k", "v"},
				nostr.Tag{"l", "m"},
				nostr.Tag{"m", "n"},
			},
		}
		ignoreKinds := []int{}
		onlyKinds := []int{0}

		fn := policies.PreventTooManyIndexableTags(1, ignoreKinds, onlyKinds)
		ok, msg := fn(context.Background(), event)

		assert.Equal(t, true, ok)
		assert.Equal(t, "too many indexable tags", msg)
	})
}
