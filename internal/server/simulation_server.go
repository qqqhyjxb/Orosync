package server

import (
	"Orosync/internal/raft"
	"Orosync/internal/rpc/pb/simulation"
	"context"
	"fmt"
	"github.com/jinzhu/copier"
)

type SimulationService struct {
	simulation.UnimplementedSimulationServiceServer
}

func (s *SimulationService) DroneStatus(ctx context.Context,
	req *simulation.DroneStatusRequest) (*simulation.DroneStatusResponse, error) {

	resp := &simulation.DroneStatusResponse{}

	err := copier.Copy(raft.GlobalNode.UAV, req)
	if err != nil {
		resp.Code = raft.FailCode
		resp.Msg = err.Error()
		return resp, err
	}

	fmt.Printf("监听到的无人机状态:%v\n", raft.GlobalNode.UAV)

	resp.Code = raft.SuccessCode
	resp.Msg = "success"
	return resp, nil
}

func (s *SimulationService) DroneSwarmChange(ctx context.Context,
	req *simulation.DroneSwarmChangeRequest) (*simulation.DroneSwarmChangeResponse, error) {

	resp := &simulation.DroneSwarmChangeResponse{}

	if raft.GlobalNode.Role != raft.Leader {
		resp.Code = raft.NotLeaderCode
		resp.Msg = "not leader"
		//resp.Address = raft.GlobalNode.Logs[raft.GlobalNode.NewLogIndex].UavMap[raft.GlobalNode.LeaderUid].Address
		//resp.Port = raft.GlobalNode.Logs[raft.GlobalNode.NewLogIndex].UavMap[raft.GlobalNode.LeaderUid].Port
	}

	resp, err := raft.GlobalLeader.DroneSwarmChange(ctx, req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
