package hub

import (
	"github.com/sirupsen/logrus"
	"net"
	"time"
	"ws-server/outbound"
	"ws-server/service"
)

type CheckUserHandle struct {
	outbound.AdaptorClientHandle

	userService service.UserService
}

func (c CheckUserHandle) CallbackCreate(_, _, password, _ string, _ *outbound.WsClient, _ net.Conn) {
	userInfo := c.userService.GetByToken(password)
	if userInfo == nil || userInfo.Expire.Before(time.Now()) {
		logrus.Errorf("账号已经过期 %s", userInfo.Expire.Format(time.DateOnly))
		panic("账号过期失效")
	}
}

func NewCheck(s service.UserService) CheckUserHandle {
	return CheckUserHandle{
		userService: s,
	}
}
