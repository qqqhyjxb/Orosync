package client

import (
	"Orosync/internal/rpc/pb/raft"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var GlobalRaftClient RaftClient

type RaftClient struct {
}

type RaftServiceClient struct {
	Conn   *grpc.ClientConn
	Client raft.RaftServiceClient
}

func (r *RaftClient) StartClient(target string) *RaftServiceClient {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := raft.NewRaftServiceClient(conn)

	return &RaftServiceClient{
		Client: c,
		Conn:   conn,
	}
}
