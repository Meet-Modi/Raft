package kvstore

import (
	"context"
	"log"
	"strings"
	"sync"

	"Raft/consensus"
	pb "Raft/proto/kvstore"
)

type KVStore struct {
	mu              sync.Mutex
	store           map[string]string
	RaftState       *consensus.RaftState
	RaftProposeChan chan string // This channel will be used to propose a new command to the Raft state.
	RaftResultChan  chan string // This channel will be used to get the result of the proposed command.
	pb.UnimplementedKVStoreServer
}

func InitialiseKVStore() (*KVStore, error) {
	RaftProposeChan := make(chan string)
	RaftResultChan := make(chan string)

	rs, err := consensus.InitialiseRaftState(RaftProposeChan, RaftResultChan)
	if err != nil {
		panic(err)
	}
	kv := &KVStore{
		store:           make(map[string]string),
		RaftState:       rs,
		RaftProposeChan: RaftProposeChan,
		RaftResultChan:  RaftResultChan,
	}
	go kv.ReadResultsFromRaft()

	return kv, nil
}

func (kv *KVStore) ReadResultsFromRaft() {
	for command := range kv.RaftResultChan {
		kv.mu.Lock()
		kv.Apply(command)
		kv.mu.Unlock()
	}
}

func (kv *KVStore) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	command := "PUT " + req.Key + " " + req.Value
	log.Printf("Received PUT request on kvserver for key: %s, value: %s", req.Key, req.Value)
	kv.RaftProposeChan <- command
	return &pb.PutResponse{Status: true}, nil
}

func (kv *KVStore) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Received GET request on kvserver for key: %s", req.Key)
	value, found := kv.store[req.Key]
	return &pb.GetResponse{Value: value, Found: found}, nil
}

func (kv *KVStore) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	command := "DELETE " + req.Key
	log.Printf("Received DELETE request on kvserver for key: %s", req.Key)
	kv.RaftProposeChan <- command
	return &pb.DeleteResponse{Error: "Couldn't complete the transaction, rollbacked!"}, nil
}

func (kv *KVStore) Apply(command string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	parts := strings.SplitN(command, " ", 3)
	commandType := parts[0]
	if commandType == "PUT" && len(parts) == 3 {
		kv.applyPut(parts[1], parts[2])
	} else if commandType == "DELETE" && len(parts) == 2 {
		kv.applyDelete(parts[1])
	}
}

func (kv *KVStore) applyPut(key, value string) {
	kv.store[key] = value
}

func (kv *KVStore) applyDelete(key string) {
	delete(kv.store, key)
}
