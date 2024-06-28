package service

import (
	"context"
	"errors"
	"github.com/bluele/gcache"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	__name__ = "clash:redis:name:token"
	// 用户信息存储对象， 使用hash
	__info__ = "clash:redis:user:info"

	__traffic__ = "clash:redis:user:traffic"
)

const (
	Byte uint64 = 1.0
	KB          = 1024 * Byte
	MB          = 1024 * KB
	GB          = 1024 * MB
	TB          = 1024 * GB
	PB          = 1024 * TB
	EB          = 1024 * PB
)

type UserInfo struct {
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Token    string    `json:"token"`
	Expire   time.Time `json:"expire"`
	Total    uint64    `json:"total"`
}

type Traffic struct {
	Download int64 `json:"download"`
	Upload   int64 `json:"upload"`
}

type UserService interface {
	GetByToken(token string) *UserInfo
	GetByName(name string) *UserInfo
	AddTotalTraffic(name string, size int)
	AddExpireTime(name string, day int)
	List() []string
	Delete(name string)
	AddUser(name string, day int) UserInfo
	UploadSize(size int64, token string)
	Download(size int64, token string)
	Traffic(token string) Traffic
	Expire(token string) bool
}

type RedisUser struct {
	r     *redis.Client
	c     context.Context
	cache gcache.Cache
}

func NewRedisService(r *redis.Client) UserService {
	c := context.TODO()
	return &RedisUser{
		r: r,
		c: c,
		cache: gcache.New(500).LFU().LoaderFunc(func(i interface{}) (interface{}, error) {
			logrus.Debugf("user cache %s", i)
			all := r.HGetAll(c, formatKey(__info__, i.(string))).Val()
			if len(all) == 0 {
				return nil, errors.New("加载数据失败")
			}
			return parseUserInfo(all), nil
		}).Expiration(time.Hour).Build(),
	}
}

func (r RedisUser) GetByToken(token string) *UserInfo {
	userInfo, err := r.cache.Get(token)
	if err != nil {
		return nil
	}
	return userInfo.(*UserInfo)
}

func (r RedisUser) GetByName(name string) *UserInfo {
	token := r.r.HGet(r.c, __name__, name)
	if token.Err() != nil {
		return nil
	}
	return r.GetByToken(token.Val())
}

func parseUserInfo(m map[string]string) *UserInfo {
	total, _ := strconv.ParseUint(m["total"], 10, 64)
	expireTime, _ := time.Parse(time.DateOnly, m["expire"])
	return &UserInfo{
		Name:     m["name"],
		Password: m["password"],
		Token:    m["token"],
		Expire:   expireTime,
		Total:    total,
	}
}

func (r RedisUser) UploadSize(size int64, token string) {
	month := time.Now().Month()
	key := formatKey(__traffic__, token, strconv.FormatInt(int64(month), 10))
	in := r.r.HIncrBy(r.c, key, "upload", size)
	if in.Err() != nil && in.Val() == size {
		r.r.Expire(r.c, key, time.Hour*24*32)
	}
}

func (r RedisUser) Download(size int64, token string) {
	month := time.Now().Month()
	key := formatKey(__traffic__, token, strconv.FormatInt(int64(month), 10))
	in := r.r.HIncrBy(r.c, key, "download", size)
	if in.Err() != nil && in.Val() == size {
		r.r.Expire(r.c, key, time.Hour*24*32)
	}
}

func (r RedisUser) AddUser(name string, day int) UserInfo {
	u, _ := uuid.NewUUID()
	info := UserInfo{
		Name:     name,
		Password: u.String(),
		Token:    u.String(),
		Expire:   time.Now().AddDate(0, 0, day),
		Total:    50 * GB,
	}

	if day < 0 {
		info.Expire = time.Unix(0, 0)
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

func (r RedisUser) Expire(token string) bool {
	value := r.r.HGet(r.c, formatKey(__info__, token), "expire")
	if value.Err() != nil || value.Val() == "" {
		return true
	}
	t, err := time.Parse(time.DateOnly, value.Val())
	if err != nil {
		return true
	}
	if t.Unix() != 0 && t.Before(time.Now()) {
		return true
	}
	return false
}

func (r RedisUser) Delete(name string) {
	userInfo := r.GetByName(name)
	if userInfo == nil {
		panic("用户不存在")
	}
	del := r.r.Del(r.c, formatKey(__info__, userInfo.Token))
	if del.Err() != nil || del.Val() <= 0 {
		panic("删除失败")
	}
	r.r.HDel(r.c, __name__, userInfo.Name)
}

func (r RedisUser) List() []string {
	all := r.r.HGetAll(r.c, __name__)
	if all.Err() != nil {
		return nil
	}
	var list []string
	for k := range all.Val() {
		list = append(list, k)
	}
	return list
}

func (r RedisUser) AddTotalTraffic(name string, size int) {
	user := r.GetByName(name)
	if user == nil {
		panic("用户不存在")
	}
	r.r.HIncrBy(r.c, formatKey(__info__, user.Token), "total", int64(size)*int64(GB))
}

func (r RedisUser) AddExpireTime(name string, day int) {
	user := r.GetByName(name)
	if user == nil {
		panic("用户不存在")
	}
	user.Expire.Add(time.Hour * 24 * time.Duration(day))
	r.r.HSet(r.c, formatKey(__info__, user.Token), user.Expire.Format(time.DateOnly))
}

func (r RedisUser) Traffic(token string) Traffic {
	month := time.Now().Month()
	key := formatKey(__traffic__, token, strconv.FormatInt(int64(month), 10))
	all := r.r.HGetAll(r.c, key)
	if all.Err() != nil {
		panic(all.Err())
	}
	val := all.Val()
	upload, _ := strconv.ParseInt(val["upload"], 10, 64)
	download, _ := strconv.ParseInt(val["download"], 10, 64)
	return Traffic{
		Upload:   upload,
		Download: download,
	}
}

func formatKey(s ...string) string {
	return strings.Join(s, ":")
}
