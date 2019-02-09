package server

import (
	"path"

	"github.com/teimurjan/go-p2p/fileutils"
	"github.com/teimurjan/go-p2p/protocol"
)

// CheckFile checking of file existence and builds response using protocol.Response
func checkFile(requestInfo *protocol.RequestInfo) protocol.Response {
	filePath := path.Join(configuration.FileSourceDir, requestInfo.FileName)
	response := protocol.Response{}
	responseInfo := protocol.ResponseInfo{}

	if !fileutils.IsFileExists(filePath) {
		response.Status = protocol.FileNotExistStatus
	} else {
		foundFileHash, error := fileutils.GetFileHash(filePath)
		if error != nil {
			response.Status = protocol.ServerSideError
		} else {
			response.Status = protocol.FileExistStatus
			responseInfo.FileName = requestInfo.FileName
			responseInfo.FileHash = foundFileHash
			response.Info = responseInfo
		}
	}
	return response
}

func retrieveChunk(request *protocol.Request) protocol.Response {
	response := protocol.Response{}
	return response
}

func process(request *protocol.Request) protocol.Response {

	var response protocol.Response

	switch request.Code {
	case protocol.CheckFileCode:
		response = checkFile(&request.Info)
	case protocol.GetChunkCode:
		response = retrieveChunk(request)
	}
	return response
}
