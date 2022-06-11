package main

import (
	"context"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gdpm/slave"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var workerNum int = 4
var ListeningAddress string = "0.0.0.0"
var ListeningPort string = "8888"
var HTTPListeningPort string = "8989"

func listenForSlave(ch chan *net.TCPConn) {
	address, _ := net.ResolveTCPAddr("tcp", strings.Join([]string{ListeningAddress, ListeningPort}, ":"))
	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Server starts listening for connection at %s.", address)
	defer func() {
		log.Printf("Server stop listening for connection at %s.", address)
		listener.Close()
	}()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err.Error())
			conn.Close()
			continue
		}
		ch <- conn
	}
}

func handleConnection(ch <-chan *net.TCPConn, pool *slave.SlavePool) {
	var conn *net.TCPConn
	for {
		conn = <-ch
		defer conn.Close()

		log.Printf("accept connection from %s\n", conn.RemoteAddr())
		newid := uuid.New()
		log.Printf("assign uuid %s to %s", newid.String(), conn.RemoteAddr())
		newidbin, _ := newid.MarshalBinary()
		_, err := conn.Write(newidbin)
		if err != nil {
			log.Printf("failed to send uuid to %s", conn.RemoteAddr())
			log.Printf("close connection from %s\n", conn.RemoteAddr())
			return
		}
		log.Printf("close connection from %s\n", conn.RemoteAddr())
		pool.AddSlave(newid.String())
	}
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan *net.TCPConn, workerNum)
	defer close(ch)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"0.0.0.0:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	pool := slave.NewSlavePool(ctx, cli)

	wg.Add(1)
	go func() {
		listenForSlave(ch)
		wg.Done()
	}()

	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			handleConnection(ch, pool)
			wg.Done()
		}()
	}
	wg.Add(1)
	go startHttpServer(pool)

	wg.Wait()
}
