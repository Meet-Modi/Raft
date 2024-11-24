package discovery

import (
	"context"
	"fmt"

	"Raft/config"
	"Raft/discovery/gossip"
	pb "Raft/proto/discovery/gossip"
)

type DiscoveryService interface {
	MakeDiscoveryRequest()
	DiscoverPeers(ctx context.Context, in *pb.DiscoveryRequest) (*pb.DiscoveryDataResponse, error)
}

func NewDiscoveryService() (DiscoveryService, error) {
	switch config.DiscoveryProtocol {
	case "gossip":
		return gossip.NewDiscoveryService(config.IsBootNode), nil
	default:
		return nil, fmt.Errorf("unknown discovery type: %s", config.DiscoveryProtocol)
	}
}
