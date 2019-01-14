package protocol

import "net"

// ConnectedID is an ID for notification about new client connection
const ConnectedID = 1

// Notification is a notification model
type Notification struct {
	ID       int64
	FromAddr *net.UDPAddr
}
