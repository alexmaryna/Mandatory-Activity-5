package Mandatory_Activity_5

import (
	"time"
)

type role int

const (
	roleLeader role = iota
	roleFollower
)

const auctionDuration = 100 * time.Second

type auctionState struct {
	start         time.Time
	end           time.Time
	registered    bool
	highestBid    int
	highestBidder string
	closed        bool
}

type resultView struct {
	finished   bool
	highestBid int
	winnerId   string
}

type logEntry struct {
	timestamp int64
	bidderId  string
	amount    int
}

type node struct {
	id      int
	role    role
	address string

	auction *auctionState

	lamportClock int64
	log          []logEntry
}

func newAuctionState() *auctionState {
	start := time.Now()
	return &auctionState{
		start: start,
		end:   start.Add(auctionDuration),
	}
}

// Bid
func (a *auctionState) clientRequest(amount int, bidderId string) (string, error) {
	now := time.Now()

	if a.closed || now.After(a.end) {
		a.closed = true
		return "fail", nil
	}

	if amount <= a.highestBid {
		return "fail", nil
	}

	a.registered = true
	a.highestBid = amount
	a.highestBidder = bidderId

	return "success", nil
}

// Result
func (a *auctionState) clientResponse() (resultView, error) {
	now := time.Now()

	if now.After(a.end) {
		a.closed = true
		return resultView{
			finished:   true,
			highestBid: a.highestBid,
			winnerId:   a.highestBidder,
		}, nil
	}

	return resultView{
		finished:   false,
		highestBid: a.highestBid,
		winnerId:   "",
	}, nil
}

//lamport helpers

// new client bid ont he leader
func (n *node) nextLamport() int64 {
	n.lamportClock++
	return n.lamportClock
}

// receiving of the replicated bid from another node
func (n *node) updateLamport(remote int64) {
	if remote > n.lamportClock {
		n.lamportClock = remote
	}
	n.lamportClock++
}

// appends to the log and applies to auctionState using clientRequest
func (n *node) commitBid(entry logEntry) (string, error) {
	n.log = append(n.log, entry)

	return n.auction.clientRequest(entry.amount, entry.bidderId)
}

// for the gRPC server to get the result
func (n *node) commitResult() (resultView, error) {
	return n.auction.clientResponse()
}

//for node creation in main
func newNode(id int, address string) *node {
	return &node{
		id:      id,
		role:    roleFollower,
		address: address,
		auction: newAuctionState(),
		log:     make([]logEntry, 0),
	}
}
