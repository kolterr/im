package network

import (
	"fmt"
	"testing"
)

type testHandler struct {
}

func (t *testHandler) Handle(conn Conn) {
	fmt.Println("new client connection:", conn.RemoteAddr().String())
	for {
		data, err := conn.ReadMsg()
		if err != nil {
			fmt.Println("read message error:", err)
			break
		}
		fmt.Println(string(data), err)
	}
}

func (t *testHandler) Close() error {
	return nil
}

func newTcpServer(options []Option) *TcpServer {
	s, _ := NewTcpServer(options...)
	return s
}

func defaultTcpServer() *TcpServer {
	options := []Option{
		Address(addr),
		MaxConn(10),
	}
	return newTcpServer(options)
}

func TestNewTcpServer_ConnNum(t *testing.T) {
	defaultServer := defaultTcpServer()
	clientsNum := 10
	clients := newClients(clientsNum, addr)
	if err := defaultServer.Start(&testHandler{}); err != nil {
		t.Error("tcp server start error", err)
	}
	for _, client := range clients {
		if err := client.Listen(); err != nil {
			t.Error("tcp client dial error", err)
		}
	}
	if defaultServer.ConnNum() != int64(clientsNum) {
		t.Error("lose client connection")
	}
}

func TestTcpServer_MaxConn(t *testing.T) {
	defaultServer := defaultTcpServer()
	clientsNum := 11
	clients := newClients(clientsNum, addr)
	if err := defaultServer.Start(&testHandler{}); err != nil {
		t.Error("tcp server start error", err)
	}
	for k, client := range clients {
		err := client.Listen()
		if k == 11 && err == nil {
			t.Error("maxConn not effective")
		}
	}
}
