package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

const (
	__name__ = "clash:redis:name:token"
	// 用户信息存储对象， 使用hash
	__info__ = "clash:redis:user:info"
)

type UserInfo struct {
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Token    string    `json:"token"`
	Expire   time.Time `json:"expire"`
	Upload   int64     `json:"upload"`
	Download int64     `json:"download"`
}

type UserService interface {
	GetByToken(token string) *UserInfo
	GetByName(name string) *UserInfo
	AddUser(name string, month int) UserInfo
	UploadSize(size int64, token string)
	Download(size int64, token string)
}

type RedisUser struct {
	r *redis.Client
	c context.Context
}

func NewRedisService(r *redis.Client) UserService {
	return &RedisUser{
		r: r,
		c: context.TODO(),
	}
}

func (r RedisUser) GetByToken(token string) *UserInfo {
	all := r.r.HGetAll(r.c, formatKey(__info__, token)).Val()
	return parseUserInfo(all)
}

func (r RedisUser) GetByName(name string) *UserInfo {
	token := r.r.HGet(r.c, __name__, name)
	if token.Err() != nil {
		return nil
	}
	all := r.r.HGetAll(r.c, formatKey(__info__, token.Val())).Val()
	if len(all) == 0 {
		return nil
	}
	return parseUserInfo(all)
}

func parseUserInfo(m map[string]string) *UserInfo {
	upload, _ := strconv.ParseInt(m["upload"], 10, 64)
	download, _ := strconv.ParseInt(m["download"], 10, 64)
	expireTime, _ := time.Parse(time.DateOnly, m["expire"])
	return &UserInfo{
		Upload:   upload,
		Download: download,
		Name:     m["name"],
		Password: m["password"],
		Token:    m["token"],
		Expire:   expireTime,
	}
}

func (r RedisUser) UploadSize(size int64, token string) {
	r.r.HIncrBy(r.c, formatKey(__info__, token), "upload", size)
}

func (r RedisUser) Download(size int64, token string) {
	r.r.HIncrBy(r.c, formatKey(__info__, token), "download", size)
}

func (r RedisUser) AddUser(name string, month int) UserInfo {
	u, _ := uuid.NewUUID()
	info := UserInfo{
		Name:     name,
		Password: u.String(),
		Token:    u.String(),
		Expire:   time.Now().AddDate(0, month, 0),
	}
	set := r.r.HSet(r.c, formatKey(__info__, u.String()),
		"name", info.Name,
		"password", info.Password,
		"token", info.Token,
		"expire", info.Expire.Format(time.DateOnly),
		"upload", "0",
		"download", "0",
	)
	if set.Err() != nil {
		panic(set.Err())
	}
	r.r.HSet(r.c, __name__, name, info.Token)
	return info
}

func formatKey(s ...string) string {
	return strings.Join(s, ":")
}
