package main

import (
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

var workerNum int = 4

var ListeningAddress string = "0.0.0.0"
var ListeningPort string = "8888"

var slavePool []string

func listenForSlave(ch chan *net.TCPConn) {
	address, _ := net.ResolveTCPAddr("tcp", strings.Join([]string{ListeningAddress, ListeningPort}, ":"))
	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	go log.Printf("Server starts listening for connection at %s.", address)
	defer func() {
		go log.Printf("Server stop listening for connection at %s.", address)
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

func handleConnection(ch <-chan *net.TCPConn) {
	var conn *net.TCPConn
	for {
		conn = <-ch
		log.Printf("accept connection from %s\n", conn.RemoteAddr())
		newid := uuid.New()
		log.Printf("assign uuid %s to %s", newid.String(), conn.RemoteAddr())
		newidbin, _ := newid.MarshalBinary()
		_, err := conn.Write(newidbin)
		if err != nil {
			log.Printf("failed to send uuid to %s", conn.RemoteAddr())
			log.Printf("close connection from %s\n", conn.RemoteAddr())
			conn.Close()
			return
		}
		log.Printf("close connection from %s\n", conn.RemoteAddr())
		conn.Close()
		slavePool = append(slavePool, newid.String())
	}
}

func main() {
	ch := make(chan *net.TCPConn, workerNum)
	defer close(ch)
	for i := 0; i < workerNum; i++ {
		go handleConnection(ch)
	}
	listenForSlave(ch)
}
