package util

import "strings"

var MasterAddress string = "http://0.0.0.0"
var MasterPortNumber string = "8989"

func URL(path string) string {
	hostName := strings.Join([]string{MasterAddress, MasterPortNumber}, ":")
	return strings.Join([]string{hostName, path}, "/")
}
