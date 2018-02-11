package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kpango/glg"
)

const (
	BaseAddress = "127.0.0.1"
)

type Configuration struct {
	Address          string `json:"address"`
	Port             int    `json:"port"`
	MotionConfigFile string `json:"motionConfigFile"`
	Username         string `json:"username"`
	Password         string `json:"password"`
}

var conf Configuration

// Load function convert a loaded JSON config file to a config struct
// return err if secret param is empty
func Load(filename string) {
	glg.Infof("Loading configuration from %s ...", filename)
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		glg.Fatal(err)
	}

	err = json.Unmarshal(raw, &conf)
	if err != nil {
		glg.Fatal(err)
	}

	glg.Debugf("Current config: %+v", conf)
}

func Get() Configuration {
	return conf
}
