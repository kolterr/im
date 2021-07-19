package network

import "net"

type TcpConn struct {
	conn    net.Conn
	decoder Decoder
}

func (t *TcpConn) Name() string {
	return Tcp
}

func (t *TcpConn) WriteMsg(args ...[]byte) error {
	return t.decoder.Write(t.conn, args...)
}

func (t *TcpConn) ReadMsg() ([]byte, error) {
	return t.decoder.Read(t.conn)
}

func (t *TcpConn) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}

func (t *TcpConn) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *TcpConn) Close() error {
	return nil
}

func NewTcpConn(conn net.Conn, decoder Decoder) *TcpConn {
	return &TcpConn{conn: conn, decoder: decoder}
}
