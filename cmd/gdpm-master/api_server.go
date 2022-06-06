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
