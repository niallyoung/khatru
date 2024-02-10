package khatru_test

import (
	"github.com/fiatjaf/khatru"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRelay(t *testing.T) {
	t.Run("constructor returns a Relay", func(t *testing.T) {
		relay := khatru.NewRelay()
		assert.NotNil(t, relay)
	})

	//t.Run("foo", func(t *testing.T) {
	//	relay := khatru.NewRelay()
	//	assert.Equal(t,
	//		&khatru.Relay{
	//			ServiceURL:"",
	//			RejectEvent: []func(context.Context, *nostr.Event) (bool, string)(nil),
	//			RejectFilter: []func(context.Context, nostr.Filter) (bool, string)(nil),
	//			RejectCountFilter: []func(context.Context, nostr.Filter) (bool, string)(nil),
	//			OverwriteDeletionOutcome: []func(context.Context, *nostr.Event, *nostr.Event) (bool, string)(nil),
	//			OverwriteResponseEvent: []func(context.Context, *nostr.Event)(nil),
	//			OverwriteFilter: []func(context.Context, *nostr.Filter)(nil),
	//			OverwriteCountFilter: []func(context.Context, *nostr.Filter)(nil),
	//			OverwriteRelayInformation: []func(context.Context, *http.Request, nip11.RelayInformationDocument) nip11.RelayInformationDocument(nil),
	//			StoreEvent: []func(context.Context, *nostr.Event) error(nil),
	//			DeleteEvent: []func(context.Context, *nostr.Event) error(nil),
	//			QueryEvents: []func(context.Context, nostr.Filter) (chan *nostr.Event, error)(nil),
	//			CountEvents: []func(context.Context, nostr.Filter) (int64, error)(nil),
	//			OnConnect: []func(context.Context)(nil),
	//			OnDisconnect: []func(context.Context)(nil),
	//			OnEventSaved: []func(context.Context, *nostr.Event)(nil),
	//			OnEphemeralEvent: []func(context.Context, *nostr.Event)(nil),
	//			Info: (*nip11.RelayInformationDocument)(0xc000132000),
	//			Log: (*log.Logger)(0xc000123500),
	//			upgrader: websocket.Upgrader{
	//				HandshakeTimeout:0,
	//				ReadBufferSize:1024,
	//				WriteBufferSize:1024,
	//				WriteBufferPool: websocket.BufferPool(nil),
	//				Subprotocols: []string(nil),
	//				Error: (func(http.ResponseWriter, *http.Request, int, error))(nil),
	//				CheckOrigin: (func(*http.Request) bool)(0x127bb20),
	//				EnableCompression: false,
	//			},
	//			clients: (*xsync.MapOf[*github.com/fasthttp/websocket.Conn,struct {}])(0xc00012e9a0),
	//			Addr:"",
	//			serveMux: (*http.ServeMux)(0xc000090e00),
	//			httpServer: (*http.Server)(nil),
	//			WriteWait:10000000000,
	//			PongWait:60000000000,
	//			PingPeriod:30000000000,
	//			MaxMessageSize:512000,
	//		}, relay)
	//})
}
