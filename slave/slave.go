package slave

import (
	"context"
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
	Services      []*service.Service
}

type SlavePool struct {
	lock       sync.Mutex
	Slaves     []*Slave
	counter    int
	EtcdClient *clientv3.Client
	Context    context.Context
}

func NewSlavePool(context context.Context, cli *clientv3.Client) *SlavePool {
	return &SlavePool{
		lock:       sync.Mutex{},
		Slaves:     make([]*Slave, 0),
		EtcdClient: cli,
		Context:    context,
	}
}

func (pool *SlavePool) ScheduleService(service *service.Service) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	// use a round-robin algo to select a slave from slave pool
	curr := pool.Slaves[pool.counter]

	ctx, cancel := context.WithCancel(pool.Context)
	defer cancel()

	kvs := service.KVs()
	for _, kv := range kvs {
		k, v := kv[0], kv[1]
		_, err := pool.EtcdClient.Put(ctx, curr.GetKeyName(service, k), v)
		if err != nil {
			log.Printf("[PUT] Failed %s Key: %s", service.Id, curr.GetKeyName(service, k))
		} else {
			log.Printf("[PUT] %s Key: %s Value: %s", service.Id, curr.GetKeyName(service, k), v)
		}
	}

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
		Id:       id,
		Services: make([]*service.Service, 0),
	}
	pool.Slaves = append(pool.Slaves, slave)
	log.Printf("%s is added to the pool.", id)
}

func (pool *SlavePool) RemoveSlave(slave *Slave) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

}
