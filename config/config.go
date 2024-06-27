package config

import (
	"bufio"
	"bytes"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"
)

var (
	Proxies  []map[string]any
	Rules    []string
	TokenMap = make(map[string]string)
	NameMap  = make(map[string]string)
	ticker   = time.NewTicker(3 * time.Second)
	path     string
)

type RowConfig struct {
	Proxies []map[string]any
	Rules   []string
}

func init() {
	go func() {
		for range ticker.C {
			if path == "" {
				continue
			}
			err := parseConfig(path)
			if err != nil {
				logrus.Errorf("parse config error %s", err)
			}
		}
	}()
}

func ParseConfig(p string) error {
	path = p
	return parseConfig(p)
}

func parseConfig(p string) error {
	err := parseProxies(p + "/application.yaml")
	if err != nil {
		return err
	}
	err = parseTokenMap(p + "/token")
	if err != nil {
		return err
	}
	return nil
}

func parseProxies(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(file))
	config := &RowConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	Proxies = config.Proxies
	Rules = config.Rules
	return nil
}

func parseTokenMap(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(file)

	newReader := bufio.NewReader(reader)
	m1 := make(map[string]string)
	m2 := make(map[string]string)
	for {
		line, _, err := newReader.ReadLine()
		if err != nil {
			break
		}
		before, after, found := strings.Cut(string(line), "=")
		if !found {
			continue
		}
		before = strings.Trim(before, " ")
		after = strings.Trim(after, " ")
		m1[before] = after
		m2[after] = before
	}
	TokenMap = m1
	NameMap = m2
	return nil
}

func Stop() {
	ticker.Stop()
}
