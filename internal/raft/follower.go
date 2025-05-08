package raft

import (
	"context"
	"fmt"
	"time"

	"Orosync/internal/client"
	"Orosync/internal/rpc/pb/raft"

	"github.com/jinzhu/copier"
)

var GlobalFollower *FollowerInstance

type FollowerInstance struct {
	Status string
}

func InitGlobalFollower() {
	GlobalFollower = &FollowerInstance{
		Status: StopStatus,
	}
}

// Start 开始follower的生命周期
func (f *FollowerInstance) Start() {
	f.Status = StartStatus
	GlobalNode.Role = Follower
	GlobalNode.IsVoted = false

	fmt.Printf("new follower: %v\n", GlobalNode.UAV.Uid)

	go f.sendLogToLeader()
	go f.heartbeat()
}

func (f *FollowerInstance) Stop() {
	f.Status = StopStatus
}

// heartbeat 检测leader心跳
func (f *FollowerInstance) heartbeat() {
	for GlobalNode.Role == Follower {
		time.Sleep(SleepTime * time.Millisecond) //心跳超时时间

		// leader心跳失效，转换为candidate开始运行
		if !GlobalNode.HasHeartbeat {
			GlobalNode.Role = Candidate

			f.Stop()

			GlobalCandidate.Start()

			return
		}

		// leader心跳正常，置监听心跳
		GlobalNode.HasHeartbeat = false
	}
}

// sendLogToLeader 定期发送自身状态给leader，无并发问题
func (f *FollowerInstance) sendLogToLeader() {
	for GlobalNode.Role == Follower {

		if GlobalNode.LeaderUid == "" {
			continue
		}

		time.Sleep(AppendLogTime * time.Millisecond) //发送log的周期

		ctx := context.Background()

		GlobalNode.UAV.Time = time.Now().String()

		uav := &raft.UAVInfo{}

		err := copier.Copy(uav, GlobalNode.UAV)
		if err != nil {
			return
		}

		req := &raft.SendUAVInfoReq{
			Uav: uav,
		}

		// 获取客户端连接（自动复用）
		clientObj, err := client.GlobalRaftClient.StartClient(
			GlobalNode.Logs.UavMap[GlobalNode.LeaderUid].Address,
		)
		if err != nil {
			fmt.Printf("连接失败: %v\n", err)
			continue
		}

		// 使用复用的连接发起请求
		resp, err := clientObj.Client.SendUAVInfo(ctx, req)
		if err != nil || resp.Code != SuccessCode {
			fmt.Printf("请求失败 uid:%v err:%v \n", GlobalNode.UAV.Uid, err)
		}
	}
}

// ReceiveLog log和心跳为一体
func (f *FollowerInstance) ReceiveLog(ctx context.Context, req *raft.AppendLogReq) (*raft.AppendLogResp, error) {
	resp := &raft.AppendLogResp{}

	// 安全性校验,log不需要完整一致性，只需要临时一致性，所以只要是最新的leader发过来的日志，无条件接受

	// 若发送请求的leader任期比节点当前任期小，代表该leader已经过期
	if req.Term < GlobalNode.Term {
		resp.Code = LeaderOutOfDate
		resp.Term = GlobalNode.Term
		resp.LeaderUid = GlobalNode.LeaderUid
		return resp, nil
	}

	// 接收日志
	err := copier.Copy(GlobalNode.Logs, req.Logs)
	if err != nil {
		resp.Code = FailCode
		resp.Term = GlobalNode.Term
		resp.LeaderUid = GlobalNode.LeaderUid
		return nil, err
	}

	// 检测到心跳
	GlobalNode.HasHeartbeat = true

	GlobalNode.LeaderUid = req.LeaderUid

	resp.Code = SuccessCode
	resp.Term = GlobalNode.Term
	resp.LeaderUid = GlobalNode.LeaderUid

	// 更新为未投票状态
	GlobalNode.IsVoted = false

	return resp, nil
}

func (f *FollowerInstance) ReceiveVote(ctx context.Context, req *raft.VoteReq) (*raft.VoteResp, error) {
	resp := &raft.VoteResp{}

	// 若请求投票的candidate任期比当前小，说明当前任期已经过期
	// 拒绝投票，并告诉candidate新leader的信息
	if req.Term < GlobalNode.Term {
		resp.Term = GlobalNode.Term
		resp.LeaderUid = GlobalNode.LeaderUid
		resp.VoteGranted = false
		return resp, nil
	}

	// 已经投过票了，拒绝投票
	if GlobalNode.IsVoted {
		resp.VoteGranted = false
		return resp, nil
	}

	// 同意投票
	resp.VoteGranted = true
	GlobalNode.IsVoted = true
	return resp, nil
}

func (f *FollowerInstance) ReceiveLogFromEachUAV(ctx context.Context, req *raft.SendUAVInfoReq) (*raft.SendUAVInfoResp, error) {
	resp := &raft.SendUAVInfoResp{
		Code:      NotLeaderCode,
		LeaderUid: GlobalNode.LeaderUid,
	}

	return resp, nil
}
