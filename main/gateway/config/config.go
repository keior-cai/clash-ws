package config

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"math"
	"os"
)

type GatewayConf struct {
	Port   int      `yaml:"port,omitempty"`
	Server []string `yaml:"server"`
}

func NewConf(file string) (*GatewayConf, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(f))
	conf := &GatewayConf{}
	err = decoder.Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (g GatewayConf) RandomServer() string {
	if len(g.Server) == 0 {
		panic("服务配置不存在")
	}
	round := math.Round(float64(len(g.Server)))
	return g.Server[int(round)]
}
