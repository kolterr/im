package network

import (
	"github.com/kolterr/hellochat/pkg/log"
	libAtomic "github.com/kolterr/lib/atomic"
	"net"
	"sync"
	"sync/atomic"
)

type Handler interface {
	Handle(conn Conn)
	Close() error
}

type TcpServer struct {
	opts      Options
	listener  net.Listener
	closeChan chan struct{}
	closed    libAtomic.Bool
	wg        sync.WaitGroup
	connNum   *int64
}

func NewTcpServer(opts ...Option) (*TcpServer, error) {
	opt := newOptions(opts...)
	listener, err := net.Listen("tcp", opt.Address)
	if err != nil {
		return nil, err
	}
	return &TcpServer{listener: listener, opts: opt, connNum: new(int64)}, nil
}

func (s *TcpServer) Start(handler Handler) error {
	closeChan := make(chan struct{}, 1)
	s.closeChan = closeChan
	go func() {
		if err := s.Listen(handler, closeChan); err != nil {
			log.Errorf("tcp pkg start error", err)
		}
	}()
	return nil
}

func (s *TcpServer) ConnNum() int64 {
	return *s.connNum
}

func (s *TcpServer) Close() {
	if s.closed.Get() {
		return
	}
	close(s.closeChan)
	s.wg.Wait()
	s.closed.Set(true)
}

func (s *TcpServer) IsClosed() bool {
	return s.closed.Get()
}

func (s *TcpServer) Listen(handler Handler, closeChan chan struct{}) error {
	go func() {
		<-closeChan
		if err := s.listener.Close(); err != nil {
			log.Errorf("handler close error", err)
		}
		if err := handler.Close(); err != nil {
			log.Errorf("handler close error", err)
		}
	}()
	defer func() {
		if err := s.listener.Close(); err != nil {
			log.Errorf("handler close error", err)
		}
		if err := handler.Close(); err != nil {
			log.Errorf("handler close error", err)
		}
	}()
	var wg sync.WaitGroup
	log.Infof("Tcp server start accept, address is:%s", s.listener.Addr().String())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			break
		}
		if s.opts.OriginCheck != nil {
			if err := s.opts.OriginCheck(conn.RemoteAddr().String()); err != nil {
				log.Errorf("forbidden origin:", conn.RemoteAddr().String())
				conn.Close()
				continue
			}
		}
		if s.opts.MaxConn > 0 && atomic.LoadInt64(s.connNum) >= int64(s.opts.MaxConn) {
			log.Errorf("to many connections")
			conn.Close()
			continue
		}
		atomic.AddInt64(s.connNum, 1)
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.Handle(NewTcpConn(conn, s.opts.Decoder))
			atomic.AddInt64(s.connNum, -1)
		}()
	}
	wg.Wait()
	return nil
}
