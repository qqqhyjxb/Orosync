package raft

import (
	"Orosync/internal/model"
	"time"
)

func InitRaft(uav *model.UAV, logs *LogsInfo) {
	// 初始化全局节点状态
	InitGlobalCandidate()
	InitGlobalLeader()
	InitGlobalFollower()
	InitGlobalNode(uav, logs)
}

func (n *Node) Start() {
	time.Sleep(10 * time.Second)

	//fmt.Println(n.UAV.Uid)
	//fmt.Println(n.UAV.Address)
	//fmt.Println(n.UAV.Port)

	n.Role = Follower
	GlobalFollower.Start()
}
