package logging

import (
	"fmt"
	"os"
	"sync"
)

type LogEntry struct {
	Term    uint64
	Index   int64
	Command string
}

type LogService struct {
	mu           sync.Mutex
	Logs         []LogEntry
	LastLogIndex int64
	LogFile      *os.File
}

func NewLogService(id string) (*LogService, error) {

	logFileName := id + ".log"
	file, err := os.Create("/app/tmp/" + logFileName)
	if err != nil {
		panic("Failed to create file: " + err.Error())
	}

	return &LogService{
		Logs:         []LogEntry{},
		LogFile:      file,
		LastLogIndex: -1,
	}, nil
}

func (ls *LogService) PersistLogEntry(entry LogEntry) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.Logs = append(ls.Logs, entry)
	ls.LastLogIndex++
	ls.LogFile.WriteString(fmt.Sprintf("%d %d %s\n", entry.Term, entry.Index, entry.Command))
}

func (ls *LogService) ShutdownHandling() error {
	err := ls.LogFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close log file: %v", err)
	}
	return nil
}
