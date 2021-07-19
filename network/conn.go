package network

import "net"

const (
	Tcp = "tcp"
	Ws = "ws"
)

// the interface for tcp of websocket
type Conn interface {
	Name() string
	WriteMsg(args ...[]byte) error
	ReadMsg() ([]byte, error)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()error
}
