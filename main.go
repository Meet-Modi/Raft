package main

import (
	"Raft/consensus"
	"Raft/network"
	"fmt"
	"os"
)

func main() {

	raftNode, err := consensus.InitialiseRaftNode()
	network.InitialiseRaftServer(raftNode)

	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Raft is here!!")

	// Block the main function indefinitely
	select {}
}
