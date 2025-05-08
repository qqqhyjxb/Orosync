#!/bin/bash

# 使用普通的变量定义 UAV 配置
# shellcheck disable=SC2034
UAV001="192.168.1.87:40000"
UAV002="192.168.1.87:40001"
UAV003="192.168.1.87:40002"
UAV004="192.168.1.87:40003"
UAV005="192.168.1.87:40004"
UAV006="192.168.1.87:40005"
UAV007="192.168.1.87:40006"
UAV008="192.168.1.87:40007"
UAV009="192.168.1.87:40008"
UAV010="192.168.1.87:40009"

# UAV005 UAV006 UAV007 UAV008 UAV009 UAV010

# 遍历 UAV 配置，启动每一个 UAV 实例
for uav in UAV001 UAV002 UAV003 UAV004 UAV005
do
    # 获取对应 UAV 的 IP 地址和端口
    address_and_port="${!uav}"
    address=$(echo $address_and_port | cut -d':' -f1)  # 获取 IP 地址
    port=$(echo $address_and_port | cut -d':' -f2)     # 获取端口号

    # 启动 Go 程序，并传递相应的参数（uid, ip, port）
    go run ./cmd/main.go -uid $uav -ip $address -port $port &

    # 等待一小段时间，确保程序能够启动
    sleep 1
done

# 等待所有后台进程完成
wait
echo "所有 UAV 服务已闭关。"
