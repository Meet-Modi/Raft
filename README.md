# Raft-based Key Value Store

A fault-tolerant key-value store based on the Raft consensus algorithm. This project implements a distributed key-value store that ensures consistency and fault tolerance using the Raft consensus protocol.

## Features

- **Fault Tolerance**: The key-value store can tolerate failures of some nodes while still providing consistent data.
- **Eventual Consistency**: Ensures that all nodes in the cluster have the same data in the state machine. It will take 1 heartbeat check to apply commands to state machine. 
- **Leader Election**: Automatically elects a leader to handle client requests.
- **Log Replication**: Replicates log entries across all nodes to ensure data consistency.

## Architecture

The system is composed of multiple nodes, each running an instance of the Raft consensus algorithm. One of the nodes is elected as the leader, which handles all client requests and replicates log entries to the follower nodes.

## Components

- **Raft Consensus Module**: Implements the Raft consensus algorithm, including leader election, log replication, and state machine management.
- **Key-Value Store**: Provides a simple key-value store interface with `Put`, `Get`, and `Delete` operations.
- **gRPC Interface**: Exposes the key-value store operations and Raft consensus operations via gRPC and protocol buffer.

## Getting Started

### Prerequisites

- Go 1.23 or later
- Docker (for running the nodes in containers)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/meet-modi/raft.git
    cd raft
    ```

2. Install the dependencies:

    ```sh
    go mod tidy
    ```
### Running the Nodes Natively

Following environment variables will have to set first, usage example can be found below:
- BOOT_NODE=true
- DISCOVERY_TYPE=gossip
- PORT=3000
- BOOT_NODE_URI=192.168.100.100:3000

### Running the Nodes on Docker

You can run the nodes using Docker Compose. The `docker-compose.yml` file is provided to set up a cluster of Raft nodes.

1. Build and start the nodes:

    ```sh
    docker compose build && docker compose up --scale raft-node=2
    ```

2. The nodes will start and elect a leader. You can check the logs to see the leader election process and log replication.

### Usage

You can interact with the key-value store using gRPC clients. The gRPC interface is defined in the `proto/kvstore/kvstore.proto` file.

#### Example gRPC Client

Here is an example of how to use the gRPC client to interact with the key-value store:

```go
package main

import (
    "context"
    "log"
    "time"

    pb "github.com/yourusername/raft-kvstore/proto/kvstore"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()

    client := pb.NewKVStoreClient(conn)

    // Put operation
    _, err = client.Put(context.Background(), &pb.PutRequest{Key: "key1", Value: "value1"})
    if err != nil {
        log.Fatalf("Failed to put key-value: %v", err)
    }

    // Get operation
    resp, err := client.Get(context.Background(), &pb.GetRequest{Key: "key1"})
    if err != nil {
        log.Fatalf("Failed to get key-value: %v", err)
    }
    log.Printf("Get response: %v", resp)

    // Delete operation
    _, err = client.Delete(context.Background(), &pb.DeleteRequest{Key: "key1"})
    if err != nil {
        log.Fatalf("Failed to delete key-value: %v", err)
    }
}