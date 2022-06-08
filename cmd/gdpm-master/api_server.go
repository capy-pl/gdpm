package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gdpm/service"
	"github.com/gdpm/slave"
	"github.com/gorilla/mux"
)

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
		err = pool.UpdateService(serviceId, instanceNum)
		if err != nil {
			http.Error(res, "error", http.StatusBadRequest)
			return
		}
	}
}

func handleDeleteService(pool *slave.SlavePool) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		serviceId := params["serviceId"]
		err := pool.DeleteService(serviceId)
		if err != nil {
			http.Error(res, "error", http.StatusBadRequest)
			return
		}
	}
}

type GetNodesResponse struct {
	Ids        []string
	ServiceNum []int
	Status     []int
	Times      []string
}

func handleGetNodes(pool *slave.SlavePool) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		encoder := json.NewEncoder(res)
		jsonResponse := GetNodesResponse{
			Ids:        make([]string, len(pool.Slaves)),
			ServiceNum: make([]int, len(pool.Slaves)),
			Status:     make([]int, len(pool.Slaves)),
			Times:      make([]string, len(pool.Slaves)),
		}
		for i := 0; i < len(pool.Slaves); i++ {
			jsonResponse.Ids[i] = pool.Slaves[i].Id
			jsonResponse.ServiceNum[i] = len(pool.Slaves[i].ServicesMap)
			jsonResponse.Times[i] = pool.Slaves[i].StartTime.Format(time.RFC822)
			jsonResponse.Status[i] = 1
		}
		encoder.Encode(jsonResponse)
	}
}

type GetNodeResponse struct {
	Ids     []string
	Command []string
	Number  []int
}

func handleGetNode(pool *slave.SlavePool) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		nodeId := params["nodeId"]
		encoder := json.NewEncoder(res)
		var slave *slave.Slave
		for _, sv := range pool.Slaves {
			if sv.Id == nodeId {
				slave = sv
			}
		}

		if slave == nil {
			http.NotFound(res, req)
			return
		}

		jsonResponse := GetNodeResponse{
			Ids:     make([]string, 0),
			Command: make([]string, 0),
			Number:  make([]int, 0),
		}

		for _, sv := range slave.ServicesMap {
			if sv != nil && len(sv.Command) > 0 {
				jsonResponse.Ids = append(jsonResponse.Ids, sv.Id)
				jsonResponse.Command = append(jsonResponse.Command, strings.Join(sv.Command, " "))
				jsonResponse.Number = append(jsonResponse.Number, sv.InstanceNum)
			}
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
	nodeHandle.HandleFunc("/{nodeId}/", handleGetNode(pool))

	http.Handle("/", r)

	server := &http.Server{
		Handler:      r,
		Addr:         strings.Join([]string{ListeningAddress, HTTPListeningPort}, ":"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	defer server.Close()

	log.Printf("http server is listening at %s:%s", ListeningAddress, HTTPListeningPort)
	log.Fatalln(server.ListenAndServe())
}
