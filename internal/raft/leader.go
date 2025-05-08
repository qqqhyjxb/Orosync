package raft

import (
	"Orosync/internal/client"
	"Orosync/internal/rpc/pb/raft"
	"Orosync/internal/rpc/pb/simulation"
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"log"
	"time"
)

var GlobalLeader *LeaderInstance

type LeaderInstance struct {
	Status string
}

func InitGlobalLeader() {
	GlobalLeader = &LeaderInstance{
		Status: StopStatus,
	}
}

// Start 开始leader的新周期
func (l *LeaderInstance) Start() {
	l.Status = StartStatus
	GlobalNode.Role = Leader

	fmt.Printf("new leader: %v\n", GlobalNode.UAV.Uid)

	ctx := context.Background()

	go l.AppendLogs(ctx)
	go l.UpdateAndPrintLogs()
}

func (l *LeaderInstance) Stop() {
	l.Status = StopStatus
}

// AppendLogs  追加日志
func (l *LeaderInstance) AppendLogs(ctx context.Context) {
	for GlobalNode.Role == Leader {
		// 周期性追加日志
		time.Sleep(AppendLogTime * time.Millisecond)

		logs := &raft.LogsInfo{}

		err := copier.Copy(logs, GlobalNode.Logs)
		if err != nil {
			log.Printf("copier.Copy LogsInfo err: %v", err)
		}

		request := &raft.AppendLogReq{
			Term:         GlobalNode.Term,
			LeaderUid:    GlobalNode.UAV.Uid,
			PrevLogIndex: GlobalNode.PrevLogIndex,
			PrevLogTerm:  GlobalNode.PrevLogTerm,
			Logs:         logs,
			LeaderCommit: GlobalNode.LeaderCommit,
		}

		for _, u := range GlobalNode.Logs.UavMap {

			// 跳过自己
			if u.Uid == GlobalNode.UAV.Uid {
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

			response, err := clientObj.Client.AppendLog(ctx, request)
			if err != nil {
				// TODO：考虑无人机一直追加日志失败有没有兜底方案
				continue
			}

			// 响应的term更大，代表有新leader产生，老leader应该退位
			if response.Term > GlobalNode.Term {
				GlobalNode.Role = Follower
				GlobalNode.Term = response.Term
				GlobalNode.LeaderUid = response.LeaderUid

				GlobalLeader.Stop()

				GlobalFollower.Start()

				return
			}
		}
	}
}

// ReceiveLogFromEachUAV 从每台无人机
func (l *LeaderInstance) ReceiveLogFromEachUAV(ctx context.Context, req *raft.SendUAVInfoReq) (*raft.SendUAVInfoResp, error) {
	if GlobalLeader.Status == StopStatus {
		return nil, fmt.Errorf("uav is not leader")
	}

	resp := &raft.SendUAVInfoResp{}

	// TODO:目前的方案是直接更新在logs中，然后下一次追加日志给各无人机也直接发logs就行

	// 存在该无人机才赋值
	if v, ok := GlobalNode.Logs.UavMap[req.Uav.Uid]; ok {
		err := copier.Copy(v, req.Uav)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("uav:%s is not exist  req.uid: %v", req.Uav.Uid, req.Uav.Uid)
	}

	resp.Code = SuccessCode
	return resp, nil
}

func (l *LeaderInstance) DroneSwarmChange(ctx context.Context,
	req *simulation.DroneSwarmChangeRequest) (*simulation.DroneSwarmChangeResponse, error) {

	resp := &simulation.DroneSwarmChangeResponse{}

	//if GlobalLeader.Status == StopStatus {
	//	resp.Code = NotLeaderCode
	//	resp.Msg = "uav is not leader"
	//	return resp, nil
	//}
	//
	//if req.ChangeType == AddUAV {
	//	for _, u := range req.DroneUidList {
	//		GlobalNode.Cluster.UidList = append(GlobalNode.Cluster.UidList, u)
	//	}
	//} else if req.ChangeType == RemoveUAV {
	//	var newUidList []string
	//	var removeUidMap map[string]bool
	//
	//	for _, u := range req.DroneUidList {
	//		removeUidMap[u] = true
	//	}
	//
	//	for _, u := range GlobalNode.Cluster.UidList {
	//		if removeUidMap[u] {
	//			continue
	//		}
	//		newUidList = append(newUidList, u)
	//	}
	//
	//	GlobalNode.Cluster.UidList = newUidList
	//} else {
	//	resp.Code = FailCode
	//	resp.Msg = "wrong change type"
	//}

	return resp, nil
}

func (l *LeaderInstance) ReceiveVote(ctx context.Context, req *raft.VoteReq) (*raft.VoteResp, error) {
	resp := &raft.VoteResp{}

	if GlobalNode.Term >= req.Term {
		resp.VoteGranted = false
		resp.Term = GlobalNode.Term
		resp.LeaderUid = GlobalNode.UAV.Uid
		return resp, nil
	}

	l.Stop()

	go GlobalFollower.Start()

	resp.VoteGranted = false
	return resp, nil
}

func (l *LeaderInstance) ReceiveLog(ctx context.Context, req *raft.AppendLogReq) (*raft.AppendLogResp, error) {
	resp := &raft.AppendLogResp{
		Code:      NotFollower,
		Term:      GlobalNode.Term,
		LeaderUid: GlobalNode.LeaderUid,
	}

	if req.Term > GlobalNode.Term {
		GlobalNode.Term = req.Term
		GlobalNode.LeaderUid = req.LeaderUid

		l.Stop()

		go GlobalFollower.Start()
	}

	return resp, nil
}

func (l *LeaderInstance) UpdateAndPrintLogs() {
	for GlobalNode.Role == Leader {

		GlobalNode.UAV.Time = time.Now().String()

		if v, ok := GlobalNode.Logs.UavMap[GlobalNode.UAV.Uid]; ok {
			err := copier.Copy(v, GlobalNode.UAV)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("uav:%s is not exist\n", GlobalNode.UAV.Uid)
		}

		time.Sleep(2 * time.Second)

		fmt.Printf("current leader: %s\n", GlobalNode.UAV.Uid)
		fmt.Printf("current term: %d\n", GlobalNode.Term)
		for _, u := range GlobalNode.Logs.UavMap {
			fmt.Println(u)
		}

		fmt.Printf("\n\n\n\n")
	}
}
