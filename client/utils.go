package client

import (
	"net"

	"github.com/teimurjan/go-p2p/protocol"
	"github.com/teimurjan/go-p2p/utilTypes"
)

// ChunkSize is a default size chunk in bytes
const ChunkSize = 1024

func getActivePeers(peers utilTypes.UDPAddrsArray) (utilTypes.UDPAddrsArray, error) {
	peersWithFile := utilTypes.NewUDPAddrsArray()

	for _, peer := range peers {
		request := &protocol.Request{Code: protocol.CheckFileCode}

		response, err := process(peer, request)

		if err != nil {
			return nil, err
		} else {

			if response.Status == protocol.FileExistStatus {
				peersWithFile.Add(peer)
			}
		}

	}

	return peersWithFile, nil
}

func getFileInfo(peer *net.UDPAddr) *protocol.ResponseInfo {
	request := &protocol.Request{Code: protocol.CheckFileCode}

	response, err := process(peer, request)

	if err != nil {
		return nil
	}

	return &response.Info
}
