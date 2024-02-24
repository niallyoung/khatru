package khatru

import (
	"net/http"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func Test_isOlder(t *testing.T) {
	scenarios := []struct {
		Name           string
		previousEvent  *nostr.Event
		nextEvent      *nostr.Event
		expectedResult bool
	}{
		{
			"next is older",
			&nostr.Event{CreatedAt: 0},
			&nostr.Event{CreatedAt: 1},
			true,
		},
		{
			"previous is older",
			&nostr.Event{CreatedAt: 1},
			&nostr.Event{CreatedAt: 0},
			false,
		},
		{
			"same age",
			&nostr.Event{CreatedAt: 1},
			&nostr.Event{CreatedAt: 1},
			false,
		},
		{
			"same age, but next has larger ID",
			&nostr.Event{CreatedAt: 1, ID: "a"},
			&nostr.Event{CreatedAt: 1, ID: "b"},
			false,
		},
		{
			"same age, but previous has larger ID",
			&nostr.Event{CreatedAt: 1, ID: "b"},
			&nostr.Event{CreatedAt: 1, ID: "a"},
			true,
		},
		{
			"same age, but previous has same ID",
			&nostr.Event{CreatedAt: 1, ID: "a"},
			&nostr.Event{CreatedAt: 1, ID: "a"},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			result := isOlder(s.previousEvent, s.nextEvent)
			assert.Equal(t, s.expectedResult, result)
		})
	}
}

func Test_getServiceBaseURL(t *testing.T) {
	scenarios := []struct {
		Name           string
		req            http.Request
		expectedString string
	}{
		{
			"host via request",
			http.Request{Host: "request-host"},
			"https://request-host",
		},
		{
			"host via header",
			http.Request{Host: "request-host", Header: http.Header{"X-Forwarded-Host": {"header-host"}}},
			"https://header-host",
		},
		{
			"localhost via request gives http",
			http.Request{Host: "localhost"},
			"http://localhost",
		},
		{
			"localhost via header gives http",
			http.Request{Host: "request-host", Header: http.Header{"X-Forwarded-Host": {"localhost"}}},
			"http://localhost",
		},
		{
			"host via request has a port number",
			http.Request{Host: "request-host:1234"},
			"http://request-host:1234",
		},
		{
			"host via header has a port number",
			http.Request{Host: "request-host", Header: http.Header{"X-Forwarded-Host": {"header-host:1234"}}},
			"http://header-host:1234",
		},
		{
			"host via request is naked IP",
			http.Request{Host: "192.168.63.1"},
			"http://192.168.63.1",
		},
		{
			"host via header is naked IP",
			http.Request{Host: "request-host", Header: http.Header{"X-Forwarded-Host": {"192.168.128.7"}}},
			"http://192.168.128.7",
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			assert.Equal(t, s.expectedString, getServiceBaseURL(&s.req))
		})
	}
}
