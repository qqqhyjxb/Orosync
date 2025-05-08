package model

var (
	GlobalUAV UAV // 无人机信息
)

func InitUAV() {
	GlobalUAV = UAV{}
}
