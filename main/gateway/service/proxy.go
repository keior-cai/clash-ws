package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"net/http"
	"time"
	"ws-server/cipher"
	"ws-server/main/gateway/config"
	"ws-server/service"
	"ws-server/utils"
)

const (
	// ConnectSuccess 发送连接成功命令
	ConnectSuccess = 4
	// PublicKeyCommand 发送public key 消息命令
	PublicKeyCommand int = 5
	// AckUserMessage 发送接到用户信息命令
	AckUserMessage = 6
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: time.Second * 10,
}

type Proxy struct {
	c         *config.GatewayConf
	publicKey string
}

func NewProxy(c *config.GatewayConf) Proxy {
	return Proxy{
		c:         c,
		publicKey: utils.RsaInstant.PublicKey,
	}
}

func (p Proxy) Start() {
	mux := chi.NewMux()
	mux.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrade, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				handler.ServeHTTP(w, r)
				return
			}
			client := ProxyClient{
				c:    upgrade,
				meta: &ConnectMeta{},
			}
			defer func() {
				_ = client.Close()
			}()
			client.sendAnyObject(&publicKey{
				Command:   PublicKeyCommand,
				PublicKey: p.publicKey,
			})
			if err != nil {
				return
			}
			info := &sendConnectInfo{}
			client.readTextMsg(info)
			// 鉴权, 限流
			client.meta.Username = info.Username
			client.meta.Password = info.Password
			client.meta.Method = info.Method
			client.meta.Random = info.Random
			client.cipher = cipher.NewCipher(info.Method, info.GetDecodeRandom())
			client.sendAnyObject(&command{
				Command: AckUserMessage,
			})
			cinfo := &connectInfo{}
			client.readTextMsg(cinfo)
			// 服务客户端
			dial, err := net.Dial("tcp", p.c.RandomServer())
			if err != nil {
				panic(err)
			}
			s := service.Client{
				Conn: dial,
			}
			err = s.WriteAddr(cinfo.Addr())
			if err != nil {
				panic(err)
			}
			client.sendAnyObject(&command{
				Command: ConnectSuccess,
			})
			ch := make(chan error)
			go func() {
				_, err2 := io.Copy(s, client)
				ch <- err2
			}()
			_, _ = io.Copy(client, s)
			<-ch
		})
	})
	_ = http.ListenAndServe(fmt.Sprintf(":%d", p.c.Port), mux)
}

type ProxyClient struct {
	io.Closer
	io.Reader
	io.Writer
	c      *websocket.Conn
	meta   *ConnectMeta
	cipher cipher.Cipher
	buff   []byte
}

type ConnectMeta struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Method   string `json:"method"`
	Random   string `json:"random"`
}

func (p ProxyClient) Close() error {
	return p.c.Close()
}

func (p ProxyClient) Read(b []byte) (int, error) {
	if len(p.buff) > 0 {
		i := copy(b, p.buff)
		p.buff = p.buff[i:]
		return i, nil
	}
	messageType, bytes, err := p.c.ReadMessage()
	if err != nil {
		return 0, err
	}
	if messageType == websocket.BinaryMessage {
		decrypt, _ := p.cipher.Decrypt(bytes)
		i := copy(b, decrypt)
		if len(decrypt) > i {
			p.buff = decrypt[i:]
		}
		return i, nil
	}
	return 0, io.EOF
}

func (p ProxyClient) Write(b []byte) (int, error) {
	encrypt, err2 := p.cipher.Encrypt(b)
	if err2 != nil {
		return 0, err2
	}
	err := p.c.WriteMessage(websocket.BinaryMessage, encrypt)
	return len(encrypt), err
}

func (p ProxyClient) sendAnyObject(o any) {
	marshal, _ := json.Marshal(o)
	err := p.sendTextMsg(string(marshal))
	if err != nil {
		panic(err)
	}
}

func (p ProxyClient) sendTextMsg(a string) error {
	encrypt := utils.AesInstant.Encrypt(a)
	return p.c.WriteMessage(websocket.TextMessage, []byte(encrypt))
}

func (p ProxyClient) readTextMsg(a any) {
	messageType, bytes, err := p.c.ReadMessage()
	if messageType != websocket.TextMessage {
		panic("msg type error")
	}
	if err != nil {
		panic(err)
	}
	str := utils.AesInstant.DecryptStr(string(bytes))
	err = json.Unmarshal([]byte(str), a)
	if err != nil {
		panic(err)

	}
}

type publicKey struct {
	Command   int    `json:"command"`
	PublicKey string `json:"publicKey"`
}

// 这里是请求信息
type sendConnectInfo struct {
	Command  int    `json:"command"`
	Method   string `json:"method"`
	Username string `json:"username"`
	Password string `json:"password"`
	Random   string `json:"random"`
}

func (s sendConnectInfo) GetDecodeRandom() string {
	decryptStr := utils.AesInstant.DecryptStr(s.Random)
	str, err := utils.RsaInstant.RsaDecryptStr(decryptStr)
	if err != nil {
		panic(err)
	}
	return str
}

type connectInfo struct {
	Command int    `json:"command"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Network string `json:"network"`
}

func (c connectInfo) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type command struct {
	Command int `json:"command"`
}
