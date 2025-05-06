package raft

import (
	"Orosync/internal/model"
	"Orosync/internal/rpc/pb/raft"
	"context"
	"fmt"
)

const (
	Leader int64 = iota
	Follower
	Candidate

	StartStatus = "start"
	StopStatus  = "stop"

	SleepTime     = 300
	AppendLogTime = 150

	SuccessCode = "0"

	FailCode = "500"

	NotLeaderCode   = "400"
	NotFollower     = "401"
	LeaderOutOfDate = "402"

	AddUAV    = 1
	RemoveUAV = 2
)

// Node   RaftNode节点
type Node struct {
	// base
	Role         int64 // 无人机角色：Leader、Follower、Candidate
	Term         int64 // 当前的任期号
	LastLogIndex int64 // 自己最后一个日志号
	LastLogTerm  int64 // 自己最后一个日志的任期号

	// leader
	LeaderUid    string // leader的Uid
	PrevLogIndex int64  // 前一个日志的日志号
	PrevLogTerm  int64  // 前一个日志的任期号
	LeaderCommit int64  // leader已提交的日志号
	HasHeartbeat bool   // leader有心跳
	IsVoted      bool   // 是否已经投票

	// log
	NewLogIndex int64
	Logs        LogsInfo

	// cluster
	Cluster ClusterInfo

	// UAV
	UAV *model.UAV
}

type ClusterInfo struct {
	UidList []string
}

type LogsInfo struct {
	index   int64                // 日志的索引
	Count   int64                // 无人机数量
	UidList []string             // 无人机uid列表
	UavMap  map[string]model.UAV // 日志
}

var GlobalNode Node

func InitGlobalNode() {
	GlobalNode = Node{
		Role: Follower,
		UAV:  model.NewUAV(),
	}
}

func (n *Node) ReceiveLog(ctx context.Context, req *raft.AppendLogReq) (*raft.AppendLogResp, error) {
	switch n.Role {
	case Leader:
		return GlobalLeader.ReceiveLog(ctx, req)
	case Follower:
		return GlobalFollower.ReceiveLog(ctx, req)
	case Candidate:
		return GlobalCandidate.ReceiveLog(ctx, req)
	default:
		return nil, fmt.Errorf("wrong role")
	}
}

func (n *Node) ReceiveVote(ctx context.Context, req *raft.VoteReq) (*raft.VoteResp, error) {
	switch n.Role {
	case Leader:
		return GlobalLeader.ReceiveVote(ctx, req)
	case Follower:
		return GlobalFollower.ReceiveVote(ctx, req)
	case Candidate:
		return GlobalCandidate.ReceiveVote(ctx, req)
	default:
		return nil, fmt.Errorf("wrong role")
	}
}

func (n *Node) ReceiveLogFromEachUAV(ctx context.Context, req *raft.SendUAVInfoReq) (*raft.SendUAVInfoResp, error) {
	switch n.Role {
	case Leader:
		return GlobalLeader.ReceiveLogFromEachUAV(ctx, req)
	case Follower:
		return GlobalFollower.ReceiveLogFromEachUAV(ctx, req)
	case Candidate:
		return GlobalCandidate.ReceiveLogFromEachUAV(ctx, req)
	default:
		return nil, fmt.Errorf("wrong role")
	}
}
