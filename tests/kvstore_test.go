package tests

import (
	"context"
	"testing"
	"time"

	pb_kvstore "Raft/proto/kvstore"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestPut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to RPC server: %v", err)
	}
	defer conn.Close()

	client := pb_kvstore.NewKVStoreClient(conn)

	putReq := &pb_kvstore.PutRequest{Key: "foo", Value: "bar"}
	_, err = client.Put(ctx, putReq)
	assert.NoError(t, err, "Put operation failed")
}

func TestGet(t *testing.T) {
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to RPC server: %v", err)
	}
	defer conn.Close()

	client := pb_kvstore.NewKVStoreClient(conn)

	getReq := &pb_kvstore.GetRequest{Key: "foo"}
	getResp, err := client.Get(context.Background(), getReq)
	assert.NoError(t, err, "Get operation failed")
	assert.True(t, getResp.Found, "Key not found")
	assert.Equal(t, "bar", getResp.Value, "Incorrect value")
}

func TestDelete(t *testing.T) {
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to RPC server: %v", err)
	}
	defer conn.Close()

	client := pb_kvstore.NewKVStoreClient(conn)

	deleteReq := &pb_kvstore.DeleteRequest{Key: "foo"}
	_, err = client.Delete(context.Background(), deleteReq)
	assert.NoError(t, err, "Delete operation failed")

	getReq := &pb_kvstore.GetRequest{Key: "foo"}
	getResp, err := client.Get(context.Background(), getReq)
	assert.NoError(t, err, "Get operation failed")
	assert.False(t, getResp.Found, "Key should not be found after delete")
}
