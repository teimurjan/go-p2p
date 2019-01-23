package models

import (
	"net"

	"github.com/teimurjan/go-p2p/protocol"
)

type Notification struct {
	Req      *protocol.Request
	Res      *protocol.Response
	FromAddr *net.UDPAddr
}

// NewReceievedNotification creates new notification instance which is received
func NewReceievedNotification(res *protocol.Response, fromAddr *net.UDPAddr) *Notification {
	return &Notification{
		nil,
		res,
		fromAddr,
	}
}

// NewNotification creates new notification instance to send
func NewNotification(req *protocol.Request, fromAddr *net.UDPAddr) *Notification {
	return &Notification{
		req,
		nil,
		fromAddr,
	}
}
