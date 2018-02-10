package motion

import (
	"testing"
)

func TestConfigParser(t *testing.T) {
	Parse("motion_test.conf")
}

func TestNotPresentParser(t *testing.T) {
	configMap, _ := Parse("motion_test.conf")

	value := configMap["this_is_not_present"]

	if value != "" {
		t.Error("Nil value")
	}
}

func TestCheck(t *testing.T) {
	configMap, _ := Parse("motion_test.conf")
	err := Check(configMap)

	if err != nil {
		t.Error(err)
	}

	if len(configMap) != 89 {
		t.Error("Configuration parameters map must contain 89 elements")
	}
}

func TestCheckInstall(t *testing.T) {
	err := CheckInstall()

	if err != nil {
		t.Error(err)
	}
}

func TestStartStop(t *testing.T) {
	Init("./motion_test.conf")

	err := Startup(false)

	if err != nil {
		t.Error(err)
	}

	if !IsStarted() {
		t.Error("not started")
	}

	err = Shutdown()

	if err != nil {
		t.Error(err)
	}
}

func TestRestart(t *testing.T) {

	Init("./motion_test.conf")

	err := Startup(false)

	if err != nil {
		t.Error(err)
	}

	if !IsStarted() {
		t.Error("not started")
	}

	err = Restart()

	if err != nil {
		t.Error(err)
	}

	err = Shutdown()

	if err != nil {
		t.Error(err)
	}
}
