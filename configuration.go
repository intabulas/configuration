package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/intabulas/configuration/hocon"
)

func ParseString(text string, includeCallback ...hocon.IncludeCallback) *Config {
	var callback hocon.IncludeCallback
	if len(includeCallback) > 0 {
		callback = includeCallback[0]
	} else {
		callback = defaultIncludeCallback
	}
	root, err := hocon.Parse(text, callback)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("")
	// fmt.Println("")
	// // fmt.Printf("%+v\n", root.Value().GetObject().GetKeys())
	// fmt.Printf("%+v\n", root.Value().GetChildObject("monitor"))
	// fmt.Println("")
	// fmt.Println("")
	return NewConfigFromRoot(root)
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParseString(string(data), defaultIncludeCallback), nil
}

func LoadConfigWithIncludeCallback(filename string, includeCallback ...hocon.IncludeCallback) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var callback hocon.IncludeCallback
	if len(includeCallback) > 0 {
		callback = includeCallback[0]
	} else {
		callback = defaultIncludeCallback
	}

	return ParseString(string(data), callback), nil
}

func FromObject(obj interface{}) (*Config, error) {
	data, err := json.Marshal(obj)

	return ParseString(string(data), defaultIncludeCallback), err
}

func defaultIncludeCallback(filename string) (*hocon.HoconRoot, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return hocon.Parse(string(data), defaultIncludeCallback)
}
