package statics

import (
	"fmt"
	"ws-server/outbound"
)

type TrafficUnit uint64

func (s TrafficUnit) String() string {
	if s <= Byte {
		return fmt.Sprintf("0B")
	} else if s < KB {
		return fmt.Sprintf("%dB", s)
	} else if s < MB {
		d := float64(s) / float64(KB)
		return fmt.Sprintf("%fKB", d)
	} else if s < GB {
		d := float64(s) / float64(MB)
		return fmt.Sprintf("%fMB", d)
	} else if s < TB {
		d := float64(s) / float64(GB)
		return fmt.Sprintf("%fGB", d)
	} else if s < PB {
		d := float64(s) / float64(TB)
		return fmt.Sprintf("%fGB", d)
	} else if s < EB {
		d := float64(s) / float64(PB)
		return fmt.Sprintf("%fTB", d)
	}
	return "0"
}

const (
	Byte TrafficUnit = 1.0
	KB               = 1024 * Byte
	MB               = 1024 * KB
	GB               = 1024 * MB
	TB               = 1024 * GB
	PB               = 1024 * TB
	EB               = 1024 * PB
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
