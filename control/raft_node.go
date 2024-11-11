package control

import (
	"Raft/logging"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

type RaftNode struct {
	id          string
	currentTerm int
	votedFor    string
	isLeader    bool
	LogService  *logging.LogService
	commitIndex int
	lastApplied int
	nextIndex   map[string]int
	matchIndex  map[string]int

	networkNodes map[string]string //this stores all the Nodes IP
}

func InitialiseRaftNode() (*RaftNode, error) {
	log, _ := logging.NewLogService()

	return &RaftNode{
		id:          generateRaftNodeName(),
		currentTerm: 0,
		votedFor:    "",
		isLeader:    false, // First initialise with follower mode
		LogService:  log,
		commitIndex: 0,
		lastApplied: 0,
	}, nil

}

func generateRaftNodeName() string {
	// TODO: It should check with other node names
	time := uint64(time.Now().UnixMilli())
	rand.Seed(time)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 21)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	result = append([]byte(strconv.FormatUint(time, 10)), result...)
	fmt.Println("Raft Node name : ", string(result))
	return string(result)
}
