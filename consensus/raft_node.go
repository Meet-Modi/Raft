package consensus

import (
	"Raft/logging"
)

type RaftNode struct {
	currentTerm int
	votedFor    string
	isLeader    bool
	LogService  *logging.LogService
	commitIndex int
	lastApplied int
	nextIndex   map[string]int
	matchIndex  map[string]int

	networkNodes map[string]string //this stores all the Nodes IP
}

func InitialiseRaftNode() (*RaftNode, error) {
	log, _ := logging.NewLogService()

	return &RaftNode{
		currentTerm: 0,
		votedFor:    "",
		isLeader:    false, // First initialise with follower mode
		LogService:  log,
		commitIndex: 0,
		lastApplied: 0,
	}, nil

}

func (r *RaftNode) getCurrentState() map[string]interface{} {
	return map[string]interface{}{
		"isLeader":    r.isLeader,
		"commitIndex": r.commitIndex,
		"votedFor":    r.votedFor,
		"lastApplied": r.lastApplied,
	}
}
