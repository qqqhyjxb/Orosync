package load_balance

import (
	"Orosync/internal/client"
	"Orosync/internal/model"
	"Orosync/internal/raft"
	"Orosync/internal/rpc/pb/simulation"
	"context"
	"fmt"
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
	for _, t := range raft.GlobalNode.UAV.Tasks {
		sortedTasks = append(sortedTasks, t)
	}

	// 按 CpuUsage 从大到小排序
	sort.Slice(sortedTasks, func(i, j int) bool {
		return sortedTasks[i].CpuUsage > sortedTasks[j].CpuUsage
	})

	// 初始化当前 CPU 使用率和卸载列表
	currentUsage := raft.GlobalNode.UAV.CPU.UsageRate
	var offloadedTasks []model.TaskInfo

	// 遍历排序后的任务，依次尝试卸载
	for _, task := range sortedTasks {
		if currentUsage <= raft.GlobalNode.CpuUsageThreshold {
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
	for _, t := range raft.GlobalNode.UAV.Tasks {
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
	for _, t := range raft.GlobalNode.UAV.Tasks {
		sortedTasks = append(sortedTasks, t)
	}

	// 按 MemoryUsage 从大到小排序
	sort.Slice(sortedTasks, func(i, j int) bool {
		return sortedTasks[i].RequireMemory > sortedTasks[j].RequireMemory
	})

	// 初始化当前 Memory 和卸载列表
	currentMemory := raft.GlobalNode.UAV.Memory.UsageRate
	var offloadedTasks []model.TaskInfo

	// 遍历排序后的任务，依次尝试卸载
	for _, task := range sortedTasks {
		if currentMemory >= raft.GlobalNode.MemoryThreshold {
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
	for _, t := range raft.GlobalNode.UAV.Tasks {
		if _, ok := offloadedIDs[t.TaskId]; !ok {
			remainingTasks = append(remainingTasks, t)
		}
	}

	return offloadedTasks, remainingTasks
}

// LocalLoadBalanceForCPU CPU局部失衡先局部负载均衡
func (b *BalanceHelper) LocalLoadBalanceForCPU(tasks []model.TaskInfo) {
	var taskAssignmentList []*simulation.TaskAssignment

	failTasks := make([]model.TaskInfo, len(tasks))

	for _, t := range tasks {
		for _, u := range raft.GlobalNode.LocalGroup.UidList {
			if v, ok := raft.GlobalNode.Logs.UavMap[u]; ok {
				if v.CPU.UsageRate+t.CpuUsage < raft.GlobalNode.CpuUsageThreshold {
					taskAssignmentList = append(taskAssignmentList, &simulation.TaskAssignment{
						TaskId: t.TaskId,
						UavUid: v.Uid,
					})
				}
			} else {
				fmt.Printf("uav:%s is not exist", u)
			}
		}
		failTasks = append(failTasks, t)
	}

	// 给仿真平台发送再分配结果
	req := &simulation.TaskAssignmentRequest{
		Results: taskAssignmentList,
	}

	clientObj, err := client.GlobalSimulationClient.StartClient(
		raft.GlobalNode.SimulationAddress,
	)
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
	}

	resp, err := clientObj.Client.TaskAssignment(context.Background(), req)
	if err != nil || !resp.Success {
		fmt.Printf("LocalLoadBalanceForCPU failed: %v\n", err)
	}

	// 给leader发送需要全局再分配的任务
	// TODO: 给leader发送需要全局再分配的任务
}

// LocalLoadBalanceForMemory Memory局部失衡先局部负载均衡
func (b *BalanceHelper) LocalLoadBalanceForMemory(tasks []model.TaskInfo) {
	var taskAssignmentList []*simulation.TaskAssignment

	failTasks := make([]model.TaskInfo, len(tasks))

	for _, t := range tasks {
		for _, u := range raft.GlobalNode.LocalGroup.UidList {
			if v, ok := raft.GlobalNode.Logs.UavMap[u]; ok {
				if v.Memory.SurplusCapacity-t.RequireMemory > raft.GlobalNode.MemoryThreshold {
					taskAssignmentList = append(taskAssignmentList, &simulation.TaskAssignment{
						TaskId: t.TaskId,
						UavUid: v.Uid,
					})
				}
			} else {
				fmt.Printf("uav:%s is not exist", u)
			}
		}
		failTasks = append(failTasks, t)
	}

	// 给仿真平台发送再分配结果
	req := &simulation.TaskAssignmentRequest{
		Results: taskAssignmentList,
	}

	clientObj, err := client.GlobalSimulationClient.StartClient(
		raft.GlobalNode.SimulationAddress,
	)
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
	}

	resp, err := clientObj.Client.TaskAssignment(context.Background(), req)
	if err != nil || !resp.Success {
		fmt.Printf("LocalLoadBalanceForCPU failed: %v\n", err)
	}

	// 给leader发送需要全局再分配的任务
	// TODO: 给leader发送需要全局再分配的任务
}
