package model

// UAV 无人机信息
type UAV struct {
	// base
	Uid     string `json:"uid"` // 唯一ID
	Time    string `json:"time" `
	Address string `json:"ip"`
	Port    string `json:"port"`

	// state
	Position PositionInfo `json:"position"`
	Battery  BatteryInfo  `json:"battery"` //电量
	CPU      CPUInfo      `json:"cpu"`
	Memory   MemoryInfo   `json:"memory"`
	Network  NetworkInfo  `json:"network"`
	Tasks    TaskList     `json:"task"`
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

// TaskList 任务信息
type TaskList struct {
	List  []TaskInfo `json:"list"`  // 列表
	Count []int8     `json:"count"` // 任务数量
}

// TaskInfo 任务信息
type TaskInfo struct {
	// TODO: 完善任务具体信息
	Type string `json:"type"` // 任务类型
}

func NewUAV() *UAV {
	return &UAV{}
}
