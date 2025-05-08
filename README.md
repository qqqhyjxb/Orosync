# Orosync
动态拓扑下无人机集群的负载均衡

## 研究背景
无人机边缘节点的动态变动会导致负载分配的不均衡和任务中断，影响系统稳定性。需要快速响应节点的加入或离开，确保任务分配合理且资源利用最大化。

## 主要开展的工作
1. 分析节点加入或离开对系统负载和任务调度的影响，明确再均衡的触发条件。
2. 考虑节点状态、任务依赖关系和通信延迟等因素，提出适应动态环境的再均衡策略，兼顾实时性和计算复杂度。
3. 通过模拟节点动态变动场景（如频繁加入/离开），验证算法的效率和鲁棒性。

## 总体目标
实现一个能够在无人机节点变动时自动进行负载再均衡的机制，提升系统的鲁棒性和适应性，确保任务的连续性和系统的稳定运行

## 架构设计
系统分为 raft集群、触发器、负载均衡三大模块

### raft集群
通过raft协议维护集群数据的一致性

### 触发器
设计系统再均衡的触发条件

### 负载均衡
再均衡算法的具体实现

## TODO：

## ing～：
- 设计无人机节点的模型
- 触发机制
- 一致性保障-raft协议改造
- grpc

## done：

## 快捷指令
### 生成pb文件
```shell 
protoc --go_out=./internal/rpc/pb   \
--go-grpc_out=./internal/rpc/pb   \
./internal/rpc/proto/raft.proto
```

```shell
protoc --go_out=./internal/rpc/pb   \
--go-grpc_out=./internal/rpc/pb   \
./internal/rpc/proto/simulation.proto
```

```shell
lsof -i :40000 | awk '/LISTEN/ {print $2}' | xargs -r sudo kill -9
```






