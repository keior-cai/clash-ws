package outbound

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
	"ws-server/cipher"
	"ws-server/config"
	"ws-server/utils"
)

const (
	ACK_USER_MESSAGE = 6
	SEND_PUBLIC      = 5
	CONNECT_SUCCESS  = 4
)

type ClientHandle interface {
	Read(id string, r []byte)

	Write(id string, w []byte)

	CallbackClose(id string)

	CallbackCreate(id, userName, host string, c *WsClient, d net.Conn)

	CallbackRegister(id string, c *WsClient)
}

type AdaptorClientHandle struct {
}

func (a AdaptorClientHandle) Read(_ string, _ []byte) {

}

func (a AdaptorClientHandle) Write(_ string, _ []byte) {

}

func (a AdaptorClientHandle) CallbackClose(_ string) {

}

func (a AdaptorClientHandle) CallbackCreate(_, _, _ string, _ *WsClient, _ net.Conn) {

}

func (a AdaptorClientHandle) CallbackRegister(_ string, _ *WsClient) {

}

type WsClient struct {
	Id       string
	conn     *websocket.Conn
	cipher   cipher.Cipher
	buff     []byte
	UserName string
	Password string
	handles  []ClientHandle
	logger   logrus.Logger
}

// 这里是请求信息
type SendConnectInfo struct {
	Command  int    `json:"command"`
	Method   string `json:"method"`
	Username string `json:"username"`
	Password string `json:"password"`
	Random   string `json:"random"`
}

type ConnectInfo struct {
	Command int    `json:"command"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Network string `json:"network"`
}

type command struct {
	Command int `json:"command"`
}

type publicKey struct {
	Command   int    `json:"command"`
	PublicKey string `json:"publicKey"`
}

func (s *WsClient) Read(p []byte) (int, error) {
	if len(s.buff) > 0 {
		i := copy(p, s.buff)
		s.buff = s.buff[i:]
		for _, d := range s.handles {
			d.Read(s.Id, p)
		}
		return i, nil
	}
	messageType, b, err := s.conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	if messageType == websocket.BinaryMessage {
		decrypt, _ := s.cipher.Decrypt(b)
		i := copy(p, decrypt)
		if len(decrypt) > i {
			s.buff = decrypt[i:]
		}
		for _, d := range s.handles {
			d.Read(s.Id, p)
		}
		return i, err
	} else if messageType == websocket.CloseMessage {
		return 0, io.EOF
	}
	return 0, io.EOF
}

func (s *WsClient) Write(p []byte) (int, error) {
	encrypt, err := s.cipher.Encrypt(p)
	if err != nil {
		return 0, err
	}
	err = s.conn.WriteMessage(websocket.BinaryMessage, encrypt)
	if err != nil {
		return 0, err
	}
	for _, d := range s.handles {
		d.Write(s.Id, encrypt)
	}
	return len(p), err
}

func (s *WsClient) Close() error {
	for _, d := range s.handles {
		d.CallbackClose(s.Id)
	}
	return s.conn.Close()
}

func (s *WsClient) Create() net.Conn {
	s.writeMessage(&publicKey{
		Command:   SEND_PUBLIC,
		PublicKey: utils.RsaInstant.PublicKey,
	})
	_, p, err := s.conn.ReadMessage()
	ss := utils.AesInstant.DecryptStr(string(p))
	info := &SendConnectInfo{}
	_ = json.Unmarshal([]byte(ss), info)
	if config.NameMap[info.Username] == "" {
		// 鉴权失败
		return nil
	}
	s.UserName = info.Username
	decryptStr := utils.AesInstant.DecryptStr(info.Random)
	str, _ := utils.RsaInstant.RsaDecryptStr(decryptStr)
	s.cipher = cipher.NewCipher(info.Method, str)
	s.writeMessage(&command{Command: ACK_USER_MESSAGE})
	_, bytes, err := s.conn.ReadMessage()
	b := utils.AesInstant.DecryptStr(string(bytes))
	connectInfo := &ConnectInfo{}
	_ = json.Unmarshal([]byte(b), connectInfo)
	if connectInfo.Host == "" || connectInfo.Port <= 0 {
		return nil
	}
	if connectInfo.Network == "" {
		connectInfo.Network = "tcp"
	}
	host := fmt.Sprintf("%s:%d", connectInfo.Host, connectInfo.Port)
	conn, err := net.DialTimeout(connectInfo.Network, host, time.Second*3)
	if err != nil {
		return nil
	}
	logrus.Infof("%s HOST = %s", s.UserName, host)
	for _, h := range s.handles {
		h.CallbackCreate(s.Id, info.Username, host, s, conn)
	}
	s.writeMessage(&command{Command: CONNECT_SUCCESS})
	return conn
}

func (s *WsClient) writeMessage(any any) {
	data, _ := json.Marshal(any)
	_ = s.conn.WriteMessage(websocket.TextMessage, []byte(utils.AesInstant.Encrypt(string(data))))
}

func NewWs(conn *websocket.Conn, handles []ClientHandle) *WsClient {
	id := uuid.New().String()
	client := &WsClient{
		Id:      id,
		conn:    conn,
		handles: handles,
	}
	for _, h := range handles {
		h.CallbackRegister(id, client)
	}
	return client
}