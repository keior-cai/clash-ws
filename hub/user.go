package hub

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"ws-server/config"
	"ws-server/service"
)

type UserHub struct {
	s service.UserService
	c config.HttpServer
}

func NewUserHub(c config.HttpServer, s service.UserService) UserHub {
	return UserHub{
		s: s,
		c: c,
	}
}

func (uh UserHub) Route(route chi.Router) {
	route.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			handler.ServeHTTP(writer, request)
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
	}, func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if uh.c.Secret == "" {
				handler.ServeHTTP(w, r)
				return
			}
			auth := r.Header.Get("Authorization")
			if auth != " "+uh.c.Secret {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			handler.ServeHTTP(w, r)
		})
	})
	route.Put("/{name}/{mouth}", func(w http.ResponseWriter, rt *http.Request) {
		name := chi.URLParam(rt, "name")
		month := chi.URLParam(rt, "mouth")
		u := uh.s.GetByName(name)
		if u != nil {
			panic("用户已存在")
		}
		atoi, _ := strconv.Atoi(month)
		user := uh.s.AddUser(name, atoi)
		_ = json.NewEncoder(w).Encode(user)
	})

	route.Get("/{name}", func(w http.ResponseWriter, rt *http.Request) {
		name := chi.URLParam(rt, "name")
		userInfo := uh.s.GetByName(name)
		if userInfo == nil {
			panic("用户不存在")
		}
		_ = json.NewEncoder(w).Encode(userInfo)
	})

	route.Get("/", func(w http.ResponseWriter, rt *http.Request) {
		_ = json.NewEncoder(w).Encode(uh.s.List())
	})

	route.Delete("/{name}", func(w http.ResponseWriter, rt *http.Request) {
		name := chi.URLParam(rt, "name")
		uh.s.Delete(name)
	})
}
