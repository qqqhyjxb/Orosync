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

	SleepTime     = 400
	AppendLogTime = 100

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
	LastVoteTerm int64  // 上一次投票的任期

	// log
	NewLogIndex int64
	Logs        *LogsInfo

	// cluster
	LocalGroup *GroupInfo

	// UAV
	UAV *model.UAV

	// 各项指标阈值
	BatteryThreshold  float32 //越小越靠近阈值
	CpuUsageThreshold float32 //越大越靠近阈值
	MemoryThreshold   float32 //越小越靠近阈值
	DistanceThreshold float32 //越大越靠近阈值

	//
	SimulationAddress string
}

type GroupInfo struct {
	UidList []string
}

type LogsInfo struct {
	Index   int64                 `yaml:"index"`
	Count   int64                 `yaml:"count"`
	UidList []string              `yaml:"uid_list"`
	UavMap  map[string]*model.UAV `yaml:"uav_map"`
}

var GlobalNode Node

func InitGlobalNode(uav *model.UAV, logs *LogsInfo) {
	GlobalNode = Node{
		Role: Follower,
		UAV:  uav,
		Logs: logs,
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

//func (n *Node) GlobalLoadBalance(ctx context.Context, req *raft.GlobalLoadBalanceReq) (*raft.GlobalLoadBalanceResp, error) {
//	switch n.Role {
//	case Leader:
//		return GlobalLeader.GlobalLoadBalance(ctx, req)
//	case Follower:
//		return GlobalFollower.GlobalLoadBalance(ctx, req)
//	case Candidate:
//		return GlobalCandidate.GlobalLoadBalance(ctx, req)
//	default:
//		return nil, fmt.Errorf("wrong role")
//	}
//}
