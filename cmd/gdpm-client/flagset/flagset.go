package flagset

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/gdpm/cmd/gdpm-client/util"
)

type BaseFlagSet interface {
	FormMap() map[string]string
	Form() url.Values
	UrlPath() string
}

type CreateNewServiceFlagSet struct {
	Fs          *flag.FlagSet
	Command     string
	InstanceNum int
}

func NewCreateServiceFlagSet() *CreateNewServiceFlagSet {
	fs := &CreateNewServiceFlagSet{
		Fs: &flag.FlagSet{},
	}
	fs.Fs.IntVar(&fs.InstanceNum, "replicas", 1, "Specify the number of replicas of the given command, default to 1.")
	fs.Fs.IntVar(&fs.InstanceNum, "r", 1, "Specify the number of replicas of the given command, default to 1.")
	return fs
}

func (fs *CreateNewServiceFlagSet) FormMap() map[string]string {
	form := make(map[string]string)
	form["Command"] = fs.Command
	form["InstanceNum"] = fmt.Sprintf("%v", fs.InstanceNum)
	return form
}

func (fs *CreateNewServiceFlagSet) Form() url.Values {
	formMap := fs.FormMap()
	form := url.Values{}
	for k, v := range formMap {
		form.Add(k, v)
	}
	return form
}

func (fs *CreateNewServiceFlagSet) UrlPath() string {
	return util.URL("service/")
}

type UpdateServiceFlagSet struct {
	Fs          *flag.FlagSet
	ServiceId   string
	InstanceNum int
}

func NewUpdateServiceFlagSet() *UpdateServiceFlagSet {
	fs := &UpdateServiceFlagSet{
		Fs: &flag.FlagSet{},
	}
	fs.Fs.IntVar(&fs.InstanceNum, "replicas", 1, "Specify the number of replicas of the given command, default to 1.")
	fs.Fs.IntVar(&fs.InstanceNum, "r", 1, "Specify the number of replicas of the given command, default to 1.")
	return fs
}

func (fs *UpdateServiceFlagSet) FormMap() map[string]string {
	form := make(map[string]string)
	form["InstanceNum"] = fmt.Sprintf("%v", fs.InstanceNum)
	return form
}

func (fs *UpdateServiceFlagSet) Form() url.Values {
	formMap := fs.FormMap()
	form := url.Values{}
	for k, v := range formMap {
		form.Add(k, v)
	}
	return form
}

func (fs *UpdateServiceFlagSet) UrlPath() string {
	return util.URL(fmt.Sprintf("service/%s/", fs.ServiceId))
}
