package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"ws-server/main/gateway/config"
	"ws-server/main/gateway/service"
)

var (
	c string
)

func init() {
	flag.StringVar(&c, "conf", "", "gateway conf")
	flag.Parse()
	if c == "" {
		panic("can't find conf")
	}
}

func main() {
	conf, err := config.NewConf(c)
	if err != nil {
		return
	}
	proxy := service.NewProxy(conf)
	go func() {
		proxy.Start()
	}()
	termSign := make(chan os.Signal, 1)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-termSign:
			return
		}
	}

}
