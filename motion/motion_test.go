package motion

import (
	"fmt"
	"testing"

	"../utils"
	"github.com/stretchr/testify/assert"
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

func TestConfigTypeMapper(t *testing.T) {
	testMap := map[string]string{
		"value1": "text",
		"value2": "off",
		"value3": "on",
		"value4": "3",
	}

	assert.Equal(t, "text", ConfigTypeMapper(testMap["value1"]))
	assert.Equal(t, false, ConfigTypeMapper(testMap["value2"]))
	assert.Equal(t, true, ConfigTypeMapper(testMap["value3"]))
	assert.Equal(t, 3, ConfigTypeMapper(testMap["value4"]))
}

func TestRegexConfigFileParser(t *testing.T) {
	testString := "#comment here\n;comment here\nhello 12\nword 11\nnullparam (null)\nonoff on\noffon off"

	testMap := utils.RegexSubmatchTypedMap(ConfigDefaultParserRegex, testString, ConfigTypeMapper)

	assert.Equal(t, 4, len(testMap))
	assert.Equal(t, 12, testMap["hello"])
	assert.Equal(t, 11, testMap["word"])
	assert.Empty(t, testMap["nullparam"])
	assert.Equal(t, true, testMap["onoff"])
	assert.Equal(t, false, testMap["offon"])
}

func TestRegexConfigList(t *testing.T) {
	testString := "#comment = here\n;comment = here\nhello = 12\nword = 11\nnullparam = (null)\nonoff = on\noffon = off"

	testMap := utils.RegexSubmatchTypedMap(ListConfigParserRegex, testString, ConfigTypeMapper)

	assert.Equal(t, 4, len(testMap))
	assert.Equal(t, 12, testMap["hello"])
	assert.Equal(t, 11, testMap["word"])
	assert.Empty(t, testMap["nullparam"])
	assert.Equal(t, true, testMap["onoff"])
	assert.Equal(t, false, testMap["offon"])
}

func TestRegexSetRegex(t *testing.T) {
	testString := "testparam = Hello\nDone"
	assert.True(t, utils.RegexMustMatch(fmt.Sprintf(SetConfigParserRegex, "testparam", "Hello"), testString))
}
