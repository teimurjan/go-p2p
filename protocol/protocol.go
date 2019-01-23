package protocol

// NewPeerCode is a code for notification about a new client connected
const NewPeerCode = 1

// ExitPeerCode is a code for notification about a client exited
const ExitPeerCode = 2

// Request is a protocol request model
type Request struct {
	Code     int64
	Data     string
	MetaData string
}

// Response is a protocol response model
type Response struct {
	Code     int64
	Data     string
	MetaData string
}
