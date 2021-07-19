package network

import (
	"github.com/kolterr/lib/wait"
	"io"
	"net"
	"time"
)

type TcpClient struct {
	opt       Options
	conn      net.Conn
	closeChan chan struct{}
	OnMessage func(buf []byte)
	wg        wait.Wait
}

func NewTcpClient(opts ...Option) *TcpClient {
	opt := newOptions(opts...)
	return &TcpClient{opt: opt, closeChan: make(chan struct{}, 1)}
}

func (c *TcpClient) dial() error {
	conn, err := net.Dial("tcp", c.opt.Address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *TcpClient) WriteMsg(args ...[]byte) error {
	return c.opt.Decoder.Write(c.conn, args...)
}

func (c *TcpClient) Listen() error {
	if err := c.dial(); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-c.closeChan:
				return
			default:
				// Waiting for the server response
				data, err := c.opt.Decoder.Read(c.conn)
				if err != nil {
					if err == io.EOF {
						break
					}
				}
				if c.OnMessage != nil {
					c.wg.Add(1)
					go func() {
						c.OnMessage(data)
						c.wg.Done()
					}()
				}
			}

		}
	}()
	return nil
}

func (c *TcpClient) Close() {
	c.closeChan <- struct{}{}
	c.conn.Close()
	c.wg.WaitWithTimeout(time.Second * 10)
}
