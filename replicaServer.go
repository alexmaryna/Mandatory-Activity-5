package Mandatory_Activity_5

import (
	pb "Mandatory-Activity-5/grpc"
	"context"
	"sync"
)

type ReplicaServer struct {
	pb.UnimplementedReplicaServer

	auction *auctionServer
	mu      sync.Mutex

	//HighestBid  int32
	//AuctionOver bool
	//Winner      string
}

func NewReplicaServer(auction *auctionServer) *ReplicaServer {
	return &ReplicaServer{
		auction: auction,
		//HighestBid:  0,
		//AuctionOver: false,
		//Winner:      "",
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

type replicaServer struct {
}

func newReplicaServer() replicaServer {
	return replicaServer{}
}

func replicaBid() {

}

func syncState() {

}

func broadcastState() {

}

func heartbeat() {

}

func handleFaliure() {

}
