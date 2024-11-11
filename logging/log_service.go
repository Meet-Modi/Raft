package logging

import (
	"errors"
	"fmt"
	"os"
)

type LogService struct {
	Logs             []LogEntry
	LastAppliedIndex int // 0-based indexing
	CommitIndex      int // 0-based indexing
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
		LastAppliedIndex: -1,
		CommitIndex:      -1,
		Logs:             []LogEntry{},
		LogFile:          file,
	}, nil
}

func (s *LogService) PersistLog(term int, startIndex int, logEntries []string) (int, int, error) {

	if startIndex != s.LastAppliedIndex+1 {
		return 0, 0, errors.New("invalid term or start index")
	}

	for _, command := range logEntries {
		entry := LogEntry{
			Term:    term,
			Index:   s.LastAppliedIndex + 1,
			Command: command,
		}
		s.Logs = append(s.Logs, entry)
		err := s.persistLogEntry(entry)
		if err != nil {
			fmt.Println("=====ERROR===== Log not persisted")
			// TODO : Handle graceful exits here!!
		}

		s.LastAppliedIndex++
		s.CommitIndex++
	}

	return s.LastAppliedIndex, s.CommitIndex, nil
}

func (s *LogService) persistLogEntry(entry LogEntry) error {
	entryString := fmt.Sprintf("%d,%d,%s\n", entry.Term, entry.Index, entry.Command)
	entryBytes := []byte(entryString)

	if _, err := s.LogFile.Write(entryBytes); err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}
