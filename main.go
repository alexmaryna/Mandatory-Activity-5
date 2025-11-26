package main

import (
	pb "Mandatory-Activity-5/grpc"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting All 3 Auction Nodes")

	go startNode(1, ":5001", true, []string{"localhost:5002", "localhost:5003"})
	time.Sleep(1 * time.Second)

	go startNode(2, ":5002", false, []string{"localhost:5001", "localhost:5003"})
	time.Sleep(1 * time.Second)

	go startNode(3, ":5003", false, []string{"localhost:5001", "localhost:5002"})
	time.Sleep(1 * time.Second)

	fmt.Println("All nodes are running")
	fmt.Println("\nNow you can run: go run client.go")
	fmt.Println("\nPress Ctrl+C to stop all nodes")
	select {}
}

func startNode(id int, address string, isPrimary bool, replicaAddresses []string) {
	auctionSrv := NewAuctionServer(60 * time.Second)
	replicaSrv := newReplicaServer(auctionSrv)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Node %d failed to listen: %v", id, err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuctionServer(grpcServer, auctionSrv)
	pb.RegisterReplicaServer(grpcServer, replicaSrv)

	fmt.Printf("Node %d started on %s (primary: %v)\n", id, address, isPrimary)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("Node %d failed to serve: %v", id, err)
		}
	}()
	time.Sleep(2 * time.Second)

	if isPrimary {
		auctionSrv.connectToReplicas(replicaAddresses)
	}
}
