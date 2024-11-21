package main

import (
	"Raft/consensus"
	"Raft/network"
	"fmt"
	"os"
)

func isBootNodeArgsPassed(argsWithoutProg []string) bool {
	if len(argsWithoutProg) > 0 {
		if argsWithoutProg[0] == "--bootNode" {
			return true
		}
	}
	return false
}

func main() {
	isBootNode := isBootNodeArgsPassed(os.Args[1:])

	raftNode, err := consensus.InitialiseRaftNode()
	network.InitialiseRaftServer(raftNode, isBootNode)

	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Raft is here!!")

	// Block the main function indefinitely
	select {}
}
