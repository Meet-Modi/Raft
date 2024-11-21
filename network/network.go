package network

import (
	"Raft/config"
	pb "Raft/proto"
	"Raft/utils"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
)

type peerData struct {
	Endpoint string
	Port     int
}

type RaftNetwork struct {
	mu                sync.Mutex
	PeerId            string
	Peers             map[string]peerData
	CurrentLeaderId   string
	CurrentLeaderPort int
	isBootNode        bool
	bootNodeEndpoint  string
	bootNodePort      int
}

func GenerateRaftPeerId() string {
	// TODO: It should check with other node names
	time := uint64(time.Now().UnixMilli())
	rand.Seed(time)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 10)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	result = append([]byte(strconv.FormatUint(time, 10)), result...)
	return string(result)
}

func (network *RaftNetwork) StartPeriodicRefresh() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				network.mu.Lock()
				network.RefreshPeerData(false)
				utils.PrintAllPeers(network)
				network.mu.Unlock()
			}
		}

	}()
}

func (network *RaftNetwork) AllocateAvailableEndpoint() error {
	network.mu.Lock()
	defer network.mu.Unlock()

	network.RefreshPeerData(true)

	if network.isBootNode {
		network.Peers[network.PeerId] = peerData{
			Port:     config.BootNodePort,
			Endpoint: config.BootNodeURI,
		}
		return nil
	}

	// Create a set of used ports for faster lookup
	usedPorts := make(map[int]struct{})
	for _, peer := range network.Peers {
		usedPorts[peer.Port] = struct{}{}
	}

	// Seed the random number generator once
	rand.Seed(uint64(time.Now().UnixMilli()))

	for {
		port := rand.Intn(config.MaxPort-config.MinPort+1) + config.MinPort
		if _, found := usedPorts[port]; !found {
			// Port is not in use, assign it
			network.Peers[network.PeerId] = peerData{
				Port:     port,
				Endpoint: config.NetworkURI,
			}
			break
		}
	}
	return nil
}

func (network *RaftNetwork) RefreshPeerData(initialisation bool) {
	// can be a boot node in initialisation phase -> network table will be empty
	// can be a non boot node in initialisation phase -> get data from boot node
	// can be a boot in normal phase -> get data from any known node
	// can be a non boot in normal phase -> get data from any known node

	if initialisation {
		if network.isBootNode {
			return
		} else {
			res, err := MakeGetRaftNetworkRequest(network.bootNodeEndpoint, network.bootNodePort, network)
			if err != nil {
				log.Fatal("Error in getting network data from boot node")
				return
			}
			UpdateRaftNetworkPeers(network, res)
		}
	} else {
		// Get data from any known node
		for key, peer := range network.Peers {
			if key == network.PeerId {
				continue
			}
			res, err := MakeGetRaftNetworkRequest(peer.Endpoint, peer.Port, network)
			if err != nil {
				log.Fatal("Error in getting network data from known node")
				return
			}
			UpdateRaftNetworkPeers(network, res)
		}
	}
}

func UpdateRaftNetworkPeers(network *RaftNetwork, response *pb.RaftNetworkDataResponse) {
	for peerId, peerD := range response.Peers {
		network.Peers[peerId] = peerData{
			Endpoint: peerD.Endpoint,
			Port:     int(peerD.Port),
		}
	}
}

func MakeGetRaftNetworkRequest(endpoint string, port int, network *RaftNetwork) (*pb.RaftNetworkDataResponse, error) {
	conn, err := grpc.Dial(endpoint+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewRaftServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Convert network.Peers to the expected type
	peers := make(map[string]*pb.PeerData)
	for key, value := range network.Peers {
		peers[key] = &pb.PeerData{
			Endpoint: value.Endpoint,
			Port:     int32(value.Port),
		}
	}

	response, err := client.GetRaftNetworkData(ctx, &pb.RaftNetworkDataRequest{FromPeerId: network.PeerId, Peers: peers})
	if err != nil {
		log.Fatalf("Error calling GetRaftNetwork: %v", err)
	}

	log.Printf("Response from server: %s", response.PeerId)
	return response, nil
}

func PrintAllPeers(network *RaftNetwork) {
	fmt.Printf("Peers in Id: %s\n", network.PeerId)
	fmt.Printf("%-20s %-20s %-10s\n", "PeerId", "Endpoint", "Port")
	fmt.Println(strings.Repeat("-", 50))
	for peerId, peer := range network.Peers {
		fmt.Printf("%-20s %-20s %-10d\n", peerId, peer.Endpoint, peer.Port)
	}
}
