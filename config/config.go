package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"sync"

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
	AppPath          string `json:"password"`
	Ssl              SSL    `json:"ssl"`
	Backup           Backup `json:"backup"`
	Notify           Notify `json:"notify"`
}

type SSL struct {
	CertFile string `json:"cert"`
	KeyFile  string `json:"key"`
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

var (
	mu   sync.Mutex
	conf Configuration
)

// Load function convert a loaded JSON config file to a config struct
// return err if secret param is empty
func Load(filename string) error {
	mu.Lock()
	defer mu.Unlock()

	glg.Infof("Loading configuration from %s ...", filename)

	if conf.IsEmpty() {
		raw, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		err = json.Unmarshal(raw, &conf)
		if err != nil {
			return err
		}

		glg.Debugf("Current config: %+v", conf)
	} else {
		return fmt.Errorf("Configuration already loaded")
	}

	return nil
}

func Unload() {
	mu.Lock()
	defer mu.Unlock()

	glg.Debug("Unloading configuration")

	conf = Configuration{}
}

func GetConfig() Configuration {
	mu.Lock()
	defer mu.Unlock()

	return conf
}

func GetSSLConfig() SSL {
	mu.Lock()
	defer mu.Unlock()

	return conf.Ssl
}

func GetBackupConfig() Backup {
	mu.Lock()
	defer mu.Unlock()

	return conf.Backup
}

func GetNotifyConfig() Notify {
	mu.Lock()
	defer mu.Unlock()

	return conf.Notify
}

func (c Configuration) IsEmpty() bool {
	return reflect.DeepEqual(c, Configuration{})
}

func (c SSL) IsEmpty() bool {
	return reflect.DeepEqual(c, SSL{})
}

func (c Notify) IsEmpty() bool {
	return reflect.DeepEqual(c, Notify{})
}

func (c Backup) IsEmpty() bool {
	return reflect.DeepEqual(c, Backup{})
}
