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

	auctionServer *auctionServer
}

func newReplicaServer(auctionServer *auctionServer) *replicaServer {
	return &replicaServer{
		auctionServer: auctionServer,
	}
}

func (r *replicaServer) ReplicateBid(ctx context.Context, in *pb.ReplicaBidState) (*pb.ReplicaAck, error) {
	r.mu.Lock()
	r.highestBid = in.HighestBid
	r.auctionOver = in.AuctionOver
	r.winner = in.Winner
	r.mu.Unlock()

	if r.auctionServer != nil {
		r.auctionServer.mu.Lock()
		r.auctionServer.highestBid = in.HighestBid
		r.auctionServer.highestBidder = in.Winner
		r.auctionServer.mu.Unlock()
	}

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
