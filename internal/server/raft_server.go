package server

import (
	"Orosync/internal/raft"
	raft2 "Orosync/internal/rpc/pb/raft"
	"context"
)

type RaftService struct {
	raft2.UnimplementedRaftServiceServer
}

func (r *RaftService) SendUAVInfo(ctx context.Context, req *raft2.SendUAVInfoReq) (*raft2.SendUAVInfoResp, error) {
	return raft.GlobalNode.ReceiveLogFromEachUAV(ctx, req)
}

func (r *RaftService) AppendLog(ctx context.Context, req *raft2.AppendLogReq) (*raft2.AppendLogResp, error) {
	return raft.GlobalNode.ReceiveLog(ctx, req)
}

func (r *RaftService) Vote(ctx context.Context, req *raft2.VoteReq) (*raft2.VoteResp, error) {
	return raft.GlobalNode.ReceiveVote(ctx, req)
}

//func (r *RaftService) GlobalLoadBalance(ctx context.Context, req *raft2.GlobalLoadBalanceReq) (*raft2.GlobalLoadBalanceResp, error) {
//	return raft.GlobalNode.GlobalLoadBalance(ctx, req)
//}
