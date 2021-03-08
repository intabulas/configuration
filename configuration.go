package configuration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/intabulas/configuration/hocon"
)

func ParseString(text string, includeCallback ...hocon.IncludeCallback) (*Config, error) {
	var callback hocon.IncludeCallback
	if len(includeCallback) > 0 {
		callback = includeCallback[0]
	} else {
		callback = defaultIncludeCallback
	}
	root, err := hocon.Parse(text, callback)
	if err != nil {
		return nil, err
	}

	return NewConfigFromRoot(root), nil
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParseString(string(data), defaultIncludeCallback)
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

	return ParseString(string(data), callback)
}

func FromObject(obj interface{}) (*Config, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return ParseString(string(data), defaultIncludeCallback)
}

func defaultIncludeCallback(filename string) (*hocon.HoconRoot, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return hocon.Parse(string(data), defaultIncludeCallback)
}
