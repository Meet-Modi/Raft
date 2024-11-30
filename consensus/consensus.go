package consensus

import (
	"Raft/config"
	"Raft/discovery"
	"Raft/logging"
	pb "Raft/proto/consensus"
	"Raft/utils"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	if config.IsBootNode {
		id = utils.GenerateRaftPeerId(true)
		LeaderId = utils.GenerateRaftPeerId(true)
	} else {
		id = utils.GenerateRaftPeerId(false)
		LeaderId = ""
	}

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

	command := fmt.Sprintf("INITIALISED SERVER %v;", rs.Id)
	rs.LogService.PersistLogEntry(logging.LogEntry{Term: rs.currentTerm, Index: 0, Command: command})
	rs.LastApplied = 0

	go rs.StartPeriodicAppends()
	go rs.CheckLastHeartbeat()
	return rs, nil
}

func (rs *RaftState) ShutdownHandling() error {
	err := rs.LogService.ShutdownHandling()
	if err != nil {
		return fmt.Errorf("failed to close log service: %v", err)
	}
	return nil
}

func (rs *RaftState) StartPeriodicAppends() {
	// For sending heartbeats
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for range ticker.C {
		if rs.Mode == "Leader" {
			rs.SendAppendEntries()
		}
	}
}

func (rs *RaftState) ApplyComandToStateMachine(command string) {
	// Apply the command to the state machine
	rs.LogService.PersistLogEntry(logging.LogEntry{Term: rs.currentTerm, Index: rs.LastApplied + 1, Command: command})
	rs.LastApplied++
}

// This function can only be invoked by a leader
func (rs *RaftState) SendAppendEntries() error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	peers := rs.DiscoveryService.Peers

	for peerId, peer := range peers {
		if peerId == rs.Id {
			continue
		}

		if _, ok := rs.nextIndex[peerId]; !ok {
			rs.nextIndex[peerId] = 0 // Initialised server command will always be there
		}
		if _, ok := rs.matchIndex[peerId]; !ok {
			rs.matchIndex[peerId] = 0 // Initialised server command will always be there
		}
		// 1. Send AppendEntries RPC to all peers
		// 2. If successful, update nextIndex and matchIndex
		// 3. If not successful, decrement nextIndex and retry
		// 4. If AppendEntries fails because of log inconsistency, decrement nextIndex and retry
		// 5. If AppendEntries fails because of network failure, retry
		// 6. If AppendEntries fails because of peer failure, retry
		// 7. If AppendEntries fails because of leader failure, retry
		// 8. If AppendEntries fails because of leader crash, retry
		// 9. If AppendEntries fails because of peer crash, retry
		// 10. If AppendEntries fails because of network crash, retry
		// 11. If AppendEntries fails because of log inconsistency, retry

		// Make a gRPC call to the peer
		// If the call is successful, update the peer list
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// conn, err := grpc.DialContext(ctx, ds.Peers[peerId].URI, grpc.WithInsecure(), grpc.WithBlock())
		// conn, err := grpc.Dial(ds.Peers[peerId].URI, grpc.WithInsecure())
		conn, err := grpc.Dial(peer.URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Error in dialing the peer %s %v", peer.URI, err)
			// TODO: Can write another parellel process to delete the inactive peers
			return nil
		}
		defer conn.Close()
		client := pb.NewConsensusServiceClient(conn)

		// Prepare the AppendEntriesRequest
		request := &pb.AppendEntriesRequest{
			Term:         rs.currentTerm,
			LeaderId:     rs.LeaderId,
			PrevLogIndex: uint64(rs.matchIndex[peerId]),
			PrevLogTerm:  rs.LogService.Logs[rs.matchIndex[peerId]].Term,
			Entries:      make([]*pb.LogEntry, 0),
			LeaderCommit: rs.CommitIndex,
		}

		request.Entries = make([]*pb.LogEntry, 0)
		for i := request.PrevLogIndex + 1; i < uint64(len(rs.LogService.Logs)); i++ {
			request.Entries = append(request.Entries, &pb.LogEntry{
				Term:    rs.LogService.Logs[i].Term,
				Index:   rs.LogService.Logs[i].Index,
				Command: rs.LogService.Logs[i].Command,
			})
		}

		// log.Printf("Sending AppendEntries RPC to %s with request %v", peerId, request)

		response, err := client.AppendEntries(ctx, request)
		if err != nil {
			log.Fatalf("Error in making discovery request %v", err)
		}

		if response.Success {
			rs.nextIndex[peerId] = rs.LastApplied
			rs.matchIndex[peerId] = rs.LastApplied
		} else {
			rs.nextIndex[peerId]--
		}
	}
	return nil
}

