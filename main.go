package main

import (
	"encoding/json"
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"ws-server/config"
	"ws-server/hub"
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
	service := hub.NewWsService(ac.Port)
	service.Add(hub.NewCheck(userService))
	service.AddRoute(func(r *chi.Mux) {
		r.Get("/subject", hub.NewSubject(userService, ac))
		r.Route("/user", func(route chi.Router) {
			route.Use(func(handler http.Handler) http.Handler {
				return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					handler.ServeHTTP(writer, request)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")
				})
			}, func(handler http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer func() {
						if er := recover(); er != nil {
							w.Header().Set("Content-Type", "text/html; charset=utf-8")
							s, _ := er.(string)
							w.WriteHeader(http.StatusBadRequest)
							_, _ = w.Write([]byte(s))
						}
					}()
					handler.ServeHTTP(w, r)
				})
			})
			route.Put("/{name}/{mouth}", func(w http.ResponseWriter, rt *http.Request) {
				name := chi.URLParam(rt, "name")
				month := chi.URLParam(rt, "mouth")
				u := userService.GetByName(name)
				if u != nil {
					panic("用户已存在")
				}
				atoi, _ := strconv.Atoi(month)
				user := userService.AddUser(name, atoi)
				_ = json.NewEncoder(w).Encode(user)
			})

			route.Get("/{name}", func(w http.ResponseWriter, rt *http.Request) {
				name := chi.URLParam(rt, "name")
				userInfo := userService.GetByName(name)
				if userInfo == nil {
					panic("用户不存在")
				}
				_ = json.NewEncoder(w).Encode(userInfo)
			})
		})
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
			return
		}
	}
}
