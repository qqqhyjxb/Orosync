package raft

func InitRaft() {
	// 初始化全局节点状态
	InitGlobalCandidate()
	InitGlobalLeader()
	InitGlobalFollower()
	InitGlobalNode()

	// 节点启动
	GlobalNode.Start()
}

func (n *Node) Start() {
	n.Role = Follower
	GlobalFollower.Start()

}

/*
	TODO: 领导者持续收集日志并更新到Log（被动型） 持续追加日志（主动型 - 心跳机制）
	TODO：跟随者持续更新自身状态到领导者（主动型） 跟随者持续监测心跳并决定是否启动选举（主动型） 跟随者持续接收最新的日志（被动型）
    TODO: 领导者对于集群节点变更的处理
*/
