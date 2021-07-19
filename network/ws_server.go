package network

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kolterr/im/log"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	WsPath = "/ws"
)

type WsServer struct {
	opts         Options
	listener     net.Listener
	connNum      int64
	upgrader     websocket.Upgrader
	handler      Handler
	closeChan    chan struct{}
	allowOrigins []string
}

func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	if s.opts.MaxConn > 0 && atomic.LoadInt64(&s.connNum) >= int64(s.opts.MaxConn) {
		log.Error("to many connections")
		conn.Close()
		return
	}
	atomic.AddInt64(&s.connNum, 1)
	s.handler.Handle(NewWsConn(conn, s.opts.Decoder))
	atomic.AddInt64(&s.connNum, -1)
}

func NewWsServer(opts ...Option) (*WsServer, error) {
	opt := newOptions(opts...)
	listener, err := net.Listen("tcp", opt.Address)
	if err != nil {
		return nil, err
	}
	return &WsServer{opts: opt, listener: listener, upgrader: websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
		CheckOrigin: func(r *http.Request) bool {
			if opt.OriginCheck != nil {
				if err := opt.OriginCheck(r.RemoteAddr); err != nil {
					return false
				}
			}
			return true
		},
	}}, nil
}

func (s *WsServer) Start(ctx context.Context, handler Handler) error {
	closeChan := make(chan struct{}, 1)
	s.closeChan = closeChan
	go func() {
		if err := s.run(ctx, handler, closeChan); err != nil {
			log.Errorf("ws pkg start error", err)
		}
	}()
	return nil
}

func (s *WsServer) run(ctx context.Context, handler Handler, closeChan chan struct{}) error {
	httpServer := http.Server{
		Addr:      s.opts.Address,
		Handler:   s,
		TLSConfig: s.opts.TlsConfig,
	}

	startError := make(chan error, 1)
	go func(startError chan error) {
		if err := httpServer.Serve(s.listener); err != nil {
			startError <- err
		}
	}(startError)
	for {
		select {
		case <-closeChan:
			ctx := context.Background()
			httpServer.Shutdown(ctx)
		case err := <-startError:
			fmt.Println(err)
		}

	}
}

func (s *WsServer) Close() {
	s.listener.Close()
	close(s.closeChan)
}
