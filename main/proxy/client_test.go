package main

import (
	"context"
	"net"
	"net/http"
	"testing"
	"ws-server/service"
)

func TestClient(t *testing.T) {

	h := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dial, _ := net.Dial("tcp", ":8082")
				client := service.Client{
					Conn: dial,
				}
				err := client.WriteAddr(addr)
				if err != nil {
					return nil, err
				}
				return client, nil
			},
		},
	}
	resp, err := h.Get("https://www.baidu.com")
	if err != nil {
		t.Log(err)
	}
	t.Log(resp.Status)
}
