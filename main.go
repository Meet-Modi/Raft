package main

import (
	"Raft/network"
	"fmt"
)

func main() {
	network.InitialiseRaftServer()
	fmt.Println("Raft is here!!")

	// Block the main function indefinitely
	select {}
}
