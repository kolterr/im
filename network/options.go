package network

import (
	"crypto/tls"
	"github.com/kolterr/hellochat/pkg/broker"
	"github.com/kolterr/hellochat/pkg/discover"
)

type Options struct {
	Broker      broker.Broker
	Discover    discover.Discover
	Decoder     Decoder
	Address     string
	MaxConn     int
	OriginCheck func(addr string) error
	TlsConfig   *tls.Config
}

type Option func(options *Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Decoder: DefaultDecoder,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Broker(broker broker.Broker) Option {
	return func(options *Options) {
		options.Broker = broker
	}
}

func Discover(discover discover.Discover) Option {
	return func(options *Options) {
		options.Discover = discover
	}
}

func Decode(decoder Decoder) Option {
	return func(options *Options) {
		options.Decoder = decoder
	}
}

func Address(addr string) Option {
	return func(options *Options) {
		options.Address = addr
	}
}

func MaxConn(num int) Option {
	return func(options *Options) {
		options.MaxConn = num
	}
}

func CheckOrigin(fn func(addr string) error) Option {
	return func(options *Options) {
		options.OriginCheck = fn
	}
}

func TlsConfig(tlsConfig *tls.Config) Option {
	return func(options *Options) {
		options.TlsConfig = tlsConfig
	}
}
