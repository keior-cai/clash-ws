package main

import (
	"flag"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"syscall"
	"ws-server/config"
	"ws-server/hub"
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
	r := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/clash.log", // 日志文件路径
		MaxSize:    100,                // 每个日志文件最大为 10 MB
		MaxBackups: 30,                 // 保留最近的 3 个日志文件
		MaxAge:     28,                 // 保留最近 28 天的日志
		Compress:   true,               // 压缩旧的日志文件
		LocalTime:  true,
	})
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		DisableQuote:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
	err := config.ParseConfig(c)
	if err != nil {
		return
	}
	service := hub.NewWsService(r, int16(8081))
	go func() {
		_ = service.Start()
	}()
	termSign := make(chan os.Signal, 1)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-termSign:
			config.Stop()
			return
		}
	}
}
