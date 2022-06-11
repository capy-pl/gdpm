package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gdpm/cmd/gdpm-client/flagset"
	"github.com/gdpm/cmd/gdpm-client/util"
)

var MasterAddress string = "http://0.0.0.0"
var MasterPortNumber string = "8989"

func CreateService(args []string) {
	fs := flagset.NewCreateServiceFlagSet()
	err := fs.Fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	if fs.Fs.NArg() == 0 {
		log.Fatal("Please provide the executable.")
	}

	fs.Command = strings.Join(fs.Fs.Args(), " ")
	log.Printf("%s %v\n", fs.Command, fs.InstanceNum)
	PostService(fs)
}

func PostService(flg flagset.BaseFlagSet) {
	client := &http.Client{}
	form := flg.Form()
	path := flg.UrlPath()

	log.Printf("[service] post request, url %s\n", path)

	res, err := client.PostForm(path, form)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	log.Printf("[service] post response, status %v\n", res.StatusCode)
	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyStr := string(bodyBytes)
		log.Printf("[service] post response, body %s\n", bodyStr)
	} else {
		log.Fatalln(res.StatusCode)
	}
}

func UpdateNewService(args []string) {
	fs := flagset.NewUpdateServiceFlagSet()
	err := fs.Fs.Parse(args)
	if err != nil {
		log.Fatalln(err)
	}
	if fs.Fs.NArg() == 0 {
		log.Fatalln("missing instance id")
	}
	id := fs.Fs.Arg(0)
	fs.ServiceId = id
	log.Printf("[update] instance id %s, replica numbers %v", fs.ServiceId, fs.InstanceNum)
	PostService(fs)
}

func DeleteService(args []string) {
	if len(args) < 1 {
		log.Fatalln("please provide a service id")
	}
	id := args[0]
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", util.URL(fmt.Sprintf("service/%s/", id)), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[delete] service id = %s\n", id)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("[delete] response status = %s\n", res.Status)
}

type GetNodesResponse struct {
	Ids        []string
	ServiceNum []int
	Status     []int
	Times      []string
}

func ListNodes(args []string) {
	client := &http.Client{}
	res, err := client.Get(util.URL("node/"))
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(res.Body)
	nodeList := GetNodesResponse{}
	decoder.Decode(&nodeList)
	fmt.Printf("%-36s %-15s %-19s\n", "node id", "service number", "start time")
	for i := 0; i < len(nodeList.Ids); i++ {
		fmt.Printf("%-36s %-15d %-19s\n", nodeList.Ids[i], nodeList.ServiceNum[i], nodeList.Times[i])
	}
}

type ListNodeServicesResponse struct {
	Ids     []string
	Command []string
	Number  []int
}

func ListNodeServices(args []string) {
	client := &http.Client{}
	res, err := client.Get(util.URL(fmt.Sprintf("node/%s/", args[0])))
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(res.Body)
	svList := &ListNodeServicesResponse{}
	decoder.Decode(&svList)
	fmt.Printf("%-36s %-10s %s\n", "service id", "number", "command")
	for i := 0; i < len(svList.Ids); i++ {
		fmt.Printf("%-36s %-10v %s\n", svList.Ids[i], svList.Number[i], svList.Command[i])
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please input subcommand.")
	}
	switch os.Args[1] {
	case "create":
		CreateService(os.Args[2:])
	case "update":
		UpdateNewService(os.Args[2:])
	case "delete":
		DeleteService(os.Args[2:])
	case "get":
		if len(os.Args) < 3 {
			log.Fatal("Please input a valid target.")
		}
		switch os.Args[2] {
		case "nodes":
			ListNodes(os.Args[2:])
		case "node":
			ListNodeServices(os.Args[3:])
		default:
			log.Fatal("Not a valid target.")
		}
	default:
		log.Fatal("Not a valid command.")
	}
}
