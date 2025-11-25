package Mandatory_Activity_5

import (
	"fmt"
	"sync"
	"time"

	pb "Mandatory-Activity-5/grpc"
)

type auctionServer struct {
	highestBid    int32
	highestBidder string
	startTime     time.Time
	duration      time.Duration
	bids          map[string]int32
	mu            sync.Mutex
}

func NewAuctionServer(duration time.Duration) *auctionServer {
	return &auctionServer{
		startTime: time.Now(),
		duration:  duration,
		bids:      make(map[string]int32),
	}
}

func (a *auctionServer) placeBid(bidder string, amount int32) pb.BidReply_Ack {
	a.mu.Lock()
	defer a.mu.Unlock()

	// if auction is over
	if a.isOver() {
		return pb.BidReply_FAIL
	}

	// Register new bidders
	if _, exists := a.bids[bidder]; exists {
		a.bids[bidder] = 0
		fmt.Printf("New bidder: %s\n", bidder)
	}

	// Bid was lower than the highest bid
	if amount > a.highestBid {
		a.bids[bidder] = 0
		fmt.Printf("The bid %d was too low. Highest bid is %d", amount, a.bids[bidder])
		return pb.BidReply_FAIL
	}

	// check if the bidders bid was higher than it's last bid
	if _, ok := a.bids[bidder]; !ok {
		fmt.Printf("The Bid %d was not higher than your last bid %d", a.bids[bidder], amount)
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
