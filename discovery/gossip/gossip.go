package gossip

import (
	"Raft/config"
	basediscovery "Raft/discovery/base_discovery"
	pb "Raft/proto/discovery"
	"Raft/utils"
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
	URI string
}

type GossipDiscoveryService struct {
	basediscovery.BaseDiscoveryService
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
		BootNodeId: utils.GenerateRaftPeerId(true),
	}

	// First add boot node to the peer list
	discoveryService.Peers[discoveryService.BootNodeId] = PeerData{
		URI: config.BootNodeURI,
	}

	go discoveryService.StartPeriodicDiscovery()
	return discoveryService
}

func (ds *GossipDiscoveryService) StartPeriodicDiscovery() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			log.Println("=========Starting periodic discovery=========")
			ds.StartDiscovery()
			ds.PrintAllPeers()
			log.Println("=========Finished periodic discovery=========")
		}
	}
}

func (ds *GossipDiscoveryService) StartDiscovery() {
	// This function calls DiscoverPeers on any known node
	ds.mu.Lock()
	defer ds.mu.Unlock()

	knownPeers := make([]string, 0, len(ds.Peers))
	for k := range ds.Peers {
		knownPeers = append(knownPeers, k)
	}

	if len(knownPeers) == 1 && knownPeers[0] == ds.PeerId {
		log.Println("%s %s No known peers to make a discovery request", knownPeers, ds.PeerId)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// conn, err := grpc.DialContext(ctx, ds.Peers[peerId].URI, grpc.WithInsecure(), grpc.WithBlock())
	// conn, err := grpc.Dial(ds.Peers[peerId].URI, grpc.WithInsecure())
	conn, err := grpc.Dial(ds.Peers[peerId].URI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error in dialing the peer %s %v", ds.Peers[peerId].URI, err)
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

	response, err := client.ServeDiscoverPeers(ctx, request)
	if err != nil {
		log.Fatalf("Error in making discovery request %v", err)
	}

	for peerId, peerData := range response.Peers {
		ds.Peers[peerId] = PeerData{URI: peerData.Uri}
	}
}

func (ds *GossipDiscoveryService) ServeDiscoverPeers(ctx context.Context, in *pb.DiscoveryRequest) (*pb.DiscoveryDataResponse, error) {
	// Extract the client's IP address from the gRPC context
	p, ok := peer.FromContext(ctx)
	var clientIP string
	if ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			clientIP = addr.IP.String()
		}
	}

	// Lock the mutex to ensure thread safety
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.Peers[in.PeerId] = PeerData{URI: clientIP + ":" + config.Port}

	for k, v := range in.MyPeers {
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

func (ds *GossipDiscoveryService) PrintAllPeers() {
	log.Println("+++++++++Peer list++++++++++")
	for k, v := range ds.Peers {
		log.Printf("PeerId: %s, URI: %s", k, v.URI)
	}
	log.Println("+++++++++Peer list++++++++++")
}
