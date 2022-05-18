package main

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var HostAddress string = "0.0.0.0"
var HostPort string = "8888"

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

// func executeCommand(args []string) {
// 	command := exec.Command(args[0], args[1:]...)
// 	for true {
// 		if err := command.Start(); err != nil {
// 			log.Panic(err)
// 			break
// 		}

// 		if err := command.Wait(); err != nil {

// 		}
// 	}
// }

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
				log.Printf("PUT: %s %s\n", string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				log.Printf("DELETE %s %s\n", string(event.Kv.Key), string(event.Kv.Value))
			}
		}
	}
}
