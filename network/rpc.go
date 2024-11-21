package network

import (
	"Raft/config"
	"Raft/consensus"
	"Raft/logging"
	pb "Raft/proto"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type AppendEntriesRequest struct {
	Term         int
	LeaderId     int
	prevLogIndex int
	prevLogTerm  int
	entries      []logging.LogEntry
	leaderCommit int
}

type AppendEntriesResponse struct {
	term    int
	success bool
}

type RaftRPCService struct {
	raftNode    *consensus.RaftNode
	raftNetwork *RaftNetwork
	pb.UnimplementedRaftServiceServer
}

func InitialiseRaftServer(rn *consensus.RaftNode, isBootNode bool) (*RaftRPCService, error) {

	RaftRPCService := &RaftRPCService{
		raftNode: rn,
		raftNetwork: &RaftNetwork{
			PeerId:           GenerateRaftPeerId(),
			Peers:            make(map[string]peerData),
			isBootNode:       isBootNode,
			bootNodeEndpoint: config.BootNodeURI,
			bootNodePort:     config.BootNodePort,
		},
	}

	err := RaftRPCService.raftNetwork.AllocateAvailableEndpoint()
	if err != nil {
		fmt.Println("Error in allocating endpoint")
		return nil, err
	}

	endpoint := RaftRPCService.raftNetwork.Peers[RaftRPCService.raftNetwork.PeerId].endpoint
	port := strconv.Itoa(RaftRPCService.raftNetwork.Peers[RaftRPCService.raftNetwork.PeerId].port)

	lis, err := net.Listen("tcp", endpoint+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRaftServiceServer(s, RaftRPCService)
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	fmt.Println("Raft Node ", RaftRPCService.raftNetwork.PeerId, " is listening on port : ", RaftRPCService.raftNetwork.Peers[RaftRPCService.raftNetwork.PeerId].port)

	return RaftRPCService, nil
}

type ExampleArgs struct{}

func (r *RaftRPCService) GetCurrentTime(args *ExampleArgs, reply *uint64) error {
	*reply = uint64(time.Now().UnixNano())
	return nil
}

func (s *RaftRPCService) GetRaftNetworkData(ctx context.Context, req *pb.RaftNetworkDataRequest) (*pb.RaftNetworkDataResponse, error) {
	log.Printf("Received request from: %s", req.FromPeerId)

	s.raftNetwork.mu.Lock()
	defer s.raftNetwork.mu.Unlock()

	// Add or update the requesting node in the table
	if len(req.Peers) == 0 {
		log.Printf("Requesting node %s has empty peer table", req.FromPeerId)
	} else {

		if _, exists := s.raftNetwork.Peers[req.FromPeerId]; !exists {
			s.raftNetwork.Peers[req.FromPeerId] = peerData{
				endpoint: req.Peers[req.FromPeerId].Endpoint,
				port:     int(req.Peers[req.FromPeerId].Port),
			}
		} else {
			s.raftNetwork.Peers[req.FromPeerId] = peerData{
				endpoint: req.Peers[req.FromPeerId].Endpoint,
				port:     int(req.Peers[req.FromPeerId].Port),
			}
		}

		// Update any entries in the table
		for peerId, peerD := range req.Peers {
			s.raftNetwork.Peers[peerId] = peerData{
				endpoint: peerD.Endpoint,
				port:     int(peerD.Port),
			}
		}
	}

	// Convert s.raftNetwork.Peers to the expected type
	peers := make(map[string]*pb.PeerData)
	for key, value := range s.raftNetwork.Peers {
		peers[key] = &pb.PeerData{
			Endpoint: value.endpoint,
			Port:     int32(value.port),
		}
	}

	return &pb.RaftNetworkDataResponse{
		PeerId:           s.raftNetwork.PeerId,
		Peers:            peers,
		IsBootNode:       s.raftNetwork.isBootNode,
		BootNodeEndpoint: s.raftNetwork.bootNodeEndpoint,
		BootNodePort:     int32(s.raftNetwork.bootNodePort),
	}, nil
}
