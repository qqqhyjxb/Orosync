package raft

import (
	"Orosync/internal/client"
	"Orosync/internal/model"
	r "Orosync/internal/rpc/pb/raft"
	"Orosync/internal/rpc/pb/simulation"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"math"
	"sort"
)

type BalanceHelper struct{}

var GlobalBalance *BalanceHelper

func init() {
	GlobalBalance = &BalanceHelper{}
}

// OffloadCpuTask 卸载 CPU 型任务，返回 (被卸载的任务列表, 保留的任务列表)
func (b *BalanceHelper) OffloadCpuTask() ([]model.TaskInfo, []model.TaskInfo) {
	// 复制原始任务列表（避免直接操作原始数据）
	var sortedTasks []model.TaskInfo
	for _, t := range GlobalNode.UAV.Tasks {
		sortedTasks = append(sortedTasks, t)
	}

	// 按 CpuUsage 从大到小排序
	sort.Slice(sortedTasks, func(i, j int) bool {
		return sortedTasks[i].CpuUsage > sortedTasks[j].CpuUsage
	})

	// 初始化当前 CPU 使用率和卸载列表
	currentUsage := GlobalNode.UAV.CPU.UsageRate
	var offloadedTasks []model.TaskInfo

	// 遍历排序后的任务，依次尝试卸载
	for _, task := range sortedTasks {
		if currentUsage <= GlobalNode.CpuUsageThreshold {
			break // 已达阈值，停止卸载
		}
		currentUsage -= task.CpuUsage
		offloadedTasks = append(offloadedTasks, task)
	}

	// 构建保留任务列表（排除已卸载的任务）
	offloadedIDs := make(map[int32]struct{})
	for _, t := range offloadedTasks {
		offloadedIDs[t.TaskId] = struct{}{}
	}

	var remainingTasks []model.TaskInfo
	for _, t := range GlobalNode.UAV.Tasks {
		if _, ok := offloadedIDs[t.TaskId]; !ok {
			remainingTasks = append(remainingTasks, t)
		}
	}

	return offloadedTasks, remainingTasks
}

// OffloadMemoryTask 卸载 Memory 型任务，返回 (被卸载的任务列表, 保留的任务列表)
func (b *BalanceHelper) OffloadMemoryTask() ([]model.TaskInfo, []model.TaskInfo) {
	// 复制原始任务列表（避免直接操作原始数据）
	var sortedTasks []model.TaskInfo
	for _, t := range GlobalNode.UAV.Tasks {
		sortedTasks = append(sortedTasks, t)
	}

	// 按 RequireMemory 从大到小排序
	sort.Slice(sortedTasks, func(i, j int) bool {
		return sortedTasks[i].RequireMemory > sortedTasks[j].RequireMemory
	})

	// 初始化当前 Memory 和卸载列表
	currentMemory := GlobalNode.UAV.Memory.UsageRate
	var offloadedTasks []model.TaskInfo

	// 遍历排序后的任务，依次尝试卸载
	for _, task := range sortedTasks {
		if currentMemory >= GlobalNode.MemoryThreshold {
			break // 已达阈值，停止卸载
		}
		currentMemory += task.RequireMemory
		offloadedTasks = append(offloadedTasks, task)
	}

	// 构建保留任务列表（排除已卸载的任务）
	offloadedIDs := make(map[int32]struct{})
	for _, t := range offloadedTasks {
		offloadedIDs[t.TaskId] = struct{}{}
	}

	var remainingTasks []model.TaskInfo
	for _, t := range GlobalNode.UAV.Tasks {
		if _, ok := offloadedIDs[t.TaskId]; !ok {
			remainingTasks = append(remainingTasks, t)
		}
	}

	return offloadedTasks, remainingTasks
}

