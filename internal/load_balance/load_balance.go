package load_balance

import (
	"Orosync/internal/model"
	"Orosync/internal/raft"
	"Orosync/internal/rpc/pb/simulation"
	"fmt"
	"github.com/jinzhu/copier"
	"sort"
)

type BalanceHelper struct{}

var GlobalBalance *BalanceHelper

func InitGlobalBalance() {
	GlobalBalance = &BalanceHelper{}
}

// OffloadCpuTask 卸载cpu型任务
func (b *BalanceHelper) OffloadCpuTask() {
	var newTasks []model.TaskInfo

	for _, t := range raft.GlobalNode.UAV.Tasks {
		newTasks = append(newTasks, t)
	}

	// 使用 sort.Slice 直接排序
	sort.Slice(newTasks, func(i, j int) bool {
		return newTasks[i].CpuUsage > newTasks[j].CpuUsage
	})

	uav := &model.UAV{}

	err := copier.Copy(uav, raft.GlobalNode.UAV)
	if err != nil {
		return
	}

	for _, t := range newTasks {
		uav.CPU.UsageRate = uav.CPU.UsageRate - t.CpuUsage
		//TODO
	}
}

// OffloadMemoryTask 卸载memory型任务
func (b *BalanceHelper) OffloadMemoryTask() {
	var newTasks []model.TaskInfo

	for _, t := range raft.GlobalNode.UAV.Tasks {
		newTasks = append(newTasks, t)
	}

	// 使用 sort.Slice 直接排序
	sort.Slice(newTasks, func(i, j int) bool {
		return newTasks[i].CpuUsage > newTasks[j].CpuUsage
	})

	uav := &model.UAV{}

	err := copier.Copy(uav, raft.GlobalNode.UAV)
	if err != nil {
		return
	}

	for _, t := range newTasks {
		uav.CPU.UsageRate = uav.CPU.UsageRate - t.CpuUsage
		//TODO
	}
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

	// 发送最新的局部分配结果
	req := simulation.TaskAssignmentRequest{
		Results: taskAssignmentList,
	}

	// 获取客户端连接（自动复用）
}

// LocalLoadBalanceForMemory Memory局部失衡先局部负载均衡
func (b *BalanceHelper) LocalLoadBalanceForMemory(tasks []model.TaskInfo) {
	var taskAssignmentList []simulation.TaskAssignment

	failTasks := make([]model.TaskInfo, len(tasks))

	for _, t := range tasks {
		for _, u := range raft.GlobalNode.LocalGroup.UidList {
			if v, ok := raft.GlobalNode.Logs.UavMap[u]; ok {
				if v.Memory.SurplusCapacity-t.RequireMemory > raft.GlobalNode.MemoryThreshold {
					taskAssignmentList = append(taskAssignmentList, simulation.TaskAssignment{
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
}
