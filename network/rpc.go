package network

import (
	"Raft/config"
	"Raft/consensus"
	pb_consenus "Raft/proto/consensus"
	pb_discovery "Raft/proto/discovery"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	pb_discovery.RegisterDiscoveryServiceServer(server, rs.DiscoveryService)
	pb_consenus.RegisterConsensusServiceServer(server, rs)

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

	// ShutdownHandling(rs)

	return &RaftServer{
		RaftState: rs,
	}, nil
}

func ShutdownHandling(rs *consensus.RaftState) {
	// Ensure ShutdownHandling is called on panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			err := rs.ShutdownHandling()
			if err != nil {
				log.Fatalf("Failed to shut down Raft state: %v", err)
			}
			os.Exit(1)
		}
	}()

	// Set up signal handling to gracefully shut down
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Received signal: %v. Shutting down...", sig)
		err := rs.ShutdownHandling()
		if err != nil {
			log.Fatalf("Failed to shut down Raft state: %v", err)
		}
		os.Exit(0)
	}()
}
