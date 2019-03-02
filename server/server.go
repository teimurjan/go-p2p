package server

import (
	"encoding/json"
	"net"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-p2p/fileutils"
	"github.com/teimurjan/go-p2p/protocol"
	"github.com/teimurjan/go-p2p/utils"
)

// Server is a P2P server interface
type Server interface {
	Start()
	listenTCP()
	acceptConnections(l net.Listener)
	handleRequest(request *protocol.Request) *protocol.Response
}

type server struct {
	port          string
	fileSourceDir string
	logger        *logrus.Logger
}

// NewServer creates new server instance
func NewServer(port string, fileSourceDir string, logger *logrus.Logger) Server {
	if !fileutils.IsFileExists(fileSourceDir) {
		err := os.Mkdir(fileSourceDir, os.ModePerm)
		if err != nil {
			logger.Printf("fatal: Error was occurred while trying to create directory %v\n", err)
			os.Exit(1)
		}
	}

	return &server{port, fileSourceDir, logger}
}

func (s *server) Start() {
	s.listenTCP()
}

func (s *server) listenTCP() {
	l, _ := net.Listen("tcp", ":"+s.port)
	defer l.Close()

	s.logger.Println("Server has started.")

	for {
		s.acceptConnections(l)
	}
}

func (s *server) acceptConnections(l net.Listener) {
	conn, err := l.Accept()
	if err != nil {
		s.logger.Errorf("Cannot accept connection %s", err)
		os.Exit(1)
	}
	defer conn.Close()

	request, err := utils.ReadRequestFromTCP(conn)
	if err != nil {
		s.logger.Error(err)
		return
	}

	response := s.handleRequest(request)
	json, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err)
		return
	}
	conn.Write(json)
}

func (s *server) handleRequest(request *protocol.Request) *protocol.Response {
	var response protocol.Response

	switch request.Code {
	case protocol.CheckFileCode:
		response = s.checkFile(request)
	case protocol.GetChunkCode:
		response = s.retrieveChunk(request)
	}

	return &response
}

func (s *server) checkFile(request *protocol.Request) protocol.Response {
	filePath := path.Join(s.fileSourceDir, request.Info.FileName)

	response := protocol.Response{}
	responseInfo := protocol.ResponseInfo{}

	if !fileutils.IsFileExists(filePath) {
		response.Status = protocol.FileNotExistStatus
	} else {
		foundFileHash, err := fileutils.GetFileHash(filePath)
		if err != nil {
			response.Status = protocol.ServerSideError
			s.logger.Error(err)
		} else {
			response.Status = protocol.FileExistStatus
			responseInfo.FileName = request.Info.FileName
			responseInfo.FileHash = foundFileHash
			responseInfo.FileSize = fileutils.GetFileSize(filePath)
			response.Info = responseInfo
		}
	}
	return response
}

func (s *server) retrieveChunk(request *protocol.Request) protocol.Response {
	response := s.checkFile(request)
	if response.Status == protocol.FileExistStatus {
		filePath := path.Join(s.fileSourceDir, request.Info.FileName)
		bytes, err := fileutils.GetFileChunk(filePath, request.Info.ChunkIndex, request.Info.ChunkSize)
		if err != nil {
			response.Status = protocol.ChunkNotSentStatus
			s.logger.Error(err)
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
