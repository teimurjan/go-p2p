package models

import (
	"net"

	"github.com/teimurjan/go-p2p/protocol"
)

// Notification is a notification model
type Notification struct {
	Req      *protocol.Request
	FromAddr *net.UDPAddr
}
