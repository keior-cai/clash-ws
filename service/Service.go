package service

import (
	"bufio"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"net"
)

type Service struct {
	net.Conn
}

func (s Service) Connect() net.Conn {
	addr, err := s.readAddr()
	if err != nil {
		panic(err)
	}
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	return dial
}

func (s Service) readAddr() (string, error) {
	reader := bufio.NewReader(s)
	one, err := reader.ReadByte()
	if err != nil {
		return "", nil
	}
	tow, err := reader.ReadByte()
	if err != nil {
		return "", nil
	}
	bytes := []byte{one, tow}
	size := int(binary.BigEndian.Uint16(bytes))
	l := make([]byte, size)
	_, err = reader.Read(l)
	if err != nil {
		return "", err
	}
	addr := string(l)
	logrus.Debugf("addr = %s", addr)
	return addr, nil
}

type Client struct {
	net.Conn
}

func (c Client) WriteAddr(addr string) error {
	size := len(addr)
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(size))
	_, err := c.Write(bytes)
	if err != nil {
		return err
	}
	_, err = c.Write([]byte(addr))
	return err
}
