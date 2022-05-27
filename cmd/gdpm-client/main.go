package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var MasterAddress string = "http://0.0.0.0"
var MasterPortNumber string = "8989"

func toURL(path string) string {
	hostName := strings.Join([]string{MasterAddress, MasterPortNumber}, ":")
	return strings.Join([]string{hostName, path}, "/")
}

type CreateNewServiceFlagSet struct {
	fs          *flag.FlagSet
	Command     string
	InstanceNum int
}

func NewCreateServiceFlagSet() *CreateNewServiceFlagSet {
	fs := &CreateNewServiceFlagSet{
		fs: &flag.FlagSet{},
	}
	fs.fs.IntVar(&fs.InstanceNum, "replicas", 1, "Specify the number of replicas of the given command, default to 1.")
	fs.fs.IntVar(&fs.InstanceNum, "r", 1, "Specify the number of replicas of the given command, default to 1.")
	return fs
}

func CreateService(args []string) {
	fs := NewCreateServiceFlagSet()
	err := fs.fs.Parse(args)
	if err != nil {
		log.Fatal(err)
	}
	if fs.fs.NArg() == 0 {
		log.Fatal("Please provide the executable.")
	}
	log.Println(fs.fs.NArg())
	fs.Command = strings.Join(fs.fs.Args(), " ")
	log.Printf("%s %v\n", fs.Command, fs.InstanceNum)
	PostNewService(fs.Command, fs.InstanceNum)
}

func PostNewService(command string, instancenum int) {
	client := &http.Client{}
	form := url.Values{}
	form.Add("Command", command)
	form.Add("InstanceNum", fmt.Sprintf("%v", instancenum))
	log.Printf("[create] post request, command %s ,replicas %v \n", command, instancenum)
	res, err := client.PostForm(toURL("service/"), form)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[create] post response status %v\n", res.StatusCode)
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyStr := string(bodyBytes)
		log.Printf("[create] post response, service id %s\n", bodyStr)
	} else {
		log.Fatalln(res.StatusCode)
	}
}

func UpdateNewService(args []string) {

}

func ListNodes(args []string) {

}

func ListNodeServices(args []string) {

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
