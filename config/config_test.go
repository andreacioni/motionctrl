package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvidedConfig(t *testing.T) {
	Load("../config.json")
	Unload()
}

func TestConfig1(t *testing.T) {
	Load("test_config_1.json")

	require.Equal(t, Backup{}, GetBackupConfig())
	require.Equal(t, Notify{}, GetNotifyConfig())

	Unload()
}

func TestConfigNotUnloaded(t *testing.T) {
	require.NoError(t, Load("test_config_1.json"))

	require.Error(t, Load("test_config_1.json"))

	Unload()
}

func TestConfigEmpty(t *testing.T) {
	require.True(t, Configuration{}.IsEmpty())
	require.True(t, Backup{}.IsEmpty())
	require.True(t, Notify{}.IsEmpty())

	require.False(t, Configuration{Address: "101.1.1.1"}.IsEmpty())
	require.False(t, Backup{Method: "test"}.IsEmpty())
	require.False(t, Notify{Method: "test"}.IsEmpty())
}
