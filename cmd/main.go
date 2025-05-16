package main

import (
	"Orosync/internal/config"
	"Orosync/internal/model"
	"Orosync/internal/monitor"
	"Orosync/internal/raft"
	"Orosync/internal/server"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 测试用绝对路径更方便
	cfg, err := config.InitConfig("/Users/liangqiwang/GolandProjects/orosync/Orosync/etc/config.yaml")
	if err != nil {
		fmt.Println("init config err:", err)
	}

	uav := &model.UAV{}

	// 1. 从命令行参数获取无人机初始数据
	var port int
	var uid string
	var address string

	// 定义命令行参数
	flag.IntVar(&port, "port", 3000, "Port to listen on")
	flag.StringVar(&uid, "uid", "", "UID of the UAV")
	flag.StringVar(&address, "ip", "", "ip of the UAV")
	flag.Parse()

	uav.Uid = uid
	uav.Address = fmt.Sprintf("%s:%d", address, port)
	uav.Port = port

	// 2. 验证端口有效性
	if port < 1 || port > 65535 {
		log.Fatalf("Invalid port number: %d", port)
	}

	// 3. 初始化 Raft（添加端口标识用于测试）
	raft.InitRaft(uav, cfg.Logs)

	// 4. 初始化监控系统
	monitor.InitMonitor()

	// 5. 启动 Raft 在一个 goroutine 中
	go func() {
		raft.GlobalNode.Start()
	}()

	// 6. 启动 HTTP 服务器
	go func() {
		server.StartServer(port)
	}()

	// 7. 等待系统退出信号并优雅关闭
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gracefully...")

}
