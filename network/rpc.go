package network

import (
	"Raft/config"
	"Raft/consensus"
	"Raft/logging"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

type AppendEntriesRequest struct {
	Term         int
	LeaderId     int
	prevLogIndex int
	prevLogTerm  int
	entries      []logging.LogEntry
	leaderCommit int
}

type AppendEntriesResponse struct {
	term    int
	success bool
}

type RaftRPCServer struct {
	raftNode    *consensus.RaftNode
	raftNetwork *RaftNetwork
}

func InitialiseRaftServer(rn *consensus.RaftNode, isBootNode bool) (*RaftRPCServer, error) {
	raftRPCServer := &RaftRPCServer{
		raftNode: rn,
		raftNetwork: &RaftNetwork{
			PeerId:           GenerateRaftPeerId(),
			Peers:            make(map[string]peerData),
			isBootNode:       isBootNode,
			bootNodeEndpoint: config.BootNodeURI,
			bootNodePort:     config.BootNodePort,
		},
	}

	// TODO : First get data from all other peers if not bootnode.
	if !isBootNode {
		// RefreshPeerData first.
	}

	rpc.Register(raftRPCServer)
	rpc.HandleHTTP()
	err := raftRPCServer.raftNetwork.AllocateAvailableEndpoint()
	if err != nil {
		fmt.Println("Error in allocating endpoint")
		return nil, err
	}
	port := strconv.Itoa(raftRPCServer.raftNetwork.Peers[raftRPCServer.raftNetwork.PeerId].port)

	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Println("Raft Node ", raftRPCServer.raftNetwork.PeerId, " is listening on port : ", raftRPCServer.raftNetwork.Peers[raftRPCServer.raftNetwork.PeerId].port)

	http.Serve(l, nil)
	return raftRPCServer, nil
}

type ExampleArgs struct{}

func (r *RaftRPCServer) GetCurrentTime(args *ExampleArgs, reply *uint64) error {
	*reply = uint64(time.Now().UnixNano())
	return nil
}

func (r *RaftRPCServer) GetRaftNetwork(args *ExampleArgs, reply *RaftNetwork) error {
	*reply = *r.raftNetwork
	return nil
}