// LocalLoadBalanceForCPU CPU局部失衡先局部负载均衡
func (b *BalanceHelper) LocalLoadBalanceForCPU(tasks []model.TaskInfo) {
	var taskAssignmentList []*simulation.TaskAssignment

	failTasks := make([]*r.Task, len(tasks))

	for _, t := range tasks {
		for _, u := range GlobalNode.LocalGroup.UidList {
			if v, ok := GlobalNode.Logs.UavMap[u]; ok {
				if v.CPU.UsageRate+t.CpuUsage < GlobalNode.CpuUsageThreshold {
					taskAssignmentList = append(taskAssignmentList, &simulation.TaskAssignment{
						TaskId: t.TaskId,
						UavUid: v.Uid,
					})
				}
			} else {
				fmt.Printf("uav:%s is not exist", u)
			}
		}

		task := &r.Task{}

		err := copier.Copy(task, t)
		if err != nil {
			return
		}

		failTasks = append(failTasks, task)
	}

	// 给仿真平台发送再分配结果
	//req := &simulation.TaskAssignmentRequest{
	//	Results: taskAssignmentList,
	//}
	//
	//simulationClientObj, err := client.GlobalSimulationClient.StartClient(
	//	GlobalNode.SimulationAddress,
	//)
	//if err != nil {
	//	fmt.Printf("连接失败: %v\n", err)
	//}
	//
	//resp, err := simulationClientObj.Client.TaskAssignment(context.Background(), req)
	//if err != nil || !resp.Success {
	//	fmt.Printf("LocalLoadBalanceForCPU failed: %v\n", err)
	//}

	// 给leader发送需要全局再分配的任务
	raftClientObj, err := client.GlobalRaftClient.StartClient(
		GlobalNode.Logs.UavMap[GlobalNode.LeaderUid].Address,
	)
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
	}

	request := &r.GlobalLoadBalanceReq{
		TaskList: failTasks,
	}

	response, err := raftClientObj.Client.GlobalLoadBalance(context.Background(), request)
	if err != nil || response.Code != SuccessCode {
		fmt.Printf("LocalLoadBalanceForCPU failed: %v\n", err)
	}
}

// LocalLoadBalanceForMemory Memory局部失衡先局部负载均衡
func (b *BalanceHelper) LocalLoadBalanceForMemory(tasks []model.TaskInfo) {
	var taskAssignmentList []*simulation.TaskAssignment

	failTasks := make([]*r.Task, len(tasks))

	for _, t := range tasks {
		for _, u := range GlobalNode.LocalGroup.UidList {
			if v, ok := GlobalNode.Logs.UavMap[u]; ok {
				if v.Memory.SurplusCapacity-t.RequireMemory > GlobalNode.MemoryThreshold {
					taskAssignmentList = append(taskAssignmentList, &simulation.TaskAssignment{
						TaskId: t.TaskId,
						UavUid: v.Uid,
					})
				}
			} else {
				fmt.Printf("uav:%s is not exist", u)
			}
		}
		task := &r.Task{}

		err := copier.Copy(task, t)
		if err != nil {
			return
		}

		failTasks = append(failTasks, task)
	}

	// 给仿真平台发送再分配结果
	//req := &simulation.TaskAssignmentRequest{
	//	Results: taskAssignmentList,
	//}
	//
	//clientObj, err := client.GlobalSimulationClient.StartClient(
	//	GlobalNode.SimulationAddress,
	//)
	//if err != nil {
	//	fmt.Printf("连接失败: %v\n", err)
	//}
	//
	//resp, err := clientObj.Client.TaskAssignment(context.Background(), req)
	//if err != nil || !resp.Success {
	//	fmt.Printf("LocalLoadBalanceForCPU failed: %v\n", err)
	//}

	// 给leader发送需要全局再分配的任务
	raftClientObj, err := client.GlobalRaftClient.StartClient(
		GlobalNode.Logs.UavMap[GlobalNode.LeaderUid].Address,
	)
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
	}

	request := &r.GlobalLoadBalanceReq{
		TaskList: failTasks,
	}

	response, err := raftClientObj.Client.GlobalLoadBalance(context.Background(), request)
	if err != nil || response.Code != SuccessCode {
		fmt.Printf("LocalLoadBalanceForCPU failed: %v\n", err)
	}
}

// distance 计算两个位置之间的欧氏距离
func distance(a, b *model.PositionInfo) float32 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// calculateProfit 计算任务收益值
func calculateProfit(uav *model.UAV, task *model.TaskInfo, alpha, beta, gamma float32) float32 {
	// 资源匹配计算
	cpuMargin := uav.CPU.SurplusCapacity - task.RequiredCPU
	memMargin := uav.Memory.SurplusCapacity - task.RequireMemory

	// 电量惩罚计算（电量越低惩罚越大）
	batteryPenalty := GlobalNode.BatteryThreshold - uav.Battery.Capacity
	if batteryPenalty < 0 {
		batteryPenalty = 0
	}

	// 距离惩罚计算
	distancePenalty := distance(&uav.Position, &task.Target)

	// 综合收益计算
	return cpuMargin + memMargin + alpha*float32(task.Priority) - beta*batteryPenalty - gamma*distancePenalty
}

