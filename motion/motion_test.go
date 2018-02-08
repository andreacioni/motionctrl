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
}

func TestCheckInstall(t *testing.T) {
	err := CheckInstall()

	if err != nil {
		t.Error(err)
	}
}
