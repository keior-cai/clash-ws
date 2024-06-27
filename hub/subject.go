package hub

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"strings"
	"ws-server/config"
	"ws-server/service"
	subject2 "ws-server/subject"
)

type ProxyGroup struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
}

func NewSubject(s service.UserService, c *config.ClashWsConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if s.Expire(token) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("账号已到期"))
			return
		}
		u := s.GetByToken(token)
		groupName := "新加坡🇸🇬"
		group := ProxyGroup{
			Name: groupName,
			Type: "select",
		}
		var list []map[string]any
		for _, p := range c.Proxies {
			name := p["name"].(string)
			p["username"] = u.Name
			p["password"] = u.Password
			list = append(list, p)
			group.Proxies = append(group.Proxies, name)
		}
		var listRules []string
		for _, s := range c.Rules {
			listRules = append(listRules, strings.Replace(s, "{name}", groupName, -1))
		}
		subject := &subject2.Subject{
			Proxies:     c.Proxies,
			ProxyGroups: []any{group},
			Rules:       listRules,
		}
		out, _ := yaml.Marshal(subject)
		w.Header().Set("Content-Disposition", "attachment; filename=subject.yaml")
		w.Header().Set("subscription-userinfo", fmt.Sprintf("upload=%d; download=%d; total=1024000; expire=%d", u.Upload, u.Download, u.Expire.Unix()))
		w.Header().Set("profile-web-page-url", "https://chuangjie.icu")
		_, _ = w.Write(out)
	}
}
