package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kpango/glg"
)

type Authentication struct {
}

type Configuration struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var Conf Configuration

// Load function convert a loaded JSON config file to a config struct
// return err if secret param is empty
func Load(filename string) error {
	glg.Infof("Loading configuration from %s ...", filename)
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		glg.Fatal(err)
	}

	err = json.Unmarshal(raw, &Conf)
	if err != nil {
		glg.Fatal(err)
	}

	glg.Debugf("Current config: %+v", Conf)

	return nil
}
