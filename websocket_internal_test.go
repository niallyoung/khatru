package khatru

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fasthttp/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocket_WriteJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	header := http.Header{}

	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, _ := u.Upgrade(w, req, header)

	ws := WebSocket{
		conn:            conn,
		Request:         nil,
		Challenge:       "",
		AuthedPublicKey: "",
		Authed:          nil,
	}

	err := ws.WriteJSON(`{"foo": "bah"}`)
	assert.Error(t, err, "nil *Conn")

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Bad Request\n", string(data))
}

func TestWebSocket_WriteMessage(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	header := http.Header{}

	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, _ := u.Upgrade(w, req, header)

	ws := WebSocket{
		conn:            conn,
		Request:         nil,
		Challenge:       "",
		AuthedPublicKey: "",
		Authed:          nil,
	}

	err := ws.WriteMessage(1, []byte("foo"))
	assert.Error(t, err, "nil *Conn")

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Bad Request\n", string(data))
}
