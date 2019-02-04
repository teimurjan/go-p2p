package server

import (
	"encoding/json"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-p2p/protocol"
)

// Server is a P2P server interface
type Server interface {
	Start()
	listenTCP()
	acceptConnections(l net.Listener)
	handleRequest(request *protocol.Request) *protocol.Response
}

type server struct {
	port   string
	logger *logrus.Logger
}

// NewServer creates new server instance
func NewServer(port string, logger *logrus.Logger) Server {
	return &server{port, logger}
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

	inputBytes := make([]byte, 4096)
	length, err := conn.Read(inputBytes)
	if err != nil {
		s.logger.Errorf("Cannot read from connection. %s", err)
		return
	}

	var request protocol.Request
	err = json.Unmarshal(inputBytes[:length], &request)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	response := s.handleRequest(&request)
	json, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err)
		return
	}
	conn.Write(json)
}

func (s *server) handleRequest(request *protocol.Request) *protocol.Response {
	return &protocol.Response{}
}
