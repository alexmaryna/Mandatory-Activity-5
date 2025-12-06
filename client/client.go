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
	fmt.Printf("Result: %v\n", reply1.Ack)

	time.Sleep(500 * time.Millisecond)

	fmt.Println("\nPlacing bid: 200")
	reply2, _ := client.Bid(context.Background(), &pb.BidRequest{Amount: 200})
	fmt.Printf("Result: %v\n", reply2.Ack)

	time.Sleep(500 * time.Millisecond)
	
	fmt.Println("\nPlacing bid: 150 (should fail)")
	reply3, _ := client.Bid(context.Background(), &pb.BidRequest{Amount: 150})
	fmt.Printf("Result: %v\n", reply3.Ack)

	time.Sleep(500 * time.Millisecond)

	fmt.Println("\nChecking result...")
	result, _ := client.Result(context.Background(), &pb.ResultRequest{})
	fmt.Printf("Highest bid: %d, Over: %v, Winner: %s\n", result.HighestBid, result.AuctionOver, result.Winner)
}
