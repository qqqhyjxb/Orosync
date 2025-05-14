package raft

import (
	"Orosync/internal/client"
	"Orosync/internal/rpc/pb/raft"
	"context"
	"fmt"
)

var GlobalCandidate *CandidateInstance

type CandidateInstance struct {
	Status string
}

func InitGlobalCandidate() {
	GlobalCandidate = &CandidateInstance{
		Status: StopStatus,
	}
}

func (c *CandidateInstance) Start() {
	c.Status = StartStatus
	GlobalNode.Role = Candidate

	fmt.Printf("new candidate: %v\n", GlobalNode.UAV.Uid)

	ctx := context.Background()
	go c.Vote(ctx)
}

func (c *CandidateInstance) Stop() {
	c.Status = StopStatus
}

func (c *CandidateInstance) Vote(ctx context.Context) {
	for GlobalNode.Role == Candidate {

		var voteCount int64

		GlobalNode.Term = GlobalNode.Term + 1

		request := &raft.VoteReq{
			Term:         GlobalNode.Term,
			CandidateUid: GlobalNode.UAV.Uid,
			LastLogIndex: GlobalNode.LastLogIndex,
			LastLogTerm:  GlobalNode.LastLogTerm,
		}

		for _, u := range GlobalNode.Logs.UavMap {
			// 跳过自己
			if u.Uid == GlobalNode.UAV.Uid {
				voteCount++
				continue
			}

			// 获取客户端连接（自动复用）
			clientObj, err := client.GlobalRaftClient.StartClient(
				u.Address,
			)
			if err != nil {
				fmt.Printf("连接失败: %v\n", err)
				continue
			}

			response, err := clientObj.Client.Vote(ctx, request)
			if err != nil {
				// 投票错误直接跳过
				fmt.Printf("uav: %v start vote err: %v  address: %v\n", GlobalNode.UAV.Uid, err, u.Address)
				continue
			}

			fmt.Printf("u.uid: %v  uav.uid: %v\n", u.Uid, GlobalNode.UAV.Uid)
			fmt.Printf("candidate response.Term: %v\n", response.Term)
			fmt.Printf("candidate response.LeaderUid: %v\n", response.LeaderUid)
			fmt.Printf("candidate response.VoteGranted: %v\n", response.VoteGranted)

			// 已有新的leader,该candidate退出follower
			if response.Term > GlobalNode.Term {
				GlobalNode.Role = Follower
				GlobalNode.Term = response.Term
				GlobalNode.LeaderUid = response.LeaderUid

				GlobalCandidate.Stop()

				GlobalFollower.Start()

				return
			}

			// 收到投票
			if response.VoteGranted {
				voteCount++
			}
		}

		if voteCount > GlobalNode.Logs.Count/2 {
			GlobalNode.Role = Leader

			GlobalCandidate.Stop()

			GlobalLeader.Start()

			return
		}

		fmt.Printf("fail to candidate: %v  voteCount: %v\n", GlobalNode.UAV.Uid, voteCount)
	}
}

func (c *CandidateInstance) ReceiveVote(ctx context.Context, req *raft.VoteReq) (*raft.VoteResp, error) {
	resp := &raft.VoteResp{}

	if GlobalNode.Term >= req.Term {
		resp.VoteGranted = false
		resp.Term = GlobalNode.Term
		resp.LeaderUid = GlobalNode.UAV.Uid
		return resp, nil
	}

	c.Stop()

	go GlobalFollower.Start()

	resp.VoteGranted = false
	return resp, nil
}

func (c *CandidateInstance) ReceiveLogFromEachUAV(ctx context.Context, req *raft.SendUAVInfoReq) (*raft.SendUAVInfoResp, error) {
	resp := &raft.SendUAVInfoResp{
		Code:      NotLeaderCode,
		LeaderUid: "",
	}

	return resp, nil
}

func (c *CandidateInstance) ReceiveLog(ctx context.Context, req *raft.AppendLogReq) (*raft.AppendLogResp, error) {
	resp := &raft.AppendLogResp{
		Code:      NotFollower,
		Term:      GlobalNode.Term,
		LeaderUid: "",
	}

	if req.Term > GlobalNode.Term {
		GlobalNode.Term = req.Term
		GlobalNode.LeaderUid = req.LeaderUid

		c.Stop()

		go GlobalFollower.Start()
	}

	return resp, nil
}

//func (c *CandidateInstance) GlobalLoadBalance(ctx context.Context, req *raft.GlobalLoadBalanceReq) (*raft.GlobalLoadBalanceResp, error) {
//	resp := &raft.GlobalLoadBalanceResp{
//		Code: NotLeaderCode,
//	}
//
//	return resp, nil
//}
