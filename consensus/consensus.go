package consensus

import (
	"Raft/logging"
	"Raft/utils"
	"log"
	"sync"
)

type RaftState struct {
	mu          sync.Mutex
	currentTerm uint64              // Persistent state on all servers
	votedFor    string              // Persistent state on all servers
	LeaderId    string              // Persistent state on all servers
	LogService  *logging.LogService // Persistent state on all servers
	CommitIndex uint64              // Volatile state on all servers
	LastApplied uint64              // Volatile state on all servers
	nextIndex   map[string]uint64   // Volatile state on leaders, reinitialized after election
	matchIndex  map[string]uint64   // Volatile state on leaders, reinirialized after election
}

func InitialiseRaftState() (*RaftState, error) {
	logService, _ := logging.NewLogService()

	return &RaftState{
		currentTerm: 0,
		votedFor:    "",
		LeaderId:    "",
		LogService:  logService,
		CommitIndex: 0,
		LastApplied: 0,
		nextIndex:   make(map[string]uint64),
		matchIndex:  make(map[string]uint64),
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

func (rs *RaftState) SendAppendEntriesRPC(peerId string, args AppendEntriesArgs, reply *AppendEntriesReply) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	// reply := &AppendEntriesReply{}

	// Response received
	if reply.Success {
		// Update nextIndex and matchIndex
		rs.nextIndex[peerId] = args.Entries[len(args.Entries)-1].Index + 1
		rs.matchIndex[peerId] = args.Entries[len(args.Entries)-1].Index
	} else {
		// If AppendEntries fails because of log inconsistency, decrement nextIndex and retry
		rs.nextIndex[peerId]--
	}

	return nil
}

// AppendEntriesReceiverImpl is the implementation of the AppendEntries RPC
func (rs *RaftState) AppendEntriesReceiverImpl(args AppendEntriesArgs, reply *AppendEntriesReply) error {
	if args.PeerId != rs.LeaderId {
		log.Fatalf("PeerId should be same as leader ID. Expected leaderId %s %s", args.PeerId, rs.LeaderId)
	}

	rs.mu.Lock()
	defer rs.mu.Unlock()

	// Run validation checks
	// 1. Reply false if term < currentTerm (§5.1)
	if args.Term < rs.currentTerm {
		reply.Term = rs.currentTerm
		reply.Success = false
		return nil
	}

	// 2. Reply false if log doesn’t contain an entry at prevLogIndex
	// whose term matches prevLogTerm (§5.3)
	if rs.LogService.Logs[args.PrevLogIndex].Term != args.prevLogTerm {
		reply.Term = rs.currentTerm
		reply.Success = false
		return nil
	}

	// 3.  If an existing entry conflicts with a new one (same index but different terms),
	// delete the existing entry and all that follow it (§5.3)
	for _, entry := range args.Entries {
		if entry.Index < uint64(len(rs.LogService.Logs)) && rs.LogService.Logs[entry.Index].Term != entry.Term {
			rs.LogService.Logs = rs.LogService.Logs[:entry.Index]
			break
		}
	}

	// 4.  Append any new entries not already in the log.
	// TODO: Can optimise this in one check.
	for _, entry := range args.Entries {
		if entry.Index >= uint64(len(rs.LogService.Logs)) {
			rs.LogService.Logs = append(rs.LogService.Logs, entry)
		}
	}

	// 5.  If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
	if args.LeaderCommit > rs.CommitIndex {
		rs.CommitIndex = utils.Min(args.LeaderCommit, uint64(len(rs.LogService.Logs)-1))
	}

	reply.Term = rs.currentTerm
	reply.Success = true
	return nil
}

func (rs *RaftState) RequestVoteReceiverImpl(args RequestVoteArgs, reply *RequestVoteReply) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// run validation checks
	// 1. Reply false if term < currentTerm (§5.1)
	if args.Term < rs.currentTerm {
		reply.Term = rs.currentTerm
		reply.VoteGranted = false
	}

	// 2.  If votedFor is null or candidateId, and candidate’s log is at
	// least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	if rs.votedFor == "" || rs.votedFor == args.CandidateId {
		// Check if candidate's log is at least as up-to-date as receiver's log
		if args.LastLogTerm >= rs.LogService.Logs[len(rs.LogService.Logs)-1].Term && args.LastLogIndex >= rs.LogService.Logs[len(rs.LogService.Logs)-1].Index {
			reply.Term = rs.currentTerm
			reply.VoteGranted = true
			rs.votedFor = args.CandidateId
			return nil
		}
	}
	return nil
}

// func (rs *RaftState) SendRequestVoteRPC(peerId string, args RequestVoteArgs, reply *RequestVoteReply) error
