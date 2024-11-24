package network

import (
	"Raft/config"
	"Raft/consensus"
	"Raft/discovery"
	base_discovery "Raft/discovery/base_discovery"
	pb "Raft/proto/discovery"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type RaftRPCService struct {
	raftNode         *consensus.RaftNode
	discoveryService base_discovery.DiscoveryService
}

func InitialiseRaftServer(rn *consensus.RaftNode) (*RaftRPCService, error) {

	discoveryService, err := discovery.NewDiscoveryService()
	if err != nil {
		log.Fatalf("Failed to initialise discovery service: %v", err)
		return nil, err
	}

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the discovery service with the gRPC server
	pb.RegisterDiscoveryServiceServer(server, discoveryService.(pb.DiscoveryServiceServer))

	// Listen on the specified port
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %v", config.Port, err)
	}

	log.Printf("gRPC server listening on port %s", config.Port)

	// Serve gRPC server
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return &RaftRPCService{
		raftNode:         rn,
		discoveryService: discoveryService,
	}, nil
}
