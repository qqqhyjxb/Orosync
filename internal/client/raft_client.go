package client

import (
	"Orosync/internal/rpc/pb/raft"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var GlobalRaftClient *RaftClient

type RaftClient struct {
	conns sync.Map // 缓存连接: key=address, value=*RaftServiceClient
}

func init() {
	GlobalRaftClient = &RaftClient{
		conns: sync.Map{},
	}
}

type RaftServiceClient struct {
	Conn   *grpc.ClientConn
	Client raft.RaftServiceClient
}

// StartClient 改进后的 StartClient（带连接复用）
func (r *RaftClient) StartClient(target string) (*RaftServiceClient, error) {
	// 1. 尝试从缓存获取已有连接
	if cached, ok := r.conns.Load(target); ok {
		return cached.(*RaftServiceClient), nil
	}

	// 2. 无缓存时创建新连接
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc connection failed: %w", err)
	}

	// 3. 构造客户端对象
	newClient := &RaftServiceClient{
		Client: raft.NewRaftServiceClient(conn),
		Conn:   conn,
	}

	// 4. 原子性存储（防止并发重复创建）
	actual, loaded := r.conns.LoadOrStore(target, newClient)
	if loaded {
		// 如果其他协程已经创建，关闭当前冗余连接
		conn.Close()
		return actual.(*RaftServiceClient), nil
	}

	return newClient, nil
}

// CloseClient 可选：清理指定地址的连接
func (r *RaftClient) CloseClient(target string) {
	if client, ok := r.conns.LoadAndDelete(target); ok {
		client.(*RaftServiceClient).Conn.Close()
	}
}
