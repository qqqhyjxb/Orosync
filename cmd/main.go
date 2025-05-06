package main

import (
	"Orosync/internal/monitor"
	"Orosync/internal/raft"
)

func main() {
	raft.InitRaft()

	//server.StartServer()

	monitor.InitGlobalAPH()
	monitor.GlobalAPH.Init()
	monitor.GlobalAPH.Test()
}
