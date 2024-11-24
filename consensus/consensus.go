package consensus

import (
	"Raft/config"
	"Raft/discovery"
	"Raft/logging"
	pb "Raft/proto/consensus"
	"Raft/utils"
	"context"
	"log"
	"sync"
)

type RaftState struct {
	mu               sync.Mutex
	Id               string                      // Unique ID of the node
	currentTerm      uint64                      // Persistent state on all servers
	votedFor         string                      // Persistent state on all servers
	LeaderId         string                      // Persistent state on all servers
	LogService       *logging.LogService         // Persistent state on all servers
	CommitIndex      uint64                      // Volatile state on all servers
	LastApplied      uint64                      // Volatile state on all servers
	nextIndex        map[string]uint64           // Volatile state on leaders, reinitialized after election
	matchIndex       map[string]uint64           // Volatile state on leaders, reinirialized after election
	Mode             string                      // Leader, Follower, Candidate
	DiscoveryService *discovery.DiscoveryService // Reference to the discovery service
	pb.UnimplementedConsensusServiceServer
}

// This function initialises the Raft state
// If the node is a boot node, it is initialised as a leader
// else it is initialised as a follower first.
func InitialiseRaftState() (*RaftState, error) {
	logService, _ := logging.NewLogService()

	var id, initMode, LeaderId string
	if config.IsBootNode {
		id = utils.GenerateRaftPeerId(true)
		initMode = "Leader"
		LeaderId = utils.GenerateRaftPeerId(true)
	} else {
		id = utils.GenerateRaftPeerId(false)
		initMode = "Follower"
		LeaderId = ""
	}

	discoveryService := discovery.NewDiscoveryService(id)

	return &RaftState{
		Id:               id,
		currentTerm:      0,
		votedFor:         "",
		LeaderId:         LeaderId,
		LogService:       logService,
		CommitIndex:      0,
		LastApplied:      0,
		nextIndex:        make(map[string]uint64),
		matchIndex:       make(map[string]uint64),
		Mode:             initMode,
		DiscoveryService: discoveryService,
	}, nil
}

type AppendEntriesArgs struct {
	PeerId       string
	Term         uint64
	LeaderId     string // Leader's ID should be same as PeerID
	PrevLogIndex uint64
	prevLogTerm  uint64
	Entries      []logging.LogEntry
	LeaderCommit uint64
}

type AppendEntriesReply struct {
	Term    uint64
	Success bool
}

type RequestVoteArgs struct {
	Term         uint64
	CandidateId  string
	LastLogIndex uint64
	LastLogTerm  uint64
}

type RequestVoteReply struct {
	Term        uint64
	VoteGranted bool
}

// This function can only be invoked by a leader
func (rs *RaftState) SendAppendEntriesRPC() error {
	return nil
}

func (rs *RaftState) AppendEntries(ctx context.Context, in *pb.AppendEntriesRequest) (reply *pb.AppendEntriesResponse, err error) {
	if in.LeaderId != rs.LeaderId {
		log.Fatalf("PeerId should be same as leader ID. Expected leaderId %s %s", in.LeaderId, rs.LeaderId)
	}

	rs.mu.Lock()
	defer rs.mu.Unlock()
	reply = &pb.AppendEntriesResponse{}

	// Run validation checks
	// 1. Reply false if term < currentTerm (§5.1)
	if in.Term < rs.currentTerm {
		reply.Term = rs.currentTerm
		reply.Success = false
		return reply, nil
	}

	// 2. Reply false if log doesn’t contain an entry at prevLogIndex
	// whose term matches prevLogTerm (§5.3)
	if rs.LogService.Logs[in.PrevLogIndex].Term != in.PrevLogTerm {
		reply.Term = rs.currentTerm
		reply.Success = false
		return reply, nil
	}

	// 3.  If an existing entry conflicts with a new one (same index but different terms),
	// delete the existing entry and all that follow it (§5.3)
	for _, entry := range in.Logs {
		if entry.Index < uint64(len(rs.LogService.Logs)) && rs.LogService.Logs[entry.Index].Term != entry.Term {
			rs.LogService.Logs = rs.LogService.Logs[:entry.Index]
			reply.Term = rs.currentTerm
			reply.Success = false
			return reply, nil
		}
	}

	// 4.  Append any new entries not already in the log.
	// TODO: Can optimise this in one check.
	for _, entry := range in.Logs {
		if entry.Index >= uint64(len(rs.LogService.Logs)) {
			rs.LogService.Logs = append(rs.LogService.Logs, logging.LogEntry{Term: entry.Term, Index: entry.Index, Command: entry.Command})
		}
	}

	// 5.  If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
	if in.LeaderCommit > rs.CommitIndex {
		rs.CommitIndex = utils.Min(in.LeaderCommit, uint64(len(rs.LogService.Logs)-1))
	}

	reply.Term = rs.currentTerm
	reply.Success = true
	return reply, nil
}
