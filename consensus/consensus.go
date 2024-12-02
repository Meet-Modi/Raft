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
	CommitIndex       int64                       // Volatile state on all servers
	LastApplied       int64                       // Volatile state on all servers
	nextIndex         map[string]int64            // Volatile state on leaders, reinitialized after election
	matchIndex        map[string]int64            // Volatile state on leaders, reinirialized after election
	Mode              string                      // Leader, Follower, Candidate
	LogService        *logging.LogService         // Persistent state on all servers
	DiscoveryService  *discovery.DiscoveryService // Reference to the discovery service
	ElectionTimeoutMs int64                       // Election timeout in milliseconds
	LastHeartbeat     int64                       // Last heartbeat received: epoch time in milliseconds
	pb.UnimplementedConsensusServiceServer
}

func InitialiseRaftState() (*RaftState, error) {
	id := utils.GenerateRaftPeerId()

	logService, _ := logging.NewLogService(id)
	discoveryService := discovery.NewDiscoveryService(id)
	if !discoveryService.Status {
		return nil, fmt.Errorf("failed to initialise discovery service, no point in continuing")
	}

	rs := &RaftState{
		Id:                id,
		currentTerm:       0,
		votedFor:          "",
		LeaderId:          "",
		CommitIndex:       0,
		LastApplied:       0,
		nextIndex:         make(map[string]int64),
		matchIndex:        make(map[string]int64),
		Mode:              "Follower",
		LogService:        logService,
		DiscoveryService:  discoveryService,
		ElectionTimeoutMs: config.ElectionTimeoutMs,
	}
	return rs, nil
}

func (rs *RaftState) CheckLastHeartbeat() {
	if time.Now().UnixMilli()-rs.LastHeartbeat > rs.ElectionTimeoutMs {
		log.Printf("Election timeout occured at node: %s, starting election", rs.Id)
		rs.Mode = "Candidate"
		rs.LastHeartbeat = time.Now().UnixMilli()
		rs.StartElection()
	}
}

func (rs *RaftState) StartElection() {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.currentTerm++
	rs.votedFor = rs.Id
	rs.LastHeartbeat = time.Now().UnixMilli()
	voteCh := make(chan *pb.RequestVoteResponse, len(rs.DiscoveryService.Peers))
	votesInFavour := 1
	votesReceived := 1
	rs.BroadcastVoteRequest(voteCh)

	for voteResponse := range voteCh {
		votesReceived++
		if voteResponse.VoteGranted && rs.Mode == "Candidate" { // It is possible that a new leader was elected and this node is still receiving votes
			votesInFavour++
		} else {

			if voteResponse.Term > rs.currentTerm {
				// Abandon election and convert to follower
				rs.currentTerm = voteResponse.Term
				rs.Mode = "Follower"
				rs.LastHeartbeat = time.Now().UnixMilli()
				close(voteCh)
				return
			}
		}

		if votesInFavour > len(rs.DiscoveryService.Peers)/2 {
			rs.Mode = "Leader"
			close(voteCh)
			return
		}

		if votesReceived == len(rs.DiscoveryService.Peers) {
			close(voteCh)
			return
		}
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
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			log.Println("No longer accepting vote result, seems election is abandoned, coward!", r)
		}
	}()

	conn, err := grpc.Dial(peer.URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		voteCh <- &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: false}
		log.Printf("Error in dialing the peer %s %v", peer.URI, err)
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
		log.Printf("Error in making Request Vote %v", err)
		voteCh <- &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: false}
	}

	voteCh <- response
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
		rs.votedFor = ""
		rs.LastHeartbeat = time.Now().UnixMilli()
	}

	lastLogIndex := rs.LastApplied
	lastLogTerm := rs.LogService.Logs[lastLogIndex].Term

	if (rs.votedFor == "" || rs.votedFor == request.CandidateId) &&
		(request.LastLogTerm > lastLogTerm || (request.LastLogTerm == lastLogTerm && request.LastLogIndex >= lastLogIndex)) {
		rs.votedFor = request.CandidateId
		return &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: true}, nil
	}

	return &pb.RequestVoteResponse{Term: rs.currentTerm, VoteGranted: false}, nil
}
