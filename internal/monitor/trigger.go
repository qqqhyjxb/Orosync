package monitor

import (
	"Orosync/internal/raft"
	"fmt"
	"time"
)

var GlobalTriggerMechanism *TriggerMechanism

// TriggerMechanism 触发机制结构体
type TriggerMechanism struct {
	// 存储每个指标的阈值
	Thresholds     []float64 // 指标阈值
	Weights        []float64 // 各指标权重
	ScoreThreshold float64   // 综合评分阈值
	LagCount       int       // 连续触发次数
	TriggerHistory []bool    // 记录历史触发状态（最多保存三次）
}

// InitTriggerMechanism 初始化触发机制
func InitTriggerMechanism(weights []float64, thresholds []float64, scoreThreshold float64, lagCount int) {
	GlobalTriggerMechanism = &TriggerMechanism{
		Thresholds:     thresholds,
		Weights:        weights,
		ScoreThreshold: scoreThreshold,
		LagCount:       lagCount,
		TriggerHistory: make([]bool, lagCount), // 滞后区间记录，最多保存lagCount次
	}
}

// CalculateCompositeScore 计算综合评分
func (t *TriggerMechanism) CalculateCompositeScore(metrics []float64) float64 {
	var score float64

	score += GlobalZScoreCalculator.CalculateBatteryZScore(metrics[0]) * t.Weights[0]
	score += GlobalZScoreCalculator.CalculateCPUZScore(metrics[1]) * t.Weights[1]
	score += GlobalZScoreCalculator.CalculateMemoryZScore(metrics[2]) * t.Weights[2]
	score += GlobalZScoreCalculator.CalculateNetworkDelayZScore(metrics[3]) * t.Weights[3]
	score += GlobalZScoreCalculator.CalculateDistanceZScore(metrics[4]) * t.Weights[4]

	return score
}

// LevelOneTrigger 第一级触发机制（单指标触发）
func (t *TriggerMechanism) LevelOneTrigger(metrics []float64) bool {
	for i, metric := range metrics {
		// 电量： 1 - 当前电量
		if i == 0 {
			metric = 1 - metric
		}

		if metric >= t.Thresholds[i] {
			fmt.Printf("Level One Trigger: Indicator %d exceeded threshold (%.2f >= %.2f)\n", i+1, metric, t.Thresholds[i])
			return true
		}
	}
	return false
}

// LevelTwoTrigger 第二级触发机制（综合评分触发）
func (t *TriggerMechanism) LevelTwoTrigger(metrics []float64) bool {
	score := t.CalculateCompositeScore(metrics)
	if score >= t.ScoreThreshold {
		fmt.Printf("Level Two Trigger: Composite score (%.2f) exceeded threshold (%.2f)\n", score, t.ScoreThreshold)
		return true
	}
	return false
}

// CheckLagPeriod 滞后区间触发机制（连续多次触发）
func (t *TriggerMechanism) CheckLagPeriod() bool {
	// 检查是否已连续满足触发条件
	if len(t.TriggerHistory) == t.LagCount {
		count := 0
		for _, triggered := range t.TriggerHistory {
			if triggered {
				count++
			}
		}
		return count == t.LagCount // 如果历史记录中连续满足条件的次数等于lagCount，则触发
	}
	return false
}

// EvaluateTrigger 执行触发判定
func (t *TriggerMechanism) EvaluateTrigger() {
	for {
		// 每300毫秒触发一次检查
		time.Sleep(300 * time.Millisecond)

		var metrics []float64
		metrics = append(metrics, float64(raft.GlobalNode.UAV.Battery.Capacity))
		metrics = append(metrics, float64(raft.GlobalNode.UAV.CPU.UsageRate))
		metrics = append(metrics, float64(raft.GlobalNode.UAV.Memory.UsageRate))
		metrics = append(metrics, float64(raft.GlobalNode.UAV.Network.Delay))
		metrics = append(metrics, float64(0))

		// 检查第一级触发
		levelOneTriggered := t.LevelOneTrigger(metrics)

		// 检查第二级触发
		levelTwoTriggered := t.LevelTwoTrigger(metrics)

		// 如果任何一级触发，则更新历史并检查滞后区间
		if levelOneTriggered || levelTwoTriggered {
			t.TriggerHistory = append(t.TriggerHistory[1:], true) // 记录当前触发
			if t.CheckLagPeriod() {
				fmt.Println("Lag Period Trigger: Executing rebalancing task due to consecutive triggers.")
				// TODO：触发负载均衡
			}
		} else {
			// 如果没有触发，更新历史状态
			t.TriggerHistory = append(t.TriggerHistory[1:], false)
		}

		// 该轮次检查没有触发负载均衡
	}
}
