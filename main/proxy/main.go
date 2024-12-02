package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"ws-server/service"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "port", 8082, "service port")
	flag.Parse()
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	for {
		accept, err := listen.Accept()
		if err != nil {
			continue
		}
		go func() {
			s := service.Service{
				Conn: accept,
			}
			defer func() {
				_ = s.Close()
			}()
			connect := s.Connect()
			ch := make(chan error)
			go func() {
				_, err2 := io.Copy(s, connect)
				ch <- err2
			}()
			_, _ = io.Copy(connect, s)
			<-ch
		}()
	}
}