// selectBestFromCandidates 处理平局选择最优无人机
func selectBestFromCandidates(candidates []*model.UAV, task *model.TaskInfo) *model.UAV {
	// 第一级筛选：剩余电量最高
	maxBattery := float32(-1)
	var batteryCandidates []*model.UAV
	for _, u := range candidates {
		if u.Battery.Capacity > maxBattery {
			maxBattery = u.Battery.Capacity
			batteryCandidates = []*model.UAV{u}
		} else if u.Battery.Capacity == maxBattery {
			batteryCandidates = append(batteryCandidates, u)
		}
	}

	// 如果唯一候选直接返回
	if len(batteryCandidates) == 1 {
		return batteryCandidates[0]
	}

	// 第二级筛选：距离最近
	minDistance := float32(math.MaxFloat32)
	var selected *model.UAV
	for _, u := range batteryCandidates {
		d := distance(&u.Position, &task.Target)
		if d < minDistance {
			minDistance = d
			selected = u
		}
	}
	return selected
}

func (b *BalanceHelper) GlobalLoadBalance(tasks []model.TaskInfo) error {
	// 参数校验
	if len(tasks) == 0 {
		return errors.New("tasks is empty")
	}

	// 初始化算法参数
	const (
		maxRetries         = 3            // 最大重试次数
		epsilon            = float32(0.1) // 价格增量
		alpha      float32 = 0.5          // 优先级权重
		beta       float32 = 0.3          // 电量惩罚权重
		gamma      float32 = 0.2          // 距离惩罚权重
	)

	var assignmentList []*simulation.TaskAssignment

	// 遍历所有待分配任务
	for taskIdx := range tasks {
		task := &tasks[taskIdx]
		retryCount := 0
		assigned := false

		// 价格调整循环
		for retryCount < maxRetries && !assigned {
			var bestUAV *model.UAV
			maxProfit := float32(-math.MaxFloat32)
			var candidates []*model.UAV

			// 遍历集群所有无人机
			for _, uid := range GlobalNode.LocalGroup.UidList {
				uav, exists := GlobalNode.Logs.UavMap[uid]
				if !exists {
					continue
				}

				// 基础资源检查
				if uav.CPU.SurplusCapacity < task.RequiredCPU ||
					uav.Memory.SurplusCapacity < task.RequireMemory {
					continue
				}

				// 计算收益值
				currentProfit := calculateProfit(uav, task, alpha, beta, gamma)

				// 考虑价格因素（示例：价格越高优先级越高）
				currentProfit += task.Cost["price"] // 假设价格存储在cost map中

				// 更新最佳候选
				if currentProfit > maxProfit {
					maxProfit = currentProfit
					bestUAV = uav
					candidates = []*model.UAV{uav}
				} else if currentProfit == maxProfit {
					candidates = append(candidates, uav)
				}
			}

			// 处理平局情况
			if len(candidates) > 1 {
				bestUAV = selectBestFromCandidates(candidates, task)
			}

			if bestUAV != nil {
				// 执行任务分配
				bestUAV.Tasks = append(bestUAV.Tasks, *task)

				// 更新无人机资源状态
				bestUAV.CPU.SurplusCapacity -= task.RequiredCPU
				bestUAV.Memory.SurplusCapacity -= task.RequireMemory
				bestUAV.CPU.UsageRate += task.CpuUsage
				bestUAV.Memory.UsageRate += task.RequireMemory / bestUAV.Memory.Capacity

				// 记录分配结果
				assignmentList = append(assignmentList, &simulation.TaskAssignment{
					TaskId: task.TaskId,
					UavUid: bestUAV.Uid,
				})

				// 更新任务状态
				task.Status = "assigned"
				assigned = true
			} else {
				// 价格调整策略
				task.Cost["price"] += epsilon
				retryCount++
			}
		}

		if !assigned {
			return fmt.Errorf("task %d failed to assign after %d retries", task.TaskId, maxRetries)
		}
	}

	// 发送分配结果到仿真平台
	if len(assignmentList) > 0 {
		simClient, err := client.GlobalSimulationClient.StartClient(GlobalNode.SimulationAddress)
		if err != nil {
			return fmt.Errorf("simulation connection failed: %v", err)
		}

		req := &simulation.TaskAssignmentRequest{
			Results: assignmentList,
		}

		if resp, err := simClient.Client.TaskAssignment(context.Background(), req); err != nil || !resp.Success {
			return fmt.Errorf("task assignment failed: %v", err)
		}
	}

	// 更新集群日志
	GlobalNode.Logs.Count += int64(len(tasks))
	return nil
}
