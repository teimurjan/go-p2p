package fileutils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// IsFileExists indicates whether a file exists or not
func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetFileSize returns file size
func GetFileSize(path string) int64 {
	f, _ := os.Stat(path)
	return f.Size()
}

// GetFileHash returns file hash
func GetFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// GetFileChunk returns file chunk
func GetFileChunk(path string, chunkIndex int64, chunkSize int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileSize := GetFileSize(path)

	chunkBytesStartPosition := chunkSize * chunkIndex

	if chunkBytesStartPosition >= fileSize {
		return nil, nil
	}
	if chunkBytesStartPosition+chunkSize > fileSize {
		chunkSize = fileSize - chunkBytesStartPosition
	}
	chunkBuffer := make([]byte, chunkSize)
	_, err = f.ReadAt(chunkBuffer, chunkBytesStartPosition)

	if err != nil {
		return nil, err
	}

	return chunkBuffer, nil
}
