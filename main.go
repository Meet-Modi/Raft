package main

import (
	"Raft/consensus"
	"Raft/network"
	"fmt"
	"os"
)

func main() {

	RaftState, err := consensus.InitialiseRaftState()
	network.InitialiseRaftServer(RaftState)

	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Raft is here!!")

	// Block the main function indefinitely
	select {}
}
