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
	Backup           Backup `json:"backup"`
	Notify           Notify `json:"notify"`
}

type Backup struct {
	When           string `json:"when"`
	Method         string `json:"method"`
	EncryptionKey  string `json:"encryptionKey"`
	Archive        bool   `json:"archive"`
	FilePerArchive int    `json:"filePerArchive"`
}

type Notify struct {
	Method  string   `json:"method"`
	Token   string   `json:"token"`
	To      []string `json:"to"`
	Message string   `json:"message"`
	Photo   int      `json:"photo"`
}

var conf Configuration

// Load function convert a loaded JSON config file to a config struct
// return err if secret param is empty
func Load(filename string) error {
	glg.Infof("Loading configuration from %s ...", filename)
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, &conf)
	if err != nil {
		return err
	}

	glg.Debugf("Current config: %+v", conf)

	return nil
}

func GetConfig() Configuration {
	return conf
}

func GetBackupConfig() Backup {
	return conf.Backup
}

func GetNotifyConfig() Notify {
	return conf.Notify
}
