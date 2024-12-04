package kvstore

import (
	"context"
	"sync"

	pb "Raft/proto/kvstore"
)

type KVStore struct {
	mu    sync.Mutex
	store map[string]string
	pb.UnimplementedKVStoreServer
}

func InitialiseKVStore() *KVStore {
	return &KVStore{
		store: make(map[string]string),
	}
}

func (kv *KVStore) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	// if !kv.RaftState.StoreAndReplicateWAL(fmt.Sprintf("PUT %s %s;", req.Key, req.Value)) {
	// 	return &pb.PutResponse{Error: "Couldn't complete the transaction, rollbacked!"}, nil
	// }
	kv.store[req.Key] = req.Value
	return &pb.PutResponse{Status: true}, nil
}

func (kv *KVStore) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	// if !kv.RaftState.StoreAndReplicateWAL(fmt.Sprintf("GET %s;", req.Key)) {
	// 	return &pb.GetResponse{Error: "Couldn't complete the transaction, rollbacked!"}, nil
	// }
	value, found := kv.store[req.Key]
	return &pb.GetResponse{Value: value, Found: found}, nil
}

func (kv *KVStore) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	// kv.RaftState.StoreAndReplicateWAL(fmt.Sprintf("DELETE %s;", req.Key))
	delete(kv.store, req.Key)
	return &pb.DeleteResponse{Error: "Couldn't complete the transaction, rollbacked!"}, nil
}
