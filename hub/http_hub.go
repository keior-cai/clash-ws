package hub

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type WsServer struct {
	port int16

	redis *redis.Client
}

func NewWsService(c *redis.Client, port int16) WsServer {
	return WsServer{
		redis: c,
		port:  port,
	}
}

func (s WsServer) Start() error {
	r := chi.NewRouter()
	r.Get("/", NewService(NewCheck(s.redis)))
	r.Get("/subject", NewSubject())
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), r)
}