// This function is called by the leader to send heartbeats and log entries to all followers
func (rs *RaftState) AppendEntries(ctx context.Context, in *pb.AppendEntriesRequest) (reply *pb.AppendEntriesResponse, err error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if in.LeaderId != rs.LeaderId && rs.LeaderId != "" {
		log.Fatalf("PeerId should be same as leader ID. Expected leaderId %s %s", in.LeaderId, rs.LeaderId)
	}

	// TODO : Change this after leader election
	rs.LeaderId = in.LeaderId
	rs.LastHeartbeat = time.Now()

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
	for _, entry := range in.Entries {
		if entry.Index < uint64(len(rs.LogService.Logs)) && rs.LogService.Logs[entry.Index].Term != entry.Term {
			rs.LogService.Logs = rs.LogService.Logs[:entry.Index]
			reply.Term = rs.currentTerm
			reply.Success = false
			return reply, nil
		}
	}

	// 4.  Append any new entries not already in the log.
	// TODO: Can optimise this in one check.
	for _, entry := range in.Entries {
		if entry.Index >= uint64(len(rs.LogService.Logs)) {
			rs.LogService.PersistLogEntry(logging.LogEntry{Term: entry.Term, Index: entry.Index, Command: entry.Command})
		}
	}

	// 5.  If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
	if in.LeaderCommit > rs.CommitIndex {
		rs.CommitIndex = utils.Min(in.LeaderCommit, uint64(len(rs.LogService.Logs)-1))
	}

	reply.Term = rs.currentTerm
	reply.Success = true

	// log.Printf("Id: %v , Logs after AppendEntries RPC %v", rs.Id, rs.LogService.Logs)
	return reply, nil
}

func (rs *RaftState) CheckLastHeartbeat() {
	ticker := time.NewTicker(time.Duration(rs.ElectionTimeoutMs) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		// Check if the last heartbeat received is greater than the election timeout
		if rs.DiscoveryService.Status && rs.Mode != "Leader" && time.Since(rs.LastHeartbeat) > time.Duration(rs.ElectionTimeoutMs)*time.Millisecond {
			rs.StartElection()
			log.Println("Election timeout")
		}
	}
}

func (rs *RaftState) StartElection() {
	rs.mu.Lock()
	rs.currentTerm++
	rs.votedFor = rs.Id
	rs.Mode = "Candidate"
	rs.mu.Unlock()

	votes := 1
	Peers := rs.DiscoveryService.Peers
	voteCh := make(chan bool, len(Peers))

	for peerId, peer := range Peers {
		if peerId == rs.Id {
			continue
		}

		go rs.sendRequestVote(peer, voteCh)
	}

	for range Peers {
		if <-voteCh {
			votes++
		}
		if votes > len(Peers)/2 {
			rs.mu.Lock()
			rs.Mode = "Leader"
			rs.mu.Unlock()
			go rs.StartPeriodicAppends()
			return
		}
	}

	rs.mu.Lock()
	rs.Mode = "Follower"
	rs.mu.Unlock()
}

func (rs *RaftState) sendRequestVote(peer discovery.PeerData, voteCh chan<- bool) {
	conn, err := grpc.Dial(peer.URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		voteCh <- false
		return
	}
	defer conn.Close()

	client := pb.NewConsensusServiceClient(conn)

	rs.mu.Lock()
	lastLogIndex := rs.LastApplied
	lastLogTerm := rs.currentTerm - 1
	currentTerm := rs.currentTerm
	rs.mu.Unlock()

	req := &pb.RequestVoteRequest{
		Term:         currentTerm,
		CandidateId:  rs.Id,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  lastLogTerm,
	}

	resp, err := client.RequestVote(context.Background(), req)
	if err != nil || !resp.VoteGranted {
		voteCh <- false
		return
	}

	voteCh <- true
}

func (rs *RaftState) RequestVote(ctx context.Context, in *pb.RequestVoteRequest) (reply *pb.RequestVoteResponse, err error) {
	log.Println("Received RequestVote RPC")
	rs.mu.Lock()
	defer rs.mu.Unlock()

	reply = &pb.RequestVoteResponse{}

	// 1. Reply false if term < currentTerm (§5.1)
	if in.Term < rs.currentTerm {
		reply.Term = rs.currentTerm
		reply.VoteGranted = false
		return reply, nil
	}

	// 2. If votedFor is null or candidateId, and candidate’s log is at
	// least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	if (rs.votedFor == "" || rs.votedFor == in.CandidateId) && (in.LastLogTerm >= rs.currentTerm && in.LastLogIndex >= rs.LastApplied) {
		rs.votedFor = in.CandidateId
		reply.Term = rs.currentTerm
		reply.VoteGranted = true
		return reply, nil
	}

	reply.Term = rs.currentTerm
	reply.VoteGranted = false
	return reply, nil
}