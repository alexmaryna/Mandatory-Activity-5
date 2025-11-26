package main

import (
	pb "Mandatory-Activity-5/grpc"
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type auctionServer struct {
	pb.UnimplementedAuctionServer

	mu            sync.Mutex
	highestBid    int32
	highestBidder string
	startTime     time.Time
	duration      time.Duration
	bids          map[string]int32

	replicaClients map[string]pb.ReplicaClient
	replicaMu      sync.RWMutex
}

func NewAuctionServer(duration time.Duration) *auctionServer {
	return &auctionServer{
		startTime:      time.Now(),
		duration:       duration,
		bids:           make(map[string]int32),
		replicaClients: make(map[string]pb.ReplicaClient),
	}
}

func (a *auctionServer) Bid(ctx context.Context, req *pb.BidRequest) (*pb.BidReply, error) {
	bidder := "client"

	ack := a.placeBid(bidder, req.Amount)

	if ack == pb.BidReply_SUCCESS {
		go a.replicateToBackups()
	}

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

	lastBid, existed := a.bids[bidder]
	if !existed {
		fmt.Printf("New bidder: %s\n", bidder)
	}

	if amount <= lastBid {
		fmt.Printf("The bid %d was not higehr than your last bid %d\n", amount, lastBid)
		return pb.BidReply_FAIL
	}

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

func (a *auctionServer) setState(highestBid int32, highestBidder string, startTime time.Time, bids map[string]int32) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.highestBid = highestBid
	a.highestBidder = highestBidder
	a.startTime = startTime
	a.bids = bids
}

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

func (a *auctionServer) connectToReplicas(replicaAddresses []string) {
	a.replicaMu.Lock()
	defer a.replicaMu.Unlock()

	for _, addr := range replicaAddresses {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("Failed to connect to replica: %s\n", addr)
			continue
		}
		client := pb.NewReplicaClient(conn)
		a.replicaClients[addr] = client
		fmt.Printf("Connected to replica: %s\n", addr)
	}
}

func (a *auctionServer) replicateToBackups() {
	a.mu.Lock()
	state := &pb.ReplicaBidState{
		HighestBid:  a.highestBid,
		AuctionOver: a.isOver(),
		Winner:      a.highestBidder,
	}
	a.mu.Unlock()

	a.replicaMu.Lock()
	defer a.replicaMu.Unlock()

	for addr, client := range a.replicaClients {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_, err := client.ReplicateBid(ctx, state)
		cancel()

		if err != nil {
			fmt.Printf("Failed to replicate bid: %s\n", addr)
		} else {
			fmt.Printf("Successfully replicated bid: %s\n", addr)
		}
	}
}
