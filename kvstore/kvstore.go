package kvstore

import (
	"context"
	"fmt"
	"log"
	"sync"

	"Raft/consensus"
	pb "Raft/proto/kvstore"
)

type KVStore struct {
	mu        sync.Mutex
	store     map[string]string
	RaftState *consensus.RaftState
	pb.UnimplementedKVStoreServer
}

func InitialiseKVStore() *KVStore {
	rs, err := consensus.InitialiseRaftState()
	if err != nil {
		panic("Failed to initialise Raft state: " + err.Error())
	}

	return &KVStore{
		RaftState: rs,
		store:     make(map[string]string),
	}
}

func (kv *KVStore) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	log.Printf("PUT %s %s;", req.Key, req.Value)
	kv.store[req.Key] = req.Value
	kv.RaftState.ApplyComandToStateMachine(fmt.Sprintf("PUT %s %s;", req.Key, req.Value))
	log.Printf("PUT %s %s;", req.Key, req.Value)
	return &pb.PutResponse{Status: true}, nil
}

func (kv *KVStore) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	value, found := kv.store[req.Key]
	kv.RaftState.ApplyComandToStateMachine(fmt.Sprintf("GET %s;", req.Key))
	return &pb.GetResponse{Value: value, Found: found}, nil
}

func (kv *KVStore) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.store, req.Key)
	kv.RaftState.ApplyComandToStateMachine(fmt.Sprintf("DELETE %s;", req.Key))
	return &pb.DeleteResponse{Status: true}, nil
}
