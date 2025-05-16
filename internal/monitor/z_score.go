package monitor

import (
	"math"
)

var GlobalZScoreCalculator *ZScoreCalculator

// ZScoreCalculator 用于计算 Z-score 的结构体，包含各个指标的历史数据
type ZScoreCalculator struct {
	BatteryHistoryData      *HistoryData
	CPUHistoryData          *HistoryData
	MemoryHistoryData       *HistoryData
	NetworkDelayHistoryData *HistoryData
	DistanceHistoryData     *HistoryData
}

type HistoryData struct {
	Data             []float64 // 存储历史数据
	HistorySize      int       // 历史数据的大小
	CurrentDataIndex int       // 当前数据的索引
}

// InitZScoreCalculator 初始化 ZScoreCalculator 实例
func InitZScoreCalculator(historySize int) {
	GlobalZScoreCalculator = &ZScoreCalculator{
		BatteryHistoryData: &HistoryData{
			Data:             make([]float64, historySize),
			HistorySize:      historySize,
			CurrentDataIndex: 0,
		},
		CPUHistoryData: &HistoryData{
			Data:             make([]float64, historySize),
			HistorySize:      historySize,
			CurrentDataIndex: 0,
		},
		MemoryHistoryData: &HistoryData{
			Data:             make([]float64, historySize),
			HistorySize:      historySize,
			CurrentDataIndex: 0,
		},
		NetworkDelayHistoryData: &HistoryData{
			Data:             make([]float64, historySize),
			HistorySize:      historySize,
			CurrentDataIndex: 0,
		},
		DistanceHistoryData: &HistoryData{
			Data:             make([]float64, historySize),
			HistorySize:      historySize,
			CurrentDataIndex: 0,
		},
	}
}

// UpdateBatteryHistoryData 更新历史数据
func (z *ZScoreCalculator) UpdateBatteryHistoryData(newData float64) {
	z.UpdateHistoryData(newData, z.BatteryHistoryData)
}

func (z *ZScoreCalculator) UpdateCPUHistoryData(newData float64) {
	z.UpdateHistoryData(newData, z.CPUHistoryData)
}

func (z *ZScoreCalculator) UpdateMemoryHistoryData(newData float64) {
	z.UpdateHistoryData(newData, z.MemoryHistoryData)
}

func (z *ZScoreCalculator) UpdateNetworkDelayHistoryData(newData float64) {
	z.UpdateHistoryData(newData, z.NetworkDelayHistoryData)
}

func (z *ZScoreCalculator) UpdateDistanceHistoryData(newData float64) {
	z.UpdateHistoryData(newData, z.DistanceHistoryData)
}

// UpdateHistoryData 通用的更新方法，更新历史数据并循环更新索引
func (z *ZScoreCalculator) UpdateHistoryData(newData float64, historyData *HistoryData) {
	// 更新历史数据并循环更新索引
	historyData.Data[historyData.CurrentDataIndex] = newData
	historyData.CurrentDataIndex = (historyData.CurrentDataIndex + 1) % historyData.HistorySize
}

// Mean 计算均值
func (z *ZScoreCalculator) Mean(historyData *HistoryData) float64 {
	var sum float64
	for _, value := range historyData.Data {
		sum += value
	}
	return sum / float64(historyData.HistorySize)
}

// StandardDeviation 计算标准差
func (z *ZScoreCalculator) StandardDeviation(historyData *HistoryData) float64 {
	mean := z.Mean(historyData)
	var sumSquares float64
	for _, value := range historyData.Data {
		sumSquares += (value - mean) * (value - mean)
	}
	return math.Sqrt(sumSquares / float64(historyData.HistorySize))
}

// CalculateBatteryZScore 计算 Z-score
func (z *ZScoreCalculator) CalculateBatteryZScore(newData float64) float64 {
	// 电池电量是反向度量，电量越低，Z-score 越大
	return z.CalculateZScore(newData, z.BatteryHistoryData, true)
}

func (z *ZScoreCalculator) CalculateCPUZScore(newData float64) float64 {
	return z.CalculateZScore(newData, z.CPUHistoryData, false)
}

func (z *ZScoreCalculator) CalculateMemoryZScore(newData float64) float64 {
	return z.CalculateZScore(newData, z.MemoryHistoryData, false)
}

func (z *ZScoreCalculator) CalculateNetworkDelayZScore(newData float64) float64 {
	return z.CalculateZScore(newData, z.NetworkDelayHistoryData, false)
}

func (z *ZScoreCalculator) CalculateDistanceZScore(newData float64) float64 {
	return z.CalculateZScore(newData, z.DistanceHistoryData, false)
}

// CalculateZScore 计算 Z-score 的通用方法
func (z *ZScoreCalculator) CalculateZScore(newData float64, historyData *HistoryData, isInverse bool) float64 {
	mean := z.Mean(historyData)
	stdDev := z.StandardDeviation(historyData)

	if stdDev == 0 {
		// 防止标准差为0
		return 0
	}

	// 反向度量（如电池电量），电量越低，Z-score应该越大
	if isInverse {
		return 1 - (newData-mean)/stdDev
	}

	// 正常度量
	return (newData - mean) / stdDev
}
