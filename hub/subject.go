package hub

import (
	"gopkg.in/yaml.v3"
	"net/http"
	"strings"
	"ws-server/config"
	subject2 "ws-server/subject"
)

type ProxyGroup struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
}

func NewSubject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		userName := config.TokenMap[token]
		if token == "" || userName == "" {
			w.WriteHeader(400)
			return
		}
		groupName := "新加坡🇸🇬"
		group := ProxyGroup{
			Name: groupName,
			Type: "select",
		}
		var list []map[string]any
		for _, p := range config.Proxies {
			name := p["name"].(string)
			p["username"] = userName
			list = append(list, p)
			group.Proxies = append(group.Proxies, name)
		}
		var listRules []string
		for _, s := range config.Rules {
			listRules = append(listRules, strings.Replace(s, "{name}", groupName, -1))
		}
		subject := &subject2.Subject{
			Proxies:     config.Proxies,
			ProxyGroups: []any{group},
			Rules:       listRules,
		}
		out, _ := yaml.Marshal(subject)
		w.Header().Set("Content-Disposition", "attachment; filename=subject.yaml")
		w.Header().Set("subscription-userinfo", "upload=1234; download=2234; total=1024000; expire=2218532293")
		w.Header().Set("profile-web-page-url", "https://chuangjie.icu")
		_, _ = w.Write(out)
	}
}
