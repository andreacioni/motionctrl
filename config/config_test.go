package config

import "testing"

func TestConfig(t *testing.T) {
	Load("github.com/andreacioni/motionctrl/config.json")
}
