package client

import (
	"Orosync/internal/rpc/pb/simulation"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var GlobalSimulationClient *SimulationClient

type SimulationClient struct {
	conns sync.Map // 缓存连接: key=address, value=*SimulationServiceClient
}

func init() {
	GlobalSimulationClient = &SimulationClient{
		conns: sync.Map{},
	}
}

type SimulationServiceClient struct {
	Conn   *grpc.ClientConn
	Client simulation.SimulationServiceClient
}

// StartClient 带连接复用的客户端启动方法
func (s *SimulationClient) StartClient(target string) (*SimulationServiceClient, error) {
	// 1. 尝试从缓存获取已有连接
	if cached, ok := s.conns.Load(target); ok {
		return cached.(*SimulationServiceClient), nil
	}

	// 2. 创建新连接
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("simulation grpc connection failed: %w", err)
	}

	// 3. 创建客户端对象
	newClient := &SimulationServiceClient{
		Client: simulation.NewSimulationServiceClient(conn),
		Conn:   conn,
	}

	// 4. 原子性存储防止并发重复创建
	actual, loaded := s.conns.LoadOrStore(target, newClient)
	if loaded {
		// 如果其他协程已经创建，关闭当前冗余连接
		conn.Close()
		return actual.(*SimulationServiceClient), nil
	}

	return newClient, nil
}

// CloseClient 关闭指定地址的连接
func (s *SimulationClient) CloseClient(target string) {
	if client, ok := s.conns.LoadAndDelete(target); ok {
		client.(*SimulationServiceClient).Conn.Close()
	}
}

// CloseAll 关闭所有连接（可选扩展方法）
func (s *SimulationClient) CloseAll() {
	s.conns.Range(func(key, value interface{}) bool {
		value.(*SimulationServiceClient).Conn.Close()
		s.conns.Delete(key)
		return true
	})
}
