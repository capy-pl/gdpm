package slave

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gdpm/service"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Slave struct {
	Id            string
	LastHeartBeat time.Time
	ServicesMap   map[string]*service.Service
}

type SlavePool struct {
	lock              sync.Mutex
	Slaves            []*Slave
	ServiceQueue      []*service.Service
	ServiceToSlaveMap map[string]*Slave
	counter           int
	EtcdClient        *clientv3.Client
	Context           context.Context
}

func NewSlavePool(context context.Context, cli *clientv3.Client) *SlavePool {
	return &SlavePool{
		lock:              sync.Mutex{},
		Slaves:            make([]*Slave, 0),
		EtcdClient:        cli,
		Context:           context,
		ServiceToSlaveMap: make(map[string]*Slave),
		ServiceQueue:      make([]*service.Service, 0),
	}
}

func (pool *SlavePool) ScheduleService(service *service.Service) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if len(pool.Slaves) == 0 {
		log.Printf("no availabe slave. put the service %s into service queue\n", service.Id)
		pool.ServiceQueue = append(pool.ServiceQueue, service)
		return
	}

	// use round-robin method to select a slave from the slave pool
	curr := pool.Slaves[pool.counter]

	ctx, cancel := context.WithCancel(pool.Context)
	defer cancel()

	kvs := service.KVs()
	log.Printf("schedule service %s to slave %s\n", service.Id, curr.Id)
	for _, kv := range kvs {
		k, v := kv[0], kv[1]
		_, err := pool.EtcdClient.Put(ctx, curr.GetKeyName(service, k), v)
		if err != nil {
			log.Printf("[PUT] Failed %s Key: %s", service.Id, curr.GetKeyName(service, k))
		} else {
			log.Printf("[PUT] %s Key: %s Value: %s", service.Id, curr.GetKeyName(service, k), v)
		}
	}

	curr.ServicesMap[service.Id] = service
	pool.ServiceToSlaveMap[service.Id] = curr

	// move the counter to the next slave
	pool.counter++

	// if a round ends, shuffle the slave's order and reset the counter to zero
	if pool.counter == len(pool.Slaves) {
		rand.Shuffle(len(pool.Slaves), func(i, j int) {
			pool.Slaves[i], pool.Slaves[j] = pool.Slaves[j], pool.Slaves[i]
		})
		pool.counter = 0
	}
}

func (slave *Slave) GetKeyName(sv *service.Service, key string) string {
	return strings.Join([]string{slave.Id, sv.Id, key}, ":")
}

func (pool *SlavePool) AddSlave(id string) {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	slave := &Slave{
		Id:          id,
		ServicesMap: make(map[string]*service.Service),
	}
	pool.Slaves = append(pool.Slaves, slave)
	log.Printf("%s is added to the pool.", id)
	for len(pool.ServiceQueue) > 0 {
		sv := pool.ServiceQueue[0]
		pool.ServiceQueue = pool.ServiceQueue[1:]
		go pool.ScheduleService(sv)
	}
}

func (pool *SlavePool) RemoveSlave(slave *Slave) {
	pool.lock.Lock()
	defer pool.lock.Unlock()
}

func (pool *SlavePool) UpdateService(serviceId string, instanceNum int) error {
	if slave, exist := pool.ServiceToSlaveMap[serviceId]; exist {
		service, exist := slave.ServicesMap[serviceId]
		if !exist {
			return errors.New("service not found")
		}
		log.Printf("[update] service id %s, instance num %v\n", serviceId, instanceNum)
		_, err := pool.EtcdClient.Put(pool.Context, slave.GetKeyName(service, "InstanceNum"), fmt.Sprintf("%v", instanceNum))
		if err != nil {
			log.Printf("[update] failed, service id %s\n", serviceId)
		} else {
			log.Printf("[update] successed, service id %s\n", serviceId)
		}
	} else {
		return errors.New("service not found")
	}
	return nil
}

func (pool *SlavePool) DeleteService(serviceId string) error {
	slave, exist := pool.ServiceToSlaveMap[serviceId]
	if !exist {
		return errors.New("service not exist")
	}
	service := slave.ServicesMap[serviceId]
	keyprefix := strings.Join([]string{slave.Id, service.Id}, ":")
	_, err := pool.EtcdClient.Delete(pool.Context, keyprefix, clientv3.WithPrefix())
	if err != nil {
		log.Printf("[delete] delete failed, service id = %s\n", service.Id)
		return err
	}

	delete(slave.ServicesMap, service.Id)
	delete(pool.ServiceToSlaveMap, service.Id)

	log.Printf("[delete] delete success, service id = %s\n", service.Id)
	return nil
}
