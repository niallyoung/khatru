package khatru

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nbd-wtf/go-nostr"
)

func TestGetListeningFilters(t *testing.T) {
	t.Run("TBC", func(t *testing.T) {

		filters := GetListeningFilters()
		emptyFilters := make(nostr.Filters, 0, listeners.Size()*2)
		assert.Equal(t, emptyFilters, filters)
		//assert.IsType(t, []nostr.Filter, filters)
	})
}
