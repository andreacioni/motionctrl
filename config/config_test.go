package config

import "testing"

func TestConfig(t *testing.T) {
	Load("../config.json")
}
