package policies_test

import (
	"context"
	"github.com/fiatjaf/khatru/policies"
	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoComplexFilters(t *testing.T) {
	t.Run("returns false when <= 2 tags", func(t *testing.T) {
		tags := make(map[string][]string)
		tags["foo"] = []string{"k", "v"}
		tags["bah"] = []string{"m", "n"}
		filter := nostr.Filter{
			Kinds: []int{1, 2, 3},
			Tags:  tags,
		}
		ok, msg := policies.NoComplexFilters(context.Background(), filter)
		assert.Equal(t, false, ok)
		assert.Equal(t, "", msg)
	})

	t.Run("returns true when > 2 tags", func(t *testing.T) {
		tags := make(map[string][]string)
		tags["foo"] = []string{"k", "v"}
		tags["bah"] = []string{"m", "n"}
		tags["yes"] = []string{"n", "o"}
		tags["no"] = []string{"o", "p"}
		filter := nostr.Filter{
			Kinds: []int{1, 2, 3},
			Tags:  tags,
		}

		ok, msg := policies.NoComplexFilters(context.Background(), filter)
		assert.Equal(t, true, ok)
		assert.Equal(t, "too many things to filter for", msg)
	})
}

func TestNoEmptyFilters(t *testing.T) {
	t.Run("return true when filters empty", func(t *testing.T) {
		filter := nostr.Filter{}

		ok, msg := policies.NoEmptyFilters(context.Background(), filter)
		assert.Equal(t, true, ok)
		assert.Equal(t, "can't handle empty filters", msg)
	})

	t.Run("return false when filters empty", func(t *testing.T) {
		tags := make(map[string][]string)
		tags["foo"] = []string{"k", "v"}
		filter := nostr.Filter{
			Kinds: []int{1, 2, 3},
			Tags:  tags,
		}

		ok, msg := policies.NoEmptyFilters(context.Background(), filter)
		assert.Equal(t, false, ok)
		assert.Equal(t, "", msg)
	})
}

func TestAntiSyncBots(t *testing.T) {
	t.Run("return false when Kinds and Author is empty", func(t *testing.T) {
		filter := nostr.Filter{}

		reject, msg := policies.AntiSyncBots(context.Background(), filter)
		assert.Equal(t, true, reject)
		assert.Equal(t, "an author must be specified to get their kind:1 notes", msg)
	})
}

func TestNoSearchQueries(t *testing.T) {
	t.Run("returns false when filter.Search is empty", func(t *testing.T) {
		filter := nostr.Filter{}

		reject, msg := policies.NoSearchQueries(context.Background(), filter)
		assert.Equal(t, false, reject)
		assert.Equal(t, "", msg)
	})

	t.Run("returns true when filter.Search is not empty", func(t *testing.T) {
		filter := nostr.Filter{Search: "non-empty"}

		reject, msg := policies.NoSearchQueries(context.Background(), filter)
		assert.Equal(t, true, reject)
		assert.Equal(t, "search is not supported", msg)
	})
}

func TestRemoveSearchQueries(t *testing.T) {
	t.Run("empties Search and sets Limit to -1 when Search not empty", func(t *testing.T) {
		filter := &nostr.Filter{Search: "non-empty"}

		policies.RemoveSearchQueries(context.Background(), filter)
		assert.Equal(t, "", filter.Search)
		assert.Equal(t, -1, filter.Limit)
	})

	t.Run("does nothing when Search is empty", func(t *testing.T) {
		filter := &nostr.Filter{Search: "", Limit: 3}

		policies.RemoveSearchQueries(context.Background(), filter)
		assert.Equal(t, "", filter.Search)
		assert.Equal(t, 3, filter.Limit)
	})
}

func TestRemoveAllButKinds(t *testing.T) {
	t.Run("should remove all kinds", func(t *testing.T) {
		filter := &nostr.Filter{Kinds: []int{1, 2, 3}}

		fn := policies.RemoveAllButKinds(1, 3)
		fn(context.Background(), filter)
		assert.Equal(t, []int{1, 3}, filter.Kinds)
	})

	t.Run("should set Limit to -1 when all kinds are removed", func(t *testing.T) {
		filter := &nostr.Filter{Kinds: []int{2}, Limit: 2}

		fn := policies.RemoveAllButKinds(1, 3)
		fn(context.Background(), filter)
		assert.Equal(t, -1, filter.Limit)
	})
}

func TestRemoveAllButTags(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tags := make(map[string][]string)
		tags["foo"] = []string{"k", "v"}
		filter := &nostr.Filter{Tags: tags}

		fn := policies.RemoveAllButTags("bah")
		fn(context.Background(), filter)
		assert.Equal(t, 0, len(filter.Tags))
		assert.Equal(t, -1, filter.Limit)
	})
}
