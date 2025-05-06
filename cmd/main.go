package main

import (
	"Orosync/internal/monitor"
	"Orosync/internal/raft"
	"Orosync/internal/server"
)

func main() {
	raft.InitRaft()

	server.StartServer()

	monitor.InitGlobalAPH()
	monitor.GlobalAPH.Init()
	monitor.GlobalAPH.Test()
}
