package server

import (
	raft2 "Orosync/internal/rpc/pb/raft"
	"Orosync/internal/rpc/pb/simulation"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func StartServer() {
	// 创建gRPC服务器
	server := grpc.NewServer()

	// 注册Service服务
	raft2.RegisterRaftServiceServer(server, &RaftService{})
	simulation.RegisterSimulationServiceServer(server, &SimulationService{})

	// 监听端口
	listen, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 启动服务器
	fmt.Println("Server started at :50052")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
