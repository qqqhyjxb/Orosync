package model

// UAV 无人机信息
type UAV struct {
	// base
	Uid     string `json:"uid"` // 唯一ID
	Time    string `json:"time" `
	Address string `json:"address"`
	Port    int    `json:"port"`

	// state
	Position PositionInfo `json:"position"`
	Battery  BatteryInfo  `json:"battery"` //电量
	CPU      CPUInfo      `json:"cpu"`
	Memory   MemoryInfo   `json:"memory"`
	Network  NetworkInfo  `json:"network"`
	Tasks    []TaskInfo   `json:"tasks"`
	Status   string       `json:"status"`
}

type PositionInfo struct {
	X float32
	Y float32
	Z float32
}

type BatteryInfo struct {
	Capacity float32 `json:"capacity"` // 电量
}

// CPUInfo cpu信息
type CPUInfo struct {
	Capacity        float32 `json:"capacity"`         // cpu处理能力
	SurplusCapacity float32 `json:"surplus_capacity"` //剩余资源
	UsageRate       float32 `json:"usage_rate"`       // cpu使用率
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Capacity        float32 `json:"capacity"`         // 内存容量
	SurplusCapacity float32 `json:"surplus_capacity"` // 剩余容量
	UsageRate       float32 `json:"usage_rate"`       // 内存使用率
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	Bandwidth float32 `json:"bandwidth"` // 网络带宽
	Delay     float32 `json:"delay"`     // 网络延迟
}

// TaskInfo 任务信息
type TaskInfo struct {
	TaskId        int32              `json:"task_id"`
	Type          string             `json:"type"` // 任务类型
	Priority      int32              `json:"priority"`
	RequiredCPU   float32            `json:"required_cpu"`
	CpuUsage      float32            `json:"cpu_usage"`
	RequireMemory float32            `json:"require_memory"`
	Cost          map[string]float32 `json:"cost"`
	Target        PositionInfo       `json:"target"`
	TaskParents   []int32            `json:"task_parents"`
	TaskChildren  []int32            `json:"task_children"`
	Status        string             `json:"status"`
	Duration      string             `json:"duration"`
	Surplus       string             `json:"surplus"`
}

func NewUAV() *UAV {
	return &UAV{}
}
