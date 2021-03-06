package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdpm/service"
	"github.com/google/uuid"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var HostAddress string = "0.0.0.0"
var HostPort string = "8888"

var ServiceWaitingQueue []*service.Service
var ServiceMap map[string]*service.Service = make(map[string]*service.Service)

func register() string {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", strings.Join([]string{HostAddress, HostPort}, ":"))
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("Failed to connect to %s.\n", tcpAddr)
	}
	defer conn.Close()
	readbuffer := make([]byte, 16)
	_, readerr := conn.Read(readbuffer)
	if readerr != nil {
		log.Fatalln("Cannot get uuid.")
	}
	id, err := uuid.FromBytes(readbuffer)
	if err != nil {
		log.Fatalln("Cannot parse uuid.")
	}
	log.Println(id.String())
	return id.String()
}

func sendHeartBeat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}

func updateService(sv *service.Service) {
	// only process the change of instance number for now
	currInstanceNum := len(sv.Cmds.Queue)
	newInstanceNum := sv.InstanceNum
	log.Printf("%s instance number: %v -> %v\n", sv.Id, currInstanceNum, newInstanceNum)
	if newInstanceNum > currInstanceNum {
		for i := 0; i < newInstanceNum-currInstanceNum; i++ {
			go spawnCommand(sv)
		}
	} else if newInstanceNum < currInstanceNum {
		sv.Cmds.Lock.Lock()
		defer sv.Cmds.Lock.Unlock()
		for i := 0; i < currInstanceNum-newInstanceNum; i++ {
			if sv.Cmds.Queue[i].ProcessState == nil {
				sv.Cmds.Queue[i].Process.Kill()
			}
		}
		sv.Cmds.Queue = sv.Cmds.Queue[currInstanceNum-newInstanceNum:]
	}
}

func deleteService(serviceId string) error {
	if sv, exist := ServiceMap[serviceId]; exist {
		sv.Cmds.Lock.Lock()
		defer sv.Cmds.Lock.Unlock()
		defer delete(ServiceMap, serviceId)

		log.Printf("[delete] service %s exist, instance number = %v\n", sv.Id, len(sv.Cmds.Queue))
		for _, cmd := range sv.Cmds.Queue {
			cmd.Process.Kill()
		}
		log.Printf("[delete] service %s, terminate %v process(es)\n", sv.Id, len(sv.Cmds.Queue))
	} else {
		return fmt.Errorf("service %s not exist", serviceId)
	}
	return nil
}

func spawnCommand(sv *service.Service) {
	sv.Cmds.Lock.Lock()
	command := exec.Command(sv.Command[0], sv.Command[1:]...)
	sv.Cmds.Queue = append(sv.Cmds.Queue, command)
	sv.Cmds.Lock.Unlock()

	err := command.Start()
	if err != nil {
		sv.State = service.Error
	}
	err = command.Wait()
	if err != nil {
		log.Printf("process exit with error %v", err)
	}
}

func main() {
	id := register()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"0.0.0.0:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher := clientv3.NewWatcher(cli)
	watcherChan := watcher.Watch(ctx, id, clientv3.WithPrefix())

	for response := range watcherChan {
		for _, event := range response.Events {
			switch event.Type {
			case mvccpb.PUT:
				if event.PrevKv != nil {
					log.Printf("PUT Prev: %s %s\n", string(event.PrevKv.Key), string(event.PrevKv.Value))
				}
				log.Printf("PUT: %s %s\n", string(event.Kv.Key), string(event.Kv.Value))
				kv := service.ParseServiceKV(string(event.Kv.Key), string(event.Kv.Value))
				if sv, exist := ServiceMap[kv.ServiceId]; !exist && event.PrevKv == nil {
					newsv := &service.Service{
						Id:      kv.ServiceId,
						State:   service.Waiting,
						Command: []string{},
						Cmds: service.ServiceCommandQueue{
							Lock:  sync.Mutex{},
							Queue: make([]*exec.Cmd, 0),
						},
					}
					ServiceMap[newsv.Id] = newsv
					ServiceWaitingQueue = append(ServiceWaitingQueue, newsv)
					log.Printf("New Service Created: %s\n", kv.ServiceId)
				} else {
					switch kv.Key {
					case service.Id:
						sv.Id = kv.Value
					case service.Command:
						sv.Command = strings.Split(kv.Value, " ")
					case service.InstanceNum:
						sv.InstanceNum, _ = strconv.Atoi(kv.Value)
					}
				}
				sv := ServiceMap[kv.ServiceId]
				if sv.State == service.Running {
					updateService(sv)
				}
				if sv.State == service.Waiting && sv.Command != nil && len(sv.Command) > 0 && sv.InstanceNum > 0 {
					sv.State = service.Ready
					log.Printf("Service Ready To Be Scheduled: %s\n", sv.Id)
					sv.State = service.Running
					updateService(sv)
				}
			case mvccpb.DELETE:
				log.Printf("DELETE %s %s\n", string(event.Kv.Key), string(event.Kv.Value))
				kv := service.ParseServiceKV(string(event.Kv.Key), string(event.Kv.Value))
				if _, exist := ServiceMap[kv.ServiceId]; exist {
					deleteService(kv.ServiceId)
				}
			}
		}
	}
}
