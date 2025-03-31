package model

// Cluster 集群信息
type Cluster struct {
	Count        int8  `json:"count"`          // 无人机数量
	UAVs         []UAV `json:"uvas"`           // 无人机列表
	LastLogIndex int64 `json:"last_log_index"` // 最后一个日志号
	LastLogTerm  int64 `json:"last_log_term"`  // 最后一个日志的任期号
}
