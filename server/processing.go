package server

import (
	"log"
	"path"

	"github.com/teimurjan/go-p2p/fileutils"
	"github.com/teimurjan/go-p2p/protocol"
)

func getFilePath(filename string) string {
	return path.Join(configuration.FileSourceDir, filename)
}

// checkFile checking if file exists
func checkFile(requestInfo *protocol.RequestInfo) protocol.Response {
	filePath := getFilePath(requestInfo.FileName)
	response := protocol.Response{}
	responseInfo := protocol.ResponseInfo{}

	if !fileutils.IsFileExists(filePath) {
		response.Status = protocol.FileNotExistStatus
	} else {
		foundFileHash, error := fileutils.GetFileHash(filePath)
		if error != nil {
			response.Status = protocol.ServerSideError
			log.Printf("Error was occurred while evaluating hash %v", error)
		} else {
			response.Status = protocol.FileExistStatus
			responseInfo.FileName = requestInfo.FileName
			responseInfo.FileHash = foundFileHash
			responseInfo.FileSize = fileutils.GetFileSize(filePath)
			response.Info = responseInfo
		}
	}
	return response
}

// retrieveChunk tries to retrieve file chunk
func retrieveChunk(request *protocol.Request) protocol.Response {
	response := checkFile(&request.Info)
	if response.Status == protocol.FileExistStatus {
		filePath := path.Join(configuration.FileSourceDir, request.Info.FileName)
		bytes, error := fileutils.GetFileChunk(filePath, request.Info.ChunkIndex, request.Info.ChunkSize)
		if error != nil {
			response.Status = protocol.ChunkNotSentStatus
			log.Printf("Error was occured while retrieving file chunk %v", error)
		} else {
			responseInfo := protocol.ResponseInfo{}
			responseInfo.FileName = request.Info.FileName
			responseInfo.ChunkSize = int64(len(bytes))
			responseInfo.ChunkIndex = request.Info.ChunkIndex
			response.Status = protocol.ChunkSentStatus
			response.Data = bytes
			response.Info = responseInfo
		}
	}
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
