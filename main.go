package main

import (
	"fmt"
	"time"

	pb "Mandatory-Activity-5/grpc"
)

func main() {
	auction := NewAuctionServer(30 * time.Second)

	fmt.Println("Starting simple auction test")

	res1 := auction.placeBid("client1", 50)
	fmt.Println("client1 bids 50:", res1 == pb.BidReply_SUCCESS)

	res2 := auction.placeBid("client2", 80)
	fmt.Println("client2 bids 80:", res2 == pb.BidReply_SUCCESS)

	res3 := auction.placeBid("client3", 60)
	fmt.Println("client3 bids 60:", res3 == pb.BidReply_FAIL)

	res4 := auction.placeBid("client4", 80)
	fmt.Println("client4 bids 80:", res4 == pb.BidReply_FAIL)

	highest, finished, winner := auction.getResult()

	fmt.Println("auction finished:", finished)
	fmt.Println("highest bid:", highest)
	fmt.Println("winner:", winner)
}
