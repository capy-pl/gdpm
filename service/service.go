package service

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type ServiceState uint8

const (
	Unscheduled ServiceState = 0
	Waiting     ServiceState = 11
	Ready       ServiceState = 10
	Scheduled   ServiceState = 1
	Running     ServiceState = 2
	Exit        ServiceState = 3
	Error       ServiceState = 4
)

const (
	Id          = "Id"
	Command     = "Command"
	InstanceNum = "InstanceNum"
	State       = "State"
)

type ServiceCommandQueue struct {
	Lock  sync.Mutex
	Queue []*exec.Cmd
}

type Service struct {
	Id          string
	Command     []string
	InstanceNum int
	State       ServiceState
	Cmds        ServiceCommandQueue
}

type ServiceKV struct {
	SlaveId   string
	ServiceId string
	Key       string
	Value     string
}

func NewService(command string, instanceNum int) *Service {
	commands := strings.Split(command, " ")
	id := uuid.NewString()
	return &Service{
		Id:          id,
		Command:     commands,
		InstanceNum: instanceNum,
		State:       Unscheduled,
	}
}

func (service *Service) KVs() [][]string {
	rets := make([][]string, 0)
	rets = append(rets, []string{Id, service.Id})
	rets = append(rets, []string{Command, strings.Join(service.Command, " ")})
	rets = append(rets, []string{InstanceNum, fmt.Sprintf("%d", service.InstanceNum)})
	rets = append(rets, []string{State, string(service.State)})
	return rets
}

func ParseServiceKV(key string, value string) ServiceKV {
	strlist := strings.Split(key, ":")
	return ServiceKV{
		SlaveId:   strlist[0],
		ServiceId: strlist[1],
		Key:       strlist[2],
		Value:     value,
	}
}
