package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvidedConfig(t *testing.T) {
	Load("../config.json")
}

func TestConfig1(t *testing.T) {
	Load("test_config_1.json")

	require.Equal(t, Backup{}, GetBackupConfig())
	require.Equal(t, Notify{}, GetNotifyConfig())
}
