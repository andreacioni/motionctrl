package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kpango/glg"
)

type Configuration struct {
	Address       string `json:"address"`
	Port          string `json:"port"`
	RemoteAddress string `json:"remote_address"`
	ControlPort   int    `json:"control_port"`
	StreamPort    int    `json:"stream_port"`
}

var (
	Conf Configuration
)

// Load function convert a loaded JSON config file to a config struct
// return err if secret param is empty
func Load(filename string) error {
	glg.Infof("Loading configuration from %s ...", filename)
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, &Conf)
	if err != nil {
		return err
	}

	glg.Debugf("Current config: %+v", Conf)

	return nil
}
