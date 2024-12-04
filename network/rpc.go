package network

import (
	"Raft/config"
	"Raft/kvstore"
	pb_consenus "Raft/proto/consensus"
	pb_discovery "Raft/proto/discovery"
	pb_kvstore "Raft/proto/kvstore"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type RaftServer struct {
	KeyValue *kvstore.KVStore
}

func InitialiseRaftServer() (*RaftServer, error) {
	kv, err := kvstore.InitialiseKVStore()
	if err != nil {
		panic("Failed to initialise Raft state: " + err.Error())
	}

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the discovery service with the gRPC server
	pb_discovery.RegisterDiscoveryServiceServer(server, kv.RaftState.DiscoveryService)
	pb_consenus.RegisterConsensusServiceServer(server, kv.RaftState)
	pb_kvstore.RegisterKVStoreServer(server, kv)

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

	return &RaftServer{}, nil
}
