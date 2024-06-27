package hub

import (
	"github.com/sirupsen/logrus"
	"net"
	"ws-server/outbound"
	"ws-server/service"
)

type CheckUserHandle struct {
	outbound.AdaptorClientHandle

	userService service.UserService
}

func (c CheckUserHandle) CallbackCreate(_, _, password, _ string, _ *outbound.WsClient, _ net.Conn) {
	if c.userService.Expire(password) {
		logrus.Errorf("账号已经过期")
		panic("账号过期失效")
	}
}

func NewCheck(s service.UserService) CheckUserHandle {
	return CheckUserHandle{
		userService: s,
	}
}
