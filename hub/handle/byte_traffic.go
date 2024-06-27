package handle

import (
	"io"
	"sync"
	"sync/atomic"
	"time"
	"ws-server/outbound"
	"ws-server/service"
)

type TrafficHandle struct {
	io.Closer
	outbound.AdaptorClientHandle

	s           service.UserService
	t           *time.Ticker
	uploadCache *sync.Map
	download    *sync.Map
}

func NewTrafficHandle(s service.UserService) TrafficHandle {
	handle := TrafficHandle{
		s:           s,
		t:           time.NewTicker(time.Minute),
		uploadCache: &sync.Map{},
		download:    &sync.Map{},
	}
	go func() {
		for range handle.t.C {
			handle.uploadCache.Range(func(key, value any) bool {
				if key == "" {
					return true
				}
				handle.uploadCache.Delete(key)
				a := value.(*atomic.Int64)
				s.UploadSize(a.Load(), key.(string))
				return true
			})

			handle.download.Range(func(key, value any) bool {
				if key == "" {
					return true
				}
				handle.download.Delete(key)
				a := value.(*atomic.Int64)
				s.Download(a.Load(), key.(string))
				return true
			})
		}
	}()
	return handle
}

func (h TrafficHandle) Read(token string, r []byte) {
	if value, ok := h.download.Load(token); ok {
		v := value.(*atomic.Int64)
		v.Add(int64(len(r)))
		return
	}
	a := new(atomic.Int64)
	a.Add(int64(len(r)))
	h.download.Store(token, a)
}

func (h TrafficHandle) Write(token string, r []byte) {
	if value, ok := h.uploadCache.Load(token); ok {
		v := value.(*atomic.Int64)
		v.Add(int64(len(r)))
		return
	}
	a := new(atomic.Int64)
	a.Add(int64(len(r)))
	h.uploadCache.Store(token, a)
}

func (h TrafficHandle) Close() error {
	h.t.Stop()
	return nil
}
