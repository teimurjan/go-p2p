package protocol

/*
	This group of constants defines actions codes
*/
const (
	// NewPeerCode is a code for notification about a new client connected
	NewPeerCode = 1
	// ExitPeerCode is a code for notification about a client exited
	ExitPeerCode = 2
	// CheckFileCode is a code for checking file availability
	CheckFileCode = 3
	// GetChunkCode is a code for getting a chunk
	GetChunkCode = 4
)

/*
	This group of codes defines status of responses from server
*/
const (
	// Any server-side error
	ServerSideError = -1
	// FileNotExistStatus is a status when requested file does not exist
	FileNotExistStatus = 0
	// FileExistStatus is a status when requested file exists
	FileExistStatus = 1
	// ChunkNotSentStatus is a status when requested chunk could not be sent
	ChunkNotSentStatus = 2
	// ChunkSentStatus is a status when requested chunk was sent
	ChunkSentStatus = 3
)

// RequestInfo is meta data of a request
type RequestInfo struct {
	FileName   string
	ChunkIndex int64
}

// Request is a protocol request model
type Request struct {
	Code int64
	Data string
	Info RequestInfo
}

// ResponseInfo is meta data of a response
type ResponseInfo struct {
	FileName                string
	FileHash                string
	ChunkSize               int64
	ChunkBytesStartPosition int64
}

// Response is a protocol response model
type Response struct {
	Status int64
	Data   string
	Info   ResponseInfo
}
