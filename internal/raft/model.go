package raft

// RequestVoteRequest 请求投票RPC Request
type RequestVoteRequest struct {
	Term         int // 自己当前的任期号
	CandidateUid int // 自己的Uid
	LastLogIndex int // 自动最后一个日志号
	LastLogTerm  int // 自己最后一个日志的任期号
}

// RequestVoteResponse 请求投票的RPC Response
type RequestVoteResponse struct {
	Term        int  // follower自己的任期号
	VoteGranted bool // follower是否选择投票给这个请求的candidate
}
