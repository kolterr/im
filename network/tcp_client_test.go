package network

import (
	"fmt"
	"testing"
	"time"
)

const (
	addr = "127.0.0.1:65028"
)

func newClients(num int, addr string) []*TcpClient {
	clients := make([]*TcpClient, 0)
	for i := 1; i <= num; i++ {
		clients = append(clients, NewTcpClient(Address(addr)))
	}
	return clients
}

func TestTcpClient_Listen(t *testing.T) {
	clients := newClients(10, addr)
	for _, client := range clients {
		if err := client.Listen(); err != nil {
			t.Error("client listen error", err)
		}
	}
	time.Sleep(time.Second*5)
}

func TestTcpClient_WriteMsg(t *testing.T) {
	clients := newClients(10, addr)
	for _, client := range clients {
		if err := client.Listen(); err != nil {
			t.Error("client listen error", err)
		}
		fmt.Println(client.WriteMsg([]byte("Hello world!")))
	}
	time.Sleep(time.Minute)
}
