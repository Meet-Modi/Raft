package consensus

import (
	"Raft/discovery"
	"Raft/logging"
	pb "Raft/proto/consensus"
	"Raft/utils"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RaftState struct {
	mu               sync.Mutex
	Id               string                      // Unique ID of the node
	currentTerm      uint64                      // Persistent state on all servers
	votedFor         string                      // Persistent state on all servers
	LeaderId         string                      // Persistent state on all servers
	CommitIndex      int64                       // Volatile state on all servers
	LastApplied      int64                       // Volatile state on all servers
	nextIndex        map[string]int64            // Volatile state on leaders, reinitialized after election
	matchIndex       map[string]int64            // Volatile state on leaders, reinirialized after election
	Mode             string                      // Leader, Follower, Candidate
	LogService       *logging.LogService         // Persistent state on all servers
	DiscoveryService *discovery.DiscoveryService // Reference to the discovery service
	ElectionTimeout  time.Duration               // Election timeout duration
	LastHeartbeat    int64                       // Last heartbeat received: epoch time in milliseconds
	pb.UnimplementedConsensusServiceServer
}

func InitialiseRaftState() (*RaftState, error) {
	id := utils.GenerateRaftPeerId()

	logService, _ := logging.NewLogService(id)
	discoveryService := discovery.NewDiscoveryService(id)
	if !discoveryService.Status {
		return nil, fmt.Errorf("failed to initialise discovery service, no point in continuing")
	}

	// forceTime := 5 * time.Second
	// if config.IsBootNode {
	// 	forceTime = 8 * time.Second
	// }

	rs := &RaftState{
		Id:               id,
		currentTerm:      0,
		votedFor:         "",
		LeaderId:         "",
		CommitIndex:      0,
		LastApplied:      0,
		nextIndex:        make(map[string]int64),
		matchIndex:       make(map[string]int64),
		Mode:             "Follower",
		LogService:       logService,
		DiscoveryService: discoveryService,
		ElectionTimeout:  time.Duration(3+rand.Intn(6)) * time.Second,
		// ElectionTimeout: forceTime,
		LastHeartbeat: time.Now().UnixMilli(),
	}

	log.Printf("Election timeout for node %s is %v", rs.Id, rs.ElectionTimeout)

	rs.LogService.PersistLogEntry(logging.LogEntry{Term: 0, Index: 0, Command: "Initialised"})
	go rs.CheckLastHeartbeat()
	go rs.SendHeartbeats()
	return rs, nil
}

func (rs *RaftState) SendHeartbeats() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if rs.Mode == "Leader" {
			rs.mu.Lock()
			rs.BroadCastHeartbeat()
			rs.mu.Unlock()
		}
	}
}

func (rs *RaftState) CheckLastHeartbeat() {
	ticker := time.NewTicker(rs.ElectionTimeout)
	defer ticker.Stop()

	for range ticker.C {
		if rs.Mode != "Leader" {
			if time.Since(time.UnixMilli(rs.LastHeartbeat)) > rs.ElectionTimeout {
				log.Printf("Election timeout occured at node: %s, starting election", rs.Id)
				rs.Mode = "Candidate"
				rs.LastHeartbeat = time.Now().UnixMilli()
				rs.StartElection()
			}
		}
	}
}

func (rs *RaftState) StartElection() {
	rs.mu.Lock()
	rs.currentTerm++
	rs.votedFor = rs.Id
	rs.LastHeartbeat = time.Now().UnixMilli()
	rs.mu.Unlock()
	voteCh := make(chan *pb.RequestVoteResponse, len(rs.DiscoveryService.Peers))
	votesInFavour := 1
	votesReceived := 1
	rs.BroadcastVoteRequest(voteCh)

	// No need to implement timout here, since for each API call, a timeout is already implemented. It timed out, it will return false
	for voteResponse := range voteCh {
		log.Printf("Node : %s received vote response from peer %v", rs.Id, voteResponse)
		votesReceived++
		if voteResponse.VoteGranted && rs.Mode == "Candidate" { // It is possible that a new leader was elected and this node is still receiving votes
			votesInFavour++
		} else {
			if voteResponse.Term > rs.currentTerm {
				log.Printf("Node : %v received higher term from peer %v, abandoning election", rs.Id, voteResponse.Term)
				// Abandon election and convert to follower
				close(voteCh)
				rs.mu.Lock()
				rs.currentTerm = voteResponse.Term
				rs.Mode = "Follower"
				rs.LastHeartbeat = time.Now().UnixMilli()
				rs.mu.Unlock()
				break
			}
		}

		log.Printf("Node : %s Votes in favour: %d, Votes left: %d", rs.Id, votesInFavour, len(rs.DiscoveryService.Peers)-votesReceived)

		if votesInFavour > len(rs.DiscoveryService.Peers)/2 {
			close(voteCh)
			rs.BecomeLeader()
			break
		}

		if votesReceived == len(rs.DiscoveryService.Peers) {
			close(voteCh)
			break
		}
	}

	// Since this election failed due to not sufficient votes, convert to follower and get new election timeout
	if rs.Mode == "Candidate" {
		rs.mu.Lock()
		rs.Mode = "Follower"
		rs.ElectionTimeout = time.Duration(3+rand.Intn(6)) * time.Second
		rs.mu.Unlock()
		log.Printf("Node : %s election failed, new election timeout is %v", rs.Id, rs.ElectionTimeout)
	}
}

func (rs *RaftState) BroadcastVoteRequest(voteCh chan<- *pb.RequestVoteResponse) {
	for peerId, peer := range rs.DiscoveryService.Peers {
		if peerId != rs.Id {
			go rs.SendRequestVote(peer, voteCh)
		}
	}
}

