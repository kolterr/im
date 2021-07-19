package network

import (
	"github.com/gorilla/websocket"
	"net"
)

type WsConn struct {
	conn    *websocket.Conn
	decoder Decoder
}

func (w *WsConn) Name() string {
	return Ws
}

func (w *WsConn) WriteMsg(args ...[]byte) error {
	return nil
}

func (w *WsConn) ReadMsg() ([]byte, error) {
	_, b, err := w.conn.ReadMessage()
	return b, err
}

func (w *WsConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *WsConn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *WsConn) Close() error {
	return nil
}

func NewWsConn(conn *websocket.Conn, decoder Decoder) *WsConn {
	return &WsConn{conn: conn, decoder: decoder}
}
