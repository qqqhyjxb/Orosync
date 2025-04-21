package raft

import (
	"Orosync/internal/client"
	"Orosync/internal/rpc/pb/raft"
	"context"
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

	ctx := context.Background()
	go c.Vote(ctx)
}

func (c *CandidateInstance) Stop() {
	c.Status = StopStatus
}

func (c *CandidateInstance) Vote(ctx context.Context) {
	var voteCount int64

	request := &raft.VoteReq{
		Term:         GlobalNode.Term,
		CandidateUid: GlobalNode.UAV.Uid,
		LastLogIndex: GlobalNode.LastLogIndex,
		LastLogTerm:  GlobalNode.LastLogTerm,
	}

	for _, u := range GlobalNode.Logs.UavMap {

		// TODO:address需要再组装一下
		response, err := client.GlobalRaftClient.StartClient(u.Address).Client.Vote(ctx, request)
		if err != nil {
			// 投票错误直接跳过
			continue
		}

		// 已有新的leader,该candidate退出follower
		if response.Term > GlobalNode.Term {
			GlobalNode.Role = Follower
			GlobalNode.Term = response.Term
			GlobalNode.LeaderUid = response.LeaderUid

			GlobalCandidate.Stop()

			GlobalFollower.Start()
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

	resp.VoteGranted = true
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

		go GlobalCandidate.Start()
	}

	return resp, nil
}
