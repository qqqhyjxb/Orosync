package main

import (
	"Orosync/internal/raft"
	"Orosync/internal/server"
)

func main() {
	raft.InitRaft()

	server.StartServer()
}
