package consensus

import (
	"Raft/config"
	"Raft/discovery"
	"Raft/logging"
	pb "Raft/proto/consensus"
	"Raft/utils"
	"fmt"
	"sync"
	"time"
)

type RaftState struct {
	mu                sync.Mutex
	Id                string                      // Unique ID of the node
	currentTerm       uint64                      // Persistent state on all servers
	votedFor          string                      // Persistent state on all servers
	LeaderId          string                      // Persistent state on all servers
	LogService        *logging.LogService         // Persistent state on all servers
	CommitIndex       uint64                      // Volatile state on all servers
	LastApplied       uint64                      // Volatile state on all servers
	nextIndex         map[string]uint64           // Volatile state on leaders, reinitialized after election
	matchIndex        map[string]uint64           // Volatile state on leaders, reinirialized after election
	Mode              string                      // Leader, Follower, Candidate
	DiscoveryService  *discovery.DiscoveryService // Reference to the discovery service
	ElectionTimeoutMs int                         // Election timeout in milliseconds
	LastHeartbeat     time.Time                   // Last heartbeat received
	pb.UnimplementedConsensusServiceServer
}

// This function initialises the Raft state
// If the node is a boot node, it is initialised as a leader
// else it is initialised as a follower first.
func InitialiseRaftState() (*RaftState, error) {
	var id, LeaderId string
	id = utils.GenerateRaftPeerId()

	logService, _ := logging.NewLogService(id)
	discoveryService := discovery.NewDiscoveryService(id)
	initMode := "Follower"

	rs := &RaftState{
		Id:                id,
		currentTerm:       0,
		votedFor:          "",
		LeaderId:          LeaderId,
		LogService:        logService,
		CommitIndex:       0,
		LastApplied:       0,
		nextIndex:         make(map[string]uint64),
		matchIndex:        make(map[string]uint64),
		Mode:              initMode,
		DiscoveryService:  discoveryService,
		ElectionTimeoutMs: config.ElectionTimeoutMs,
	}
	return rs, nil
}

func (rs *RaftState) ShutdownHandling() error {
	err := rs.LogService.ShutdownHandling()
	if err != nil {
		return fmt.Errorf("failed to close log service: %v", err)
	}
	return nil
}
