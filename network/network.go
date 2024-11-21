package network

import (
	"Raft/config"
	"log"
	"net/rpc"
	"strconv"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

type peerData struct {
	endpoint string
	port     int
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

func (network *RaftNetwork) AllocateAvailableEndpoint() error {
	network.mu.Lock()
	defer network.mu.Unlock()

	network.RefreshPeerData()

	if network.isBootNode {
		network.Peers[network.PeerId] = peerData{
			port:     config.BootNodePort,
			endpoint: config.BootNodeURI,
		}
		return nil
	}

	// Create a set of used ports for faster lookup
	usedPorts := make(map[int]struct{})
	for _, peer := range network.Peers {
		usedPorts[peer.port] = struct{}{}
	}

	// Seed the random number generator once
	rand.Seed(uint64(time.Now().UnixMilli()))

	for {
		port := rand.Intn(config.MaxPort-config.MinPort+1) + config.MinPort
		if _, found := usedPorts[port]; !found {
			// Port is not in use, assign it
			network.Peers[network.PeerId] = peerData{
				port:     port,
				endpoint: config.NetworkURI,
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
			err := MakeRPCRequest(network.bootNodeEndpoint, network.bootNodePort, "RaftRPCServer.GetRaftNetwork", &ExampleArgs{}, network)
			if err != nil {
				log.Fatal("Error in getting network data from boot node")
			}
		}
	} else {
		// Get data from any known node
		for _, peer := range network.Peers {
			err := MakeRPCRequest(peer.endpoint, peer.port, "RaftRPCServer.GetRaftNetwork", &ExampleArgs{}, network)
			if err != nil {
				log.Fatal("Error in getting network data from known node")
			}
		}
	}
}

func MakeRPCRequest(endpoint string, port int, method string, args interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP("tcp", endpoint+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	err = client.Call(method, args, reply)
	return err
}
