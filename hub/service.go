package hub

import (
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"time"
	"ws-server/outbound"
	"ws-server/statics"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: time.Second * 10,
}

func NewService(h ...outbound.ClientHandle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrade, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		ws := outbound.NewWs(upgrade, h)
		defer func() {
			_ = ws.Close()
		}()
		create := ws.Create()
		if create == nil {
			return
		}
		traffic := statics.NewTraffic(ws)
		ch := make(chan error)
		go func() {
			_, err2 := io.Copy(traffic, create)
			ch <- err2
		}()
		_, _ = io.Copy(create, traffic)
		<-ch
	}
}
