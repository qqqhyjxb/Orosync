package server

import (
	raft2 "Orosync/internal/rpc/pb/raft"
	"Orosync/internal/rpc/pb/simulation"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func StartServer(port int) {
	// 创建gRPC服务器
	server := grpc.NewServer()

	reflection.Register(server)

	// 注册Service服务
	raft2.RegisterRaftServiceServer(server, &RaftService{})
	simulation.RegisterSimulationServiceServer(server, &SimulationService{})

	address := fmt.Sprintf(":%d", port)

	// 监听端口
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 启动服务器
	fmt.Printf("Server started at %v\n", port)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
