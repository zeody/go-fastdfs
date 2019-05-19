package go_fastdfs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type TrackerConf struct {
	ConnectTimeout     int    `yaml:"connect_timeout_in_seconds"`
	NetworkTimeout     int    `yaml:"network_timeout_in_seconds"`
	Charset            string `yaml:"charset"`
	HttpAntiStealToken bool   `yaml:"http_anti_steal_token"`
	HttpSecretKey      string `yaml:"http_secret_key"`
	HttpTrackerPort    int    `yaml:"http_tracker_http_port"`
	TrackerServers     string `yaml:"tracker_servers"`
}

func ParseYaml() *TrackerConf {
	file, err := ioutil.ReadFile("fdfs.yml")
	if err != nil {
		log.Printf("cannot get yaml file #%v", err)
	}
	conf := new(TrackerConf)
	err = yaml.Unmarshal(file, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return conf
}
