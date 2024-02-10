package policies_test

import (
	"context"
	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/fiatjaf/khatru/policies"
)

func TestPreventTimestampsInTheFuture(t *testing.T) {
	t.Run("stub", func(t *testing.T) {
		fn := policies.PreventTimestampsInTheFuture(nostr.Now())

		ok, result := fn(context.Background(), &nostr.Event{})
		assert.Equal(t, false, ok)
		assert.Equal(t, "", result)
	})
}
