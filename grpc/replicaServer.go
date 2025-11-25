package grpc

import (
	proto "Mandatory-Activity-5/grpc"
	"context"
)

type ReplicaServer struct {
	proto.UnimplementedReplicaServer

	HighestBid  int32
	AuctionOver bool
	Winner      string
}

func NewReplicaServer() *ReplicaServer {
	return &ReplicaServer{
		HighestBid:  0,
		AuctionOver: false,
		Winner:      "",
	}
}

func (r *ReplicaServer) ReplicateBid(ctx context.Context, in *proto.ReplicaBidState) (*proto.ReplicaAck, error) {
	r.HighestBid = in.HighestBid
	r.AuctionOver = in.AuctionOver
	r.Winner = in.Winner
	return &proto.ReplicaAck{Ok: true}, nil
}

func (r *ReplicaServer) SyncState(ctx context.Context, in *proto.SyncRequest) (*proto.SyncReply, error) {
	reply := &proto.SyncReply{
		HighestBid:  r.HighestBid,
		AuctionOver: r.AuctionOver,
		Winner:      r.Winner,
	}
	return reply, nil
}
