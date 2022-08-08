package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ContentRoot     string   `json:"content_root"`
	DNSAddress      string   `json:"dns_address"`
	UpstreamServers []string `json:"upstream_servers"`
	DNSPort         int      `json:"dns_port"`
	UpstreamTimeout int      `json:"upstream_timeout"`
}

var Settings Config

func (c *Config) Load() (err error) {
	// ToDo: change location
	configLocations := []string{
		"config.json",
	}
	var configBytes []byte

	for _, location := range configLocations {
		configBytes, err = ioutil.ReadFile(location)
		if err != nil {
			continue
		}
		break
	}

	if err == nil {
		err = json.Unmarshal(configBytes, &c)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}
