package network

import (
	"Raft/consensus"
	"Raft/discovery"
	"log"
)

type RaftRPCService struct {
	raftNode         *consensus.RaftNode
	discoveryService *discovery.DiscoveryService
}

func InitialiseRaftServer(rn *consensus.RaftNode, isBootNode bool) (*RaftRPCService, error) {
	discoveryService, err := discovery.NewDiscoveryService()
	if err != nil {
		log.Fatalf("Failed to initialise discovery service: %v", err)
		return nil, err
	}

	return &RaftRPCService{
		raftNode:         rn,
		discoveryService: &discoveryService,
	}, nil
}
