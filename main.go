package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	pb "Mandatory-Activity-5/grpc"

	"google.golang.org/grpc"
)

func parsePeers(peers string) []string {
	if peers == "" {
		return []string{}
	}
	parts := strings.Split(peers, ",")
	out := []string{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func main() {
	id := flag.Int("id", 1, "numeric ID of this node")
	addr := flag.String("addr", ":5001", "address this node listens on")
	peersFlag := flag.String("peers", "", "comma separated list of peers")
	isLeader := flag.Bool("leader", false, "whether this node starts as leader")
	flag.Parse()

	peers := parsePeers(*peersFlag)

	fmt.Printf(" Starting Auction Node %d", *id)
	fmt.Println("Address:", *addr)
	fmt.Println("Is leader:", *isLeader)
	fmt.Println("Peers:", peers)

	auctionSrv := NewAuctionServer(60 * time.Second)
	replicaSrv := newReplicaServer(auctionSrv)

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Node %d failed to listen on %s: %v", *id, *addr, err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuctionServer(grpcServer, auctionSrv)
	pb.RegisterReplicaServer(grpcServer, replicaSrv)

	if *isLeader {
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Println("Node is leader; connecting to replicas")
			auctionSrv.connectToReplicas(peers)
		}()
	} else {
		fmt.Println("Node starts as follower.")
	}

	fmt.Printf("Node %d listening on %s\n\n", *id, *addr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Node %d failed to serve: %v", *id, err)
	}
}
