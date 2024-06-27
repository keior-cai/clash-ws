package hub

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"ws-server/outbound"
)

type WsServer struct {
	Port int
	hs   []outbound.ClientHandle
	r    *chi.Mux
}

func NewWsService(port int) *WsServer {
	r := chi.NewRouter()
	return &WsServer{
		Port: port,
		r:    r,
	}
}

func (s WsServer) Start() error {
	s.r.Get("/", NewService(s.hs))
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.r)
}

func (s *WsServer) Add(h ...outbound.ClientHandle) {
	s.hs = append(s.hs, h...)
}

func (s *WsServer) AddRoute(f func(route *chi.Mux)) {
	f(s.r)
}
