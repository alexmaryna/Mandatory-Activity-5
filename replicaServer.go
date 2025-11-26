package main

import (
	pb "Mandatory-Activity-5/grpc"
	"context"
	"sync"
)

type replicaServer struct {
	pb.UnimplementedReplicaServer

	mu          sync.Mutex
	highestBid  int32
	auctionOver bool
	winner      string
}

func newReplicaServer() *replicaServer {
	return &replicaServer{}
}

func (r *replicaServer) ReplicateBid(ctx context.Context, in *pb.ReplicaBidState) (*pb.ReplicaAck, error) {
	r.mu.Lock()
	r.highestBid = in.HighestBid
	r.auctionOver = in.AuctionOver
	r.winner = in.Winner
	r.mu.Unlock()

	return &pb.ReplicaAck{Ok: true}, nil
}

func (r *replicaServer) SyncState(ctx context.Context, in *pb.SyncRequest) (*pb.SyncReply, error) {
	r.mu.Lock()
	reply := &pb.SyncReply{
		HighestBid:  r.highestBid,
		AuctionOver: r.auctionOver,
		Winner:      r.winner,
	}
	r.mu.Unlock()
	return reply, nil
}

/*func (r *replicaServer) update(highestBid int32, auctionOver bool, winner string) {
	r.mu.Lock()
	r.highestBid = highestBid
	r.auctionOver = auctionOver
	r.winner = winner
	r.mu.Unlock()
}

func (r *replicaServer) getState() (int32, bool, string) {
	r.mu.Lock()
	b := r.highestBid
	o := r.auctionOver
	w := r.winner
	r.mu.Unlock()
	return b, o, w
}*/
