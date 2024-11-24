package network

import (
	"Raft/config"
	"Raft/consensus"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type RaftServer struct {
	RaftState *consensus.RaftState
}

func InitialiseRaftServer() (*RaftServer, error) {
	rs, err := consensus.InitialiseRaftState()

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the discovery service with the gRPC server
	// pb_discovery.RegisterDiscoveryServiceServer(server, &rs.DiscoveryService)
	// pb_consenus.RegisterConsensusServiceServer(server, rs.UnimplementedConsensusServiceServer)

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

	return &RaftServer{
		RaftState: rs,
	}, nil
}
