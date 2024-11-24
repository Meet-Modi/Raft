package discovery

import (
	"fmt"

	"Raft/config"
	basediscovery "Raft/discovery/base_discovery"
	"Raft/discovery/gossip"
)

func NewDiscoveryService() (basediscovery.DiscoveryService, error) {
	switch config.DiscoveryProtocol {
	case "gossip":
		return gossip.NewDiscoveryService(config.IsBootNode), nil
	default:
		return nil, fmt.Errorf("unknown discovery type: %s", config.DiscoveryProtocol)
	}
}
