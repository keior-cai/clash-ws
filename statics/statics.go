package statics

import (
	"ws-server/outbound"
)

var (
	userTraffic map[string]*StaticTraffic = make(map[string]*StaticTraffic)
)

type StaticTraffic struct {
	Download int64 `json:"download"`
	Upload   int64 `json:"upload"`
	Total    int64 `json:"total"`
}

type TrafficConn struct {
	conn *outbound.WsClient
}

func (t TrafficConn) Read(p []byte) (n int, err error) {
	return t.conn.Read(p)
}

func (t TrafficConn) Write(p []byte) (n int, err error) {
	return t.conn.Write(p)
}

func NewTraffic(conn *outbound.WsClient) TrafficConn {
	return TrafficConn{
		conn: conn,
	}
}
