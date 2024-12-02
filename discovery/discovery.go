// If a node requested x and is not able to reach, it will not gossip to other nodes about it.
package discovery

import (
	"Raft/config"
	pb "Raft/proto/discovery"
	"context"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

type PeerData struct {
	URI          string
	RetriesCount int
}

type DiscoveryService struct {
	mu                sync.Mutex
	PeerId            string
	Peers             map[string]PeerData
	BootNode          PeerData
	Status            bool // This will be used to check if the service data is up and running or not.
	maxRetriesAllowed int
	pb.UnimplementedDiscoveryServiceServer
}

func NewDiscoveryService(id string) *DiscoveryService {
	discoveryService := &DiscoveryService{
		PeerId:            id,
		Peers:             make(map[string]PeerData),
		maxRetriesAllowed: config.MaxDiscoveryRetries,
		BootNode:          PeerData{URI: config.BootNodeURI, RetriesCount: 0},
		Status:            false,
	}
	discoveryService.StartDiscovery()
	go discoveryService.StartPeriodicDiscovery()
	return discoveryService
}

func (ds *DiscoveryService) sendDiscoverPeerRequest(peer PeerData, peerId string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.Dial(peer.URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Error in dialing the peer %v %v", peer, err)
		if peerId != "" {
			ds.handleDiscoveryRequestError(peerId)
		}
		return
	}
	defer conn.Close()

	client := pb.NewDiscoveryServiceClient(conn)

	// Create a DiscoveryRequest message and populate MyPeers
	myPeers := make(map[string]*pb.PeerData)
	for peerId, peerData := range ds.Peers {
		if peerData.RetriesCount == 0 {
			myPeers[peerId] = &pb.PeerData{Uri: peerData.URI}
		}
	}

	request := &pb.DiscoveryRequest{
		PeerId:  ds.PeerId,
		MyPeers: myPeers,
	}

	response, err := client.DiscoverPeers(ctx, request)
	if err != nil {
		log.Printf("Error in making discovery request %v", err)
		if peerId != "" {
			ds.handleDiscoveryRequestError(peerId)
		}
		return
	}

	for peerId, peerData := range response.Peers {
		// Possible scenarios:
		// 1. The data in the response is stale and the unreachable peer is still in the list
		// 2. The data in the response is stale and the unreachable peer is not in the list
		// 3. The data in the response is fresh and the unreachable peer was contacted successfully
		// 4. The data in the response is fresh and the unreachable peer location is changed
		// 5. Peer is not present in my list and hence, I need to add it to my list
		// Here i need to only check for 4 and 5, since 1,2,3 will be handled when it will periodically check for the peers.

		if _, ok := ds.Peers[peerId]; !ok {
			ds.Peers[peerId] = PeerData{URI: peerData.Uri, RetriesCount: 0}
		}
		if ds.Peers[peerId].RetriesCount != 0 && peerData.Uri != ds.Peers[peerId].URI {
			ds.Peers[peerId] = PeerData{URI: peerData.Uri, RetriesCount: 0}
		}
	}
}

func (ds *DiscoveryService) handleDiscoveryRequestError(peerId string) {
	tempPeerData := ds.Peers[peerId]
	tempPeerData.RetriesCount++
	ds.Peers[peerId] = tempPeerData
	log.Printf("Retries count for peer %v is %v", peerId, tempPeerData.RetriesCount)
	if tempPeerData.RetriesCount > ds.maxRetriesAllowed {
		delete(ds.Peers, peerId)
		log.Printf("Peer %v removed from the list", peerId)
	}
}

func (ds *DiscoveryService) DiscoverPeers(ctx context.Context, in *pb.DiscoveryRequest) (*pb.DiscoveryDataResponse, error) {
	// Lock the mutex to ensure thread safety
	ds.mu.Lock()
	defer ds.mu.Unlock()

	p, ok := peer.FromContext(ctx)
	var clientIP string
	if ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			clientIP = addr.IP.String() // Extract the client's IP address from the gRPC context
		}
	}

	ds.Peers[in.PeerId] = PeerData{URI: clientIP + ":" + config.Port}

	for k, v := range in.MyPeers {
		ds.Peers[k] = PeerData{URI: v.Uri}
	}

	// Prepare the peer data to be returned
	peerData := make(map[string]*pb.PeerData)
	for k, v := range ds.Peers {
		if v.RetriesCount == 0 {
			peerData[k] = &pb.PeerData{Uri: v.URI, RetriesCount: 0}
		}
	}

	return &pb.DiscoveryDataResponse{
		PeerId: ds.PeerId,
		Peers:  peerData,
	}, nil
}

func (ds *DiscoveryService) StartPeriodicDiscovery() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ds.StartDiscovery()
	}
}

func (ds *DiscoveryService) StartDiscovery() {
	requestedPeersCount := 0
	if len(ds.Peers) == 0 && !config.IsBootNode {
		ds.sendDiscoverPeerRequest(ds.BootNode, "")
		requestedPeersCount++
	} else {
		// For Boot node, at initialisation, it will not have any peers and it will just simply continue to do so until other nodes discover it.
		for peerId, peerData := range ds.Peers {
			if peerId != ds.PeerId {
				ds.sendDiscoverPeerRequest(peerData, peerId)
				requestedPeersCount++
			}
		}
	}
	log.Printf("Requested discovery from %d peers, total Nodes in network : %d", requestedPeersCount, len(ds.Peers))
	ds.Status = true
}
