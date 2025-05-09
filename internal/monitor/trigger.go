package monitor

import (
	"fmt"
)

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
func InitTriggerMechanism(weights []float64, thresholds []float64, scoreThreshold float64, lagCount int) *TriggerMechanism {
	return &TriggerMechanism{
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
	for i, metric := range metrics {
		score += metric * t.Weights[i]
	}
	return score
}

// LevelOneTrigger 第一级触发机制（单指标触发）
func (t *TriggerMechanism) LevelOneTrigger(metrics []float64) bool {
	for i, metric := range metrics {
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
func (t *TriggerMechanism) EvaluateTrigger(metrics []float64) bool {
	// 检查第一级触发
	levelOneTriggered := t.LevelOneTrigger(metrics)

	// 检查第二级触发
	levelTwoTriggered := t.LevelTwoTrigger(metrics)

	// 如果任何一级触发，则更新历史并检查滞后区间
	if levelOneTriggered || levelTwoTriggered {
		t.TriggerHistory = append(t.TriggerHistory[1:], true) // 记录当前触发
		if t.CheckLagPeriod() {
			fmt.Println("Lag Period Trigger: Executing rebalancing task due to consecutive triggers.")
			return true
		}
	} else {
		// 如果没有触发，更新历史状态
		t.TriggerHistory = append(t.TriggerHistory[1:], false)
	}

	// 返回是否触发
	return false
}

// MonitorAndTrigger 模拟主程序，调用触发机制进行判断
func (a *APH) MonitorAndTrigger(metrics []float64) {
	weights, err := a.CalculateWeights()
	if err != nil {
		fmt.Printf("Error calculating weights: %v\n", err)
	}

	// 初始化触发机制，滞后区间设置为3次连续触发
	triggerMechanism := InitTriggerMechanism(weights, []float64{0.8, 0.75, 0.7, 0.65, 0.9}, 0.8, 3)

	// 每次调用时监测指标
	if triggerMechanism.EvaluateTrigger(metrics) {
		fmt.Println("Trigger condition met! Execute rebalancing algorithm.")
	} else {
		fmt.Println("No trigger, continue monitoring.")
	}
}

func main() {
	// 模拟计算权重
	aph := &APH{}
	aph.Init()

	// 模拟传入的指标数据
	metrics := []float64{0.7, 0.85, 0.6, 0.9, 0.75} // 示例指标值：电量、CPU利用率、内存利用率、网络延迟、偏离距离

	// 监控并根据触发机制决策
	aph.MonitorAndTrigger(metrics)
}
