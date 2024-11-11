package main

import (
	"Raft/control"
	"fmt"
	"os"
)

func main() {
	_, err := control.InitialiseRaftNode()

	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Raft is here!!")
}
