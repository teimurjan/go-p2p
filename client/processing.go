package client

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/teimurjan/go-p2p/protocol"
	"github.com/teimurjan/go-p2p/utils"
)

func process(peer *net.UDPAddr, request *protocol.Request) (*protocol.Response, error) {
	conn, err := getConnWithPeer(peer)
	if err != nil {
		return nil, fmt.Errorf("Connection with %s cannot be established", peer.IP.String())
	}

	defer conn.Close()

	json, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	conn.Write(json)

	response, err := utils.ReadResponseFromTCP(conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getConnWithPeer(peer *net.UDPAddr) (net.Conn, error) {
	tcpAddr := peer.String() + ":" + string(peer.Port)
	return net.DialTimeout("tcp", tcpAddr, time.Second*2)
}
