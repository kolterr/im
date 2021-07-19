package network

import "errors"

var (
	MsgTooLong  = errors.New("message too long")
	MsgTooShort = errors.New("message too short")
)
