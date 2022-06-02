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

	log.Println(fs.Fs.NArg())
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
}

type GetNodesResponse struct {
	Ids        []string
	ServiceNum []int
	Status     []int
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
	for i := 0; i < len(nodeList.Ids); i++ {
		log.Println(nodeList.Ids[i])
	}
}

func ListNodeServices(args []string) {

}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please input subcommand.")
	}
	fmt.Printf("%v\n", os.Args)
	switch os.Args[1] {
	case "create":
		CreateService(os.Args[2:])
	case "update":
		UpdateNewService(os.Args[2:])
	case "get":
		if len(os.Args) < 3 {
			log.Fatal("Please input a valid target.")
		}
		switch os.Args[2] {
		case "nodes":
			ListNodes(os.Args[2:])
		case "services":
			ListNodeServices(os.Args[2:])
		default:
			log.Fatal("Not a valid target.")
		}
	default:
		log.Fatal("Not a valid command.")
	}
}
