package main

import (
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"ws-server/config"
	"ws-server/hub"
	"ws-server/hub/handle"
	service2 "ws-server/service"
)

var (
	c string
)

func init() {
	flag.StringVar(&c, "c", "", "config file path")
	flag.Parse()
	if c == "" {
		panic("not find config path")
	}
}

func main() {

	ac, err := config.NewConfig(c)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		DisableQuote:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.DebugLevel)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     ac.Redis.Addr,
		DB:       ac.Redis.Db,
		Username: ac.Redis.Username,
		Password: ac.Redis.Password,
	})
	userService := service2.NewRedisService(redisClient)
	service := hub.NewWsService(ac.Http.Port)
	trafficHandle := handle.NewTrafficHandle(userService)
	service.Add(handle.NewCheck(userService), trafficHandle)
	service.AddRoute(func(r *chi.Mux) {
		r.Get("/subject", hub.NewSubject(userService, ac))
		r.Route("/user", hub.NewUserHub(ac.Http, userService).Route)
	})
	go func() {
		logrus.Infof("start server :%d", service.Port)
		_ = service.Start()
	}()
	termSign := make(chan os.Signal, 1)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-termSign:
			_ = ac.Close()
			_ = trafficHandle.Close()
			return
		}
	}
}
