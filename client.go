package main

import (
	pb "Mandatory-Activity-5/grpc"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuctionClient(conn)

	fmt.Println("Testing Auction System")
	fmt.Println("\nPlacing bid: 100")
	reply1, _ := client.Bid(context.Background(), &pb.BidRequest{Amount: 100})
	fmt.Printf("   Result: %v\n", reply1.Ack)

	time.Sleep(1 * time.Second)

	fmt.Println("\nPlacing bid: 200")
	reply2, _ := client.Bid(context.Background(), &pb.BidRequest{Amount: 200})
	fmt.Printf("   Result: %v\n", reply2.Ack)

	time.Sleep(1 * time.Second)

	fmt.Println("\nPlacing bid: 150 (should fail)")
	reply3, _ := client.Bid(context.Background(), &pb.BidRequest{Amount: 150})
	fmt.Printf("   Result: %v\n", reply3.Ack)

	time.Sleep(1 * time.Second)

	fmt.Println("\nChecking result from primary...")
	result, _ := client.Result(context.Background(), &pb.ResultRequest{})
	fmt.Printf("   Highest bid: %d, Over: %v, Winner: %s\n", result.HighestBid, result.AuctionOver, result.Winner)

	fmt.Println("\nChecking result from backup node...")
	backupConn, err := grpc.Dial("localhost:5002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Could not connect to backup: %v", err)
		return
	}
	defer backupConn.Close()

	backupClient := pb.NewAuctionClient(backupConn)
	backupResult, _ := backupClient.Result(context.Background(), &pb.ResultRequest{})
	fmt.Printf("   Highest: %d, Over: %v, Winner: %s\n", backupResult.HighestBid, backupResult.AuctionOver, backupResult.Winner)

	if backupResult.HighestBid == result.HighestBid {
		fmt.Println("\n✓ Replication works! Primary and backup have same state")
	} else {
		fmt.Println("\n✗ Replication failed! States differ")
	}

}
