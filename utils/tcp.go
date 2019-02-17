package utils

import (
	"encoding/json"
	"net"

	"github.com/teimurjan/go-p2p/protocol"
)

func readBytes(conn net.Conn) (int, []byte, error) {
	inputBytes := make([]byte, 4096)
	length, err := conn.Read(inputBytes)
	if err != nil {
		return 0, nil, err
	}

	return length, inputBytes, nil
}

// ReadRequestFromTCP reads request from TCP
func ReadRequestFromTCP(conn net.Conn) (*protocol.Request, error) {
	length, inputBytes, err := readBytes(conn)
	if err != nil {
		return nil, err
	}

	var request protocol.Request
	err = json.Unmarshal(inputBytes[:length], &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

// ReadResponseFromTCP reads request from TCP
func ReadResponseFromTCP(conn net.Conn) (*protocol.Response, error) {
	length, inputBytes, err := readBytes(conn)
	if err != nil {
		return nil, err
	}

	var response protocol.Response
	err = json.Unmarshal(inputBytes[:length], &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
