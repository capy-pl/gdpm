package main

import (
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

var HostAddress string = "0.0.0.0"
var HostPort string = "8888"

func register() {
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
}

// func sendHeartBeat() {

// }

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
	register()
}
