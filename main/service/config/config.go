package config

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type ClashWsConfig struct {
	Proxies []map[string]any
	Rules   []string
	Redis   RedisConfig
	ticker  *time.Ticker
	Http    HttpServer
}

func (c ClashWsConfig) Close() error {
	c.ticker.Stop()
	return nil
}

type HttpServer struct {
	Port   int    `yaml:"port,omitempty"`
	Secret string `yaml:"secret,omitempty"`
}

type ApplicationConfig struct {
	// 服务启动端口号
	Server HttpServer  `yaml:"http"`
	Redis  RedisConfig `yaml:"redis"`
	Path   string      `yaml:"path,omitempty"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Db       int    `yaml:"db,omitempty"`
}

type RowConfig struct {
	Proxies []map[string]any
	Rules   []string
}

func parseConfig(p string) (*RowConfig, error) {
	conf, err := parseProxies(p + "row.yaml")
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func parseProxies(path string) (*RowConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(file))
	config := &RowConfig{}
	err = decoder.Decode(config)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return config, nil
}

func NewConfig(file string) (*ClashWsConfig, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(f))
	config := &ApplicationConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	if config.Path == "" {
		config.Path = "./"
	}
	setDefaultHttp(&config.Server)
	wsConfig := &ClashWsConfig{
		Http:   config.Server,
		Redis:  config.Redis,
		ticker: time.NewTicker(3 * time.Second),
	}
	go func(path string, wsConfig *ClashWsConfig) {
		Start(path, wsConfig)
	}(config.Path, wsConfig)
	return wsConfig, nil

}

func setDefaultHttp(server *HttpServer) {
	if server.Port == 0 {
		server.Port = 8081
	}
}

func Start(path string, config *ClashWsConfig) {
	for range config.ticker.C {
		if path == "" {
			continue
		}
		c, err := parseConfig(path)
		if err != nil {
			logrus.Errorf("parse config error %s", err)
		}
		config.Rules = c.Rules
		config.Proxies = c.Proxies
	}
}
