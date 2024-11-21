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

func (network *RaftNetwork) RefreshPeerData() {
	// Pick any random peer and Check the peer nodes and update the data
	raftPeers := network.Peers
	keys := make([]string, 0, len(raftPeers))
	for k := range raftPeers {
		keys = append(keys, k)
	}
	for {
		// Pick any random peer
		peerId := keys[rand.Intn(len(keys))]
		peerEndpoint := network.Peers[peerId].endpoint
		peerPort := network.Peers[peerId].port

		// Get the data from the peer
		args := ExampleArgs{}
		var reply RaftNetwork
		client, err := rpc.DialHTTP("tcp", peerEndpoint+":"+strconv.Itoa(peerPort))
		if err != nil {
			log.Fatal("dialing:", err)
		}
		err = client.Call("RaftRPCServer.GetRaftNetwork", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		log.Printf("%d", reply)
	}
}
