package hub

import (
	"context"
	"github.com/go-redis/redis/v8"
	"net"
	"time"
	"ws-server/outbound"
)

var redisUserZsetKey = "clash:ws:user"

type CheckUserHandle struct {
	outbound.AdaptorClientHandle

	client *redis.Client
}

func (c CheckUserHandle) CallbackCreate(_, username, _ string, _ *outbound.WsClient, _ net.Conn) {
	ctx := context.TODO()
	score := c.client.ZScore(ctx, redisUserZsetKey, username)
	second := time.Now().Second()
	if score.Val() < float64(second) {
		panic("账号过期失效")
	}
}

func NewCheck(client *redis.Client) CheckUserHandle {
	return CheckUserHandle{
		client: client,
	}
}
