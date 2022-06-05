package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdpm/service"
	"github.com/gdpm/slave"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func test(pool *slave.SlavePool) {
	time.Sleep(100 * time.Millisecond)
	sv := service.NewService("python test.py", 1)
	pool.ScheduleService(sv)
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

func handleCreateService(pool *slave.SlavePool) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			log.Println("New service request is received.")
			req.ParseForm()
			command := req.PostFormValue("Command")
			if len(command) == 0 {
				http.Error(res, "Command is required field", http.StatusBadRequest)
				return
			}
			log.Printf("Command is %s\n", command)
			instanceNumStr := req.PostFormValue("InstanceNum")
			if len(instanceNumStr) == 0 {
				instanceNumStr = "1"
			}
			log.Printf("Instance number is %s\n", instanceNumStr)
			instanceNum, err := strconv.Atoi(instanceNumStr)
			if err != nil {
				http.Error(res, "not a valid instanceNum", http.StatusBadRequest)
				return
			}
			sv := service.NewService(command, instanceNum)
			log.Printf("[New] %s is created. Command: %s\n", sv.Id, strings.Join(sv.Command, " "))

			pool.ScheduleService(sv)

			res.WriteHeader(http.StatusOK)
			res.Write([]byte(sv.Id))
		} else {
			http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleUpdateService(pool *slave.SlavePool) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		serviceId := params["serviceId"]
		instanceNum, err := strconv.Atoi(req.FormValue("InstanceNum"))
		if err != nil || instanceNum < 0 {
			http.Error(res, "not a valid instanceNum", http.StatusBadRequest)
			return
		}
		pool.UpdateService(serviceId, instanceNum)
	}
}

func handleDeleteService(pool *slave.SlavePool) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		serviceId := params["serviceId"]
		pool.DeleteService(serviceId)
	}
}

type GetNodesResponse struct {
	Ids        []string
	ServiceNum []int
	Status     []int
}

func handleGetNodes(pool *slave.SlavePool) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		encoder := json.NewEncoder(res)
		jsonResponse := GetNodesResponse{
			Ids:        make([]string, len(pool.Slaves)),
			ServiceNum: make([]int, len(pool.Slaves)),
			Status:     make([]int, len(pool.Slaves)),
		}
		for i := 0; i < len(pool.Slaves); i++ {
			jsonResponse.Ids[i] = pool.Slaves[i].Id
			jsonResponse.ServiceNum[i] = len(pool.Slaves[i].ServicesMap)
			jsonResponse.Status[i] = 1
		}
		encoder.Encode(jsonResponse)
	}
}

func startHttpServer(pool *slave.SlavePool) {
	r := mux.NewRouter()

	// api endpoints for services
	serviceHandle := r.PathPrefix("/service").Subrouter()
	serviceHandle.HandleFunc("/", handleCreateService(pool)).Methods("POST")
	serviceHandle.HandleFunc("/{serviceId}/", handleUpdateService(pool)).Methods("POST")
	serviceHandle.HandleFunc("/{serviceId}/", handleDeleteService(pool)).Methods("DELETE")

	// api endpoints for nodes
	nodeHandle := r.PathPrefix("/node").Subrouter()
	nodeHandle.HandleFunc("/", handleGetNodes(pool))

	http.Handle("/", r)
	server := &http.Server{
		Handler:      r,
		Addr:         strings.Join([]string{ListeningAddress, HTTPListeningPort}, ":"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("http server is listening at %s:%s", ListeningAddress, HTTPListeningPort)
	log.Fatalln(server.ListenAndServe())
	defer server.Close()
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
