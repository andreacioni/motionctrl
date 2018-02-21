package motion

import (
	"fmt"
	"testing"

	"github.com/andreacioni/motionctrl/utils"

	"github.com/stretchr/testify/require"
)

func TestConfigParser(t *testing.T) {
	parseConfig("motion_test.conf")
}

func TestNotPresentParser(t *testing.T) {
	configMap, _ := parseConfig("motion_test.conf")

	value := configMap["this_is_not_present"]

	if value != "" {
		t.Error("Nil value")
	}
}

func TestCheck(t *testing.T) {
	configMap, _ := parseConfig("motion_test.conf")
	err := checkConfig(configMap)

	if err != nil {
		t.Error(err)
	}

	if len(configMap) != 89 {
		t.Error("Configuration parameters map must contain 89 elements")
	}
}

func TestCheckInstall(t *testing.T) {
	err := checkInstall()

	if err != nil {
		t.Error(err)
	}
}

func TestStartStop(t *testing.T) {
	Init("github.com/andreacioni/motionctrl/motion_test.conf", false, false)

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

	Init("github.com/andreacioni/motionctrl/motion_test.conf", false, false)

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

func TestConfigTypeMapper(t *testing.T) {
	testMap := map[string]string{
		"value1": "text",
		"value2": "off",
		"value3": "on",
		"value4": "3",
	}

	require.Equal(t, "text", ConfigTypeMapper(testMap["value1"]))
	require.Equal(t, false, ConfigTypeMapper(testMap["value2"]))
	require.Equal(t, true, ConfigTypeMapper(testMap["value3"]))
	require.Equal(t, 3, ConfigTypeMapper(testMap["value4"]))
}

func TestRegexConfigFileParser(t *testing.T) {
	testString := "#comment here\n;comment here\nhello 12\nword 11\nnullparam (null)\nonoff on\noffon off"

	testMap := utils.RegexSubmatchTypedMap(configDefaultParserRegex, testString, ConfigTypeMapper)

	require.Equal(t, 4, len(testMap))
	require.Equal(t, 12, testMap["hello"])
	require.Equal(t, 11, testMap["word"])
	require.Empty(t, testMap["nullparam"])
	require.Equal(t, true, testMap["onoff"])
	require.Equal(t, false, testMap["offon"])
}

func TestRegexConfigList(t *testing.T) {
	testString := "#comment = here\n;comment = here\nhello = 12\nword = 11\nnullparam = (null)\nonoff = on\noffon = off"

	testMap := utils.RegexSubmatchTypedMap(listConfigParserRegex, testString, ConfigTypeMapper)

	require.Equal(t, 4, len(testMap))
	require.Equal(t, 12, testMap["hello"])
	require.Equal(t, 11, testMap["word"])
	require.Empty(t, testMap["nullparam"])
	require.Equal(t, true, testMap["onoff"])
	require.Equal(t, false, testMap["offon"])
}

func TestRegexSetRegex(t *testing.T) {
	testString := "testparam = Hello\nDone"
	testURL := "/config/set?daemon=true"

	require.True(t, utils.RegexMustMatch(fmt.Sprintf(setConfigParserRegex, "testparam", "Hello"), testString))

	mapped := utils.RegexSubmatchTypedMap("/config/set\\?("+KeyValueRegex+"+)=("+KeyValueRegex+"+)", testURL, nil)
	require.Equal(t, 1, len(mapped))

}

func TestParticularStartAndStop(t *testing.T) {
	Init("github.com/andreacioni/motionctrl/motion_test.conf", false, false)

	require.NoError(t, Startup(false))

	require.True(t, IsStarted())

	ret, err := ConfigGet("log_level") //Changing daemon instead of 'log_level' cause Shutdown to fail

	require.NoError(t, err)
	require.Equal(t, 6, ret["log_level"].(int))

	err = ConfigSet("log_level", "5")

	require.NoError(t, err)

	ret, err = ConfigGet("log_level")

	require.NoError(t, err)
	require.Equal(t, 5, ret["log_level"].(int))

	require.NoError(t, Shutdown())
}

func TestSomeConfigs(t *testing.T) {
	conf, err = parseConfig("motion_test.conf")
	require.Equals(t, "/tmp", conf[TargetDir])
}
