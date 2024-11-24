package basediscovery

import (
	pb "Raft/proto/discovery"
	"context"
)

// BaseDiscoveryService contains common fields and methods
type BaseDiscoveryService struct {
	pb.UnimplementedDiscoveryServiceServer
}

// DiscoveryService interface defines the methods that need to be implemented
type DiscoveryService interface {
	StartDiscovery()
	ServeDiscoverPeers(ctx context.Context, in *pb.DiscoveryRequest) (*pb.DiscoveryDataResponse, error)
}
