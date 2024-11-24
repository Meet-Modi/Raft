package consensus

import (
	"Raft/logging"
)

type RaftState struct {
	currentTerm uint64              // Persistent state on all servers
	votedFor    string              // Persistent state on all servers
	LeaderId    string              // Persistent state on all servers
	LogService  *logging.LogService // Persistent state on all servers
	CommitIndex uint64              // Volatile state on all servers
	LastApplied uint64              // Volatile state on all servers
	nextIndex   map[string]uint64   // Volatile state on leaders
	matchIndex  map[string]uint64   // Volatile state on leaders
}

func InitialiseRaftState() (*RaftState, error) {
	log, _ := logging.NewLogService()

	return &RaftState{
		currentTerm: 0,
		votedFor:    "",
		LeaderId:    "",
		LogService:  log,
		CommitIndex: 0,
		LastApplied: 0,
		nextIndex:   make(map[string]uint64),
		matchIndex:  make(map[string]uint64),
	}, nil

}
