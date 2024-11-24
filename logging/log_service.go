package logging

import (
	"os"
)

type LogEntry struct {
	Term    uint64
	Index   uint64
	Command string
}

type LogService struct {
	Logs             []LogEntry
	LastAppliedIndex uint64 // 0-based indexing
	CommitIndex      uint64 // 0-based indexing
	LogFile          *os.File
}

func NewLogService() (*LogService, error) {

	logFileName := "1.txt"
	file, err := os.Create("tmp/" + logFileName)
	if err != nil {
		panic("Failed to create file: " + err.Error())
	}

	// defer file.Close()

	return &LogService{
		LastAppliedIndex: 0,
		CommitIndex:      0,
		Logs:             []LogEntry{},
		LogFile:          file,
	}, nil
}
