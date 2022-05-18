package service

import (
	"strings"

	"github.com/google/uuid"
)

type ServiceState uint8
type ServiceKey string

const (
	Unscheduled ServiceState = 0
	Scheduled   ServiceState = 1
	Running     ServiceState = 2
	Exit        ServiceState = 3
)

const (
	Id          = "Id"
	Command     = "Command"
	InstanceNum = "InstanceNum"
	State       = "State"
)

type Service struct {
	Id          string
	Command     []string
	InstanceNum int
	State       ServiceState
}

func NewService(command string, instanceNum int) *Service {
	commands := strings.Split(command, " ")
	id := uuid.NewString()
	return &Service{
		Id:          id,
		Command:     commands,
		InstanceNum: 1,
		State:       Unscheduled,
	}
}

func (service *Service) KVs() [][]string {
	rets := make([][]string, 0)
	rets = append(rets, []string{Id, service.Id})
	rets = append(rets, []string{Command, strings.Join(service.Command, " ")})
	rets = append(rets, []string{InstanceNum, string(service.InstanceNum)})
	rets = append(rets, []string{State, string(service.State)})
	return rets
}
