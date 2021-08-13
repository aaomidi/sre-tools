package config

import (
	"io/ioutil"

	"github.com/pelletier/go-toml/v2"
)

type Configuration struct {
	ZoneIdentifier string `toml:"zone_identifier"`
	Records        map[string]Record
}

type Record struct {
	Type    string
	Name    string
	Content string
	Ttl     int
	Proxied bool
	Exist   bool
}

func ReadConfigurationFile(path string) (Configuration, error) {
	var configuration Configuration
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return configuration, nil
	}

	err = toml.Unmarshal(bytes, &configuration)

	if err != nil {
		return configuration, err
	}
	return configuration, nil
}
