package model

// Cluster 集群信息
type Cluster struct {
	Count int8  `json:"count"` // 无人机数量
	UAVs  []UAV `json:"uvas"`  // 无人机列表
}
