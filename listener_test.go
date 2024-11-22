package khatru

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nbd-wtf/go-nostr"
)

func TestGetListeningFilters(t *testing.T) {
	t.Run("get listening filters with no listeners returns empty", func(t *testing.T) {

		filters := GetListeningFilters()
		emptyFilters := make(nostr.Filters, 0, listeners.Size()*2)
		assert.Equal(t, emptyFilters, filters)
		assert.Len(t, filters, 0)
	})
}
