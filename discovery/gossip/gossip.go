package gossip

import (
	"Raft/config"
	pb "Raft/proto/discovery/gossip"
	"Raft/utils"
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type PeerData struct {
	URI string
}

type GossipDiscoveryService struct {
	mu              sync.Mutex
	PeerId          string
	Peers           map[string]PeerData
	BootNodeId      string
	CurrentLeaderId string
	Servicestatus   bool // This will be used to check if the service data is correct or not.
	pb.UnimplementedDiscoveryServiceServer
}

func NewDiscoveryService(isBootNode bool) *GossipDiscoveryService {

	discoveryService := &GossipDiscoveryService{
		PeerId:     utils.GenerateRaftPeerId(isBootNode),
		Peers:      make(map[string]PeerData),
		BootNodeId: utils.GenerateRaftPeerId(isBootNode),
	}

	return discoveryService

}

func (ds *GossipDiscoveryService) RefreshPeerData(initialisation bool) {
	// can be a boot node in initialisation phase -> network table will be empty
	// can be a non boot node in initialisation phase -> get data from boot node
	// can be a boot in normal phase -> get data from any known node
	// can be a non boot in normal phase -> get data from any known node

	if initialisation {
		// First add boot node to the peer list
		ds.Peers[ds.BootNodeId] = PeerData{
			URI: config.BootNodeURI,
		}
		if ds.PeerId != ds.BootNodeId {
			// Get the peer data from the boot node
			// Add all the peer data to the peer list

		}
	}
}

func (ds *GossipDiscoveryService) MakeDiscoveryRequest() {
	// This function calls DiscoverPeers on any known node
	ds.mu.Lock()
	defer ds.mu.Unlock()

	knownPeers := make([]string, 0, len(ds.Peers))
	for k := range ds.Peers {
		knownPeers = append(knownPeers, k)
	}

	if len(knownPeers) == 1 && knownPeers[0] == ds.PeerId {
		log.Println("No known peers to make a discovery request")
		return
	}

	var peerId string
	for {
		peerId = utils.RandomKey(knownPeers)
		if peerId == ds.PeerId {
			continue
		}
		break
	}
	// Make a gRPC call to the peer
	// If the call is successful, update the peer list
	// TODO : If the call is unsuccessful, remove the peer from the list
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, ds.Peers[peerId].URI, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Error in dialing the peer", err)
		// TODO: delete(ds.Peers, peerId)
		// TODO: Can write another parellel process to delete the inactive peers
		return
	}
	defer conn.Close()

	client := pb.NewDiscoveryServiceClient(conn)

	// Create a DiscoveryRequest message and populate MyPeers
	myPeers := make(map[string]*pb.PeerData)
	for peerId, peerData := range ds.Peers {
		myPeers[peerId] = &pb.PeerData{Uri: peerData.URI}
	}

	request := &pb.DiscoveryRequest{
		PeerId:  ds.PeerId, // Use ds.PeerId
		MyPeers: myPeers,
	}

	response, err := client.DiscoverPeers(ctx, request)
	if err != nil {
		log.Fatalf("Error in making discovery request", err)
	}

	log.Printf("Received response from peer %s", response.PeerId)
	for peerId, peerData := range response.Peers {
		log.Printf("PeerId: %s, URI: %s", peerId, peerData.Uri)
	}

	// Update local peer data with the response
	for peerId, pData := range response.Peers {
		ds.Peers[peerId] = PeerData{URI: pData.Uri}
	}

}

func (ds *GossipDiscoveryService) DiscoverPeers(ctx context.Context, in *pb.DiscoveryRequest) (*pb.DiscoveryDataResponse, error) {
	// Lock the mutex to ensure thread safety
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Update the peer list with the received data
	for k, v := range in.MyPeers {
		// Skip adding/Updating self to the peer list.
		// TODO: Can remove it after testing
		if ds.PeerId == k {
			continue
		}
		ds.Peers[k] = PeerData{URI: v.Uri}
	}

	// Prepare the peer data to be returned
	peerData := make(map[string]*pb.PeerData)
	for k, v := range ds.Peers {
		peerData[k] = &pb.PeerData{Uri: v.URI}
	}

	return &pb.DiscoveryDataResponse{
		PeerId:          ds.PeerId,
		Peers:           peerData,
		CurrentLeaderId: ds.CurrentLeaderId, // Replace with actual leader ID
		BootNodeId:      ds.BootNodeId,      // Replace with actual boot node ID
	}, nil
}
