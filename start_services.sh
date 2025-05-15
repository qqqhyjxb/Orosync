#!/bin/bash

# 使用普通的变量定义 UAV 配置
# shellcheck disable=SC2034
UAV001="192.168.1.87:40001"
UAV002="192.168.1.87:40002"
UAV003="192.168.1.87:40003"
UAV004="192.168.1.87:40004"
UAV005="192.168.1.87:40005"
UAV006="192.168.1.87:40006"
UAV007="192.168.1.87:40007"
UAV008="192.168.1.87:40008"
UAV009="192.168.1.87:40009"
UAV010="192.168.1.87:40010"
UAV011="192.168.1.87:40011"
UAV012="192.168.1.87:40012"
UAV013="192.168.1.87:40013"
UAV014="192.168.1.87:40014"
UAV015="192.168.1.87:40015"
UAV016="192.168.1.87:40016"
UAV017="192.168.1.87:40017"
UAV018="192.168.1.87:40018"
UAV019="192.168.1.87:40019"
UAV020="192.168.1.87:40020"


# UAV 配置数组
UAVs=(
  UAV001
  UAV002
  UAV003
  UAV004
  UAV005
  UAV006
  UAV007
  UAV008
  UAV009
  UAV010
  UAV011
  UAV012
  UAV013
  UAV014
  UAV015
  UAV016
  UAV017
  UAV018
  UAV019
  UAV020
)

# 遍历 UAV 配置数组，启动每一个 UAV 实例
for uav in "${UAVs[@]}"
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
