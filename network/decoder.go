package network

import (
	"encoding/binary"
	"io"
	"net"
)

type Decoder interface {
	Read(conn net.Conn) ([]byte, error)
	Write(conn net.Conn, buf ...[]byte) error
}


const (
	maxMsgSize = 4096
)

// --------------
// | len | data |
// --------------
type defaultDecoder struct {
	lenMsgLen    int
	minMsgSize   uint32
	maxMsgSize   uint32 // 消息最大长度
	littleEndian bool
}

var DefaultDecoder = &defaultDecoder{maxMsgSize: maxMsgSize, minMsgSize: 1, lenMsgLen: 2}

func (d *defaultDecoder) Read(conn net.Conn) ([]byte, error) {
	var b [4]byte
	// 消息的长度
	bufMsgLen := b[:d.lenMsgLen]
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}
	var msgLen uint32
	switch d.lenMsgLen {
	case 1:
		msgLen = uint32(bufMsgLen[0])
	case 2:
		if d.littleEndian {
			msgLen = uint32(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = uint32(binary.BigEndian.Uint16(bufMsgLen))
		}
	case 4:
		if d.littleEndian {
			msgLen = binary.LittleEndian.Uint32(bufMsgLen)
		} else {
			msgLen = binary.BigEndian.Uint32(bufMsgLen)
		}
	}
	if msgLen > d.maxMsgSize {
		return nil, MsgTooLong
	} else if msgLen < d.minMsgSize {
		return nil, MsgTooShort
	}
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, err
	}
	return msgData, nil
}

func (d *defaultDecoder) Write(conn net.Conn, args ...[]byte) error {
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}
	// check len
	if msgLen > d.maxMsgSize {
		return MsgTooLong
	} else if msgLen < d.minMsgSize {
		return MsgTooShort
	}
	msg := make([]byte, uint32(d.lenMsgLen)+msgLen)
	switch d.lenMsgLen {
	case 1:
		msg[0] = byte(msgLen)
	case 2:
		if d.littleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case 4:
		if d.littleEndian {
			binary.LittleEndian.PutUint32(msg, msgLen)
		} else {
			binary.BigEndian.PutUint32(msg, msgLen)
		}
	}
	l := d.lenMsgLen
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}
	_, err := conn.Write(msg)
	return err
}