func (rs *RaftState) SendRequestVote(peer discovery.PeerData, voteCh chan<- *pb.RequestVoteResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			log.Println("No longer accepting vote result, seems election is already closed!", r)
		}
	}()

	conn, err := grpc.Dial(peer.URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		voteCh <- &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: false}
		log.Printf("Error in dialing the peer %s %v", peer.URI, err)
		return
	}
	defer conn.Close()
	client := pb.NewConsensusServiceClient(conn)

	request := &pb.RequestVoteRequest{
		Term:         rs.currentTerm,
		CandidateId:  rs.Id,
		LastLogIndex: rs.LastApplied,
		LastLogTerm:  rs.LogService.Logs[rs.LastApplied].Term,
	}

	response, err := client.RequestVote(ctx, request)
	if err != nil {
		log.Printf("Error in making Request Vote from peer: %s err : %v", rs.Id, err)
		voteCh <- &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: false}
		return
	}

	if response != nil {
		voteCh <- response
	} else {
		voteCh <- &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: false}
	}
}

func (rs *RaftState) RequestVote(ctx context.Context, request *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if request.Term < rs.currentTerm {
		return &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: false}, nil
	}

	if request.Term > rs.currentTerm {
		rs.currentTerm = request.Term
		rs.Mode = "Follower"
		rs.votedFor = "" // Reset votedFor as the term has changed
		rs.LastHeartbeat = time.Now().UnixMilli()
	}

	lastLogIndex := rs.LastApplied
	lastLogTerm := rs.LogService.Logs[lastLogIndex].Term
	// Grant vote if the candidate's log is at least as up-to-date as the voter's log
	if (rs.votedFor == "" || rs.votedFor == request.CandidateId) &&
		(request.LastLogTerm > lastLogTerm || (request.LastLogTerm == lastLogTerm && request.LastLogIndex >= lastLogIndex)) {
		rs.votedFor = request.CandidateId
		return &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: true}, nil
	}

	return &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: false}, nil
}

func (rs *RaftState) BecomeLeader() {
	log.Printf("Node : %s has been elected as leader", rs.Id)
	rs.mu.Lock()
	rs.LeaderId = rs.Id
	rs.Mode = "Leader"
	rs.LastHeartbeat = time.Now().UnixMilli()
	rs.nextIndex = make(map[string]int64)
	rs.matchIndex = make(map[string]int64)
	for peerId := range rs.DiscoveryService.Peers {
		rs.nextIndex[peerId] = rs.LastApplied + 1
		rs.matchIndex[peerId] = 0
	}
	rs.mu.Unlock()
	rs.BroadCastHeartbeat()
}

func (rs *RaftState) BroadCastHeartbeat() {
	for peerId, peer := range rs.DiscoveryService.Peers {
		if peerId != rs.Id {
			log.Printf("Node : %s sending heartbeat to peer %s", rs.Id, peerId)
			go rs.SendHeartbeat(peer)
		}
	}
}

func (rs *RaftState) SendHeartbeat(peer discovery.PeerData) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Some error in sending heartbeat!!", r)
		}
	}()

	conn, err := grpc.Dial(peer.URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Error in dialing the peer %s %v", peer.URI, err)
	}
	defer conn.Close()
	client := pb.NewConsensusServiceClient(conn)

	request := &pb.AppendEntriesRequest{
		Term:         rs.currentTerm,
		LeaderId:     rs.Id,
		PrevLogIndex: rs.LastApplied,
		PrevLogTerm:  rs.LogService.Logs[rs.LastApplied].Term,
		LeaderCommit: rs.CommitIndex,
		Entries:      []*pb.LogEntry{},
	}

	response, err := client.AppendEntries(ctx, request)
	if err != nil {
		log.Printf("Error in making AppendEntries from peer: %s err : %v", rs.Id, err)
	}

	if response.Success {
		log.Println("Heartbeat success")
	} else {
		log.Println("Heartbeat failed")
	}
}

func (rs *RaftState) AppendEntries(ctx context.Context, request *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	log.Printf("Node : %s term %v received AppendEntries from peer %s   - %v", rs.Id, rs.currentTerm, request.LeaderId, request)
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if request.Term < rs.currentTerm {
		log.Printf("Node : %s received lower term from peer %s, rejecting", rs.Id, request.LeaderId)
		return &pb.AppendEntriesResponse{Term: rs.currentTerm, Success: false}, nil
	}

	if request.PrevLogIndex > rs.LastApplied {
		return &pb.AppendEntriesResponse{Term: rs.currentTerm, Success: false}, nil
	}

	if rs.LogService.Logs[request.PrevLogIndex].Term != request.PrevLogTerm {
		return &pb.AppendEntriesResponse{Term: rs.currentTerm, Success: false}, nil
	}

	if request.Term > rs.currentTerm {
		rs.currentTerm = request.Term
		rs.Mode = "Follower"
		rs.votedFor = "" // Reset votedFor as the term has changed
		rs.LastHeartbeat = time.Now().UnixMilli()
		log.Printf("Node : %s received higher term from peer %s, converting to follower", rs.Id, request.LeaderId)
	}

	rs.LastHeartbeat = time.Now().UnixMilli()

	for _, entry := range request.Entries {
		rs.LastApplied++
		rs.LogService.PersistLogEntry(logging.LogEntry{Term: entry.Term, Index: rs.LastApplied, Command: entry.Command})
	}

	if request.LeaderCommit > rs.CommitIndex {
		rs.CommitIndex = int64(utils.Min(uint64(request.LeaderCommit), uint64(rs.LastApplied)))
	}

	return &pb.AppendEntriesResponse{Term: rs.currentTerm, Success: true}, nil
}
