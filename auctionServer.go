package main

import (
	pb "Mandatory-Activity-5/grpc"
	"context"
	"fmt"
	"sync"
	"time"
)

type auctionServer struct {
	pb.UnimplementedAuctionServer

	mu            sync.Mutex
	highestBid    int32
	highestBidder string
	startTime     time.Time
	duration      time.Duration
	bids          map[string]int32
}

func NewAuctionServer(duration time.Duration) *auctionServer {
	return &auctionServer{
		startTime: time.Now(),
		duration:  duration,
		bids:      make(map[string]int32),
	}
}

func (a *auctionServer) Bid(ctx context.Context, req *pb.BidRequest) (*pb.BidReply, error) {
	bidder := "client"

	ack := a.placeBid(bidder, req.Amount)
	return &pb.BidReply{Ack: ack}, nil
}

func (a *auctionServer) Result(ctx context.Context, req *pb.ResultRequest) (*pb.ResultReply, error) {
	highest, over, winner := a.getResult()
	return &pb.ResultReply{
		HighestBid:  highest,
		AuctionOver: over,
		Winner:      winner,
	}, nil
}

func (a *auctionServer) placeBid(bidder string, amount int32) pb.BidReply_Ack {
	a.mu.Lock()
	defer a.mu.Unlock()

	// if auction is over
	if a.isOver() {
		return pb.BidReply_FAIL
	}

	//previous bid for the bidder and 0 if new
	lastBid, existed := a.bids[bidder]
	if !existed {
		fmt.Printf("New bidder: %s\n", bidder)
	}

	//check if the bid is higher than their own last bid
	if amount <= lastBid {
		fmt.Printf("The bid %d was not higehr than your last bid %d\n", amount, lastBid)
		return pb.BidReply_FAIL
	}

	//check if the bid is higher than teh overall last bid
	if amount <= a.highestBid {
		fmt.Printf("The bid %d was too low. Highest bid is %d", amount, a.highestBid)
		return pb.BidReply_FAIL
	}

	a.highestBid = amount
	a.highestBidder = bidder
	a.bids[bidder] = amount

	fmt.Printf("Accepted bid: %s with %d", bidder, amount)
	return pb.BidReply_SUCCESS
}

func (a *auctionServer) getResult() (int32, bool, string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	isOver := a.isOver()
	return a.highestBid, isOver, a.highestBidder
}

func (a *auctionServer) isOver() bool {
	return time.Since(a.startTime) > a.duration
}

// update state from replication
func (a *auctionServer) setState(highestBid int32, highestBidder string, startTime time.Time, bids map[string]int32) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.highestBid = highestBid
	a.highestBidder = highestBidder
	a.startTime = startTime
	a.bids = bids
}

// Returns rigistred bidders
func (a *auctionServer) getBidders() []string {
	a.mu.Lock()
	defer a.mu.Unlock()

	bidders := make([]string, len(a.bids))
	for bidder := range a.bids {
		bidders = append(bidders, bidder)
	}
	return bidders
}

func (a *auctionServer) getBidCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.bids)
}
