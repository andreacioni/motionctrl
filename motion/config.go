package motion

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/andreacioni/motionctrl/utils"
	"github.com/andreacioni/motionctrl/version"
)

const (
	KeyValueRegex = "[a-zA-Z0-9_%\\/-]"

	configWriteRegex         = "Camera [0-9]+ write\nDone\n"
	configDefaultParserRegex = "(?m)^([^;#]" + KeyValueRegex + "+) (" + KeyValueRegex + "+)$"
	listConfigParserRegex    = "(?m)^(" + KeyValueRegex + "+) = (" + KeyValueRegex + "+)$"
	getConfigParserRegex     = "(" + KeyValueRegex + "+) = (" + KeyValueRegex + "+)\nDone"
	setConfigParserRegex     = "%s = %s\nDone"
)

const (
	ConfigDaemon                   = "daemon"
	ConfigSetupMode                = "setup_mode"
	ConfigWebControlPort           = "webcontrol_port"
	ConfigStreamPort               = "stream_port"
	ConfigStreamAuthMethod         = "stream_auth_method"
	ConfigStreamAuthentication     = "stream_authentication"
	ConfigWebControlHTML           = "webcontrol_html_output"
	ConfigWebControlParms          = "webcontrol_parms"
	ConfigWebControlAuthentication = "webcontrol_authentication"
	ConfigProcessIdFile            = "process_id_file"
	ConfigTargetDir                = "target_dir"
)

var (
	configReadOnlyParams = []string{ConfigDaemon,
		ConfigSetupMode,
		ConfigWebControlPort,
		ConfigStreamPort,
		ConfigStreamAuthMethod,
		ConfigStreamAuthentication,
		ConfigWebControlHTML,
		ConfigWebControlParms,
		ConfigWebControlAuthentication,
		ConfigProcessIdFile,
		ConfigTargetDir,
	}
	readOnlyConfig map[string]string
)

var ConfigTypeMapper = func(s string) interface{} {

	integer, err := strconv.Atoi(s)

	if err == nil {
		return integer
	}

	switch s {
	case "true":
		return true
	case "on":
		return true
	case "false":
		return false
	case "off":
		return false
	case "(null)":
		return nil
	default:
		return s
	}
}

var ReverseConfigTypeMapper = func(s string) interface{} {
	switch s {
	case "true":
		return "on"
	case "false":
		return "off"
	case "null":
		return ""
	default:
		return s
	}
}

func loadConfig(filename string) error {

	temp, err := parseConfig(filename)

	if err == nil {
		err = checkConfig(temp)

		if err == nil {
			readOnlyConfig = temp
		}
	}

	return err
}

func checkConfig(configMap map[string]string) error {
	webControlPort := configMap[ConfigWebControlPort]

	if webControlPort == "" {
		return fmt.Errorf("missing %s", ConfigWebControlPort)
	}

	streamPort := configMap[ConfigStreamPort]

	if streamPort == "" {
		return fmt.Errorf("'%s' parameter not found", ConfigStreamPort)
	}

	webControlHTML := configMap[ConfigWebControlHTML]

	if webControlHTML == "" || webControlHTML == "on" {
		return fmt.Errorf("'%s' parameter not found or set to 'on' (must be 'off')", ConfigWebControlHTML)
	}

	webControlParms := configMap[ConfigWebControlParms]

	if webControlParms == "" || webControlParms != "2" {
		return fmt.Errorf("web control authentication is enabled in motion config, please disable it (set to '2'). %s already has login features to protect your camera", version.Name)
	}

	webControlAuth := configMap[ConfigWebControlAuthentication]

	if webControlAuth != "" {
		return fmt.Errorf("'%s' parameter need to be commented in motion config", ConfigWebControlAuthentication)
	}

	streamAuthMethod := configMap[ConfigStreamAuthMethod]

	if streamAuthMethod != "0" {
		return fmt.Errorf("stream authentication is enabled in motion config, please disable it, %s already has login features to protect your camera", version.Name)
	}

	streamAuth := configMap[ConfigStreamAuthentication]

	if streamAuth != "" {
		return fmt.Errorf("'%s' parameter need to be commented in motion config", ConfigStreamAuthentication)
	}

	pidFile := configMap[ConfigProcessIdFile]

	if pidFile == "" {
		return fmt.Errorf("'%s' parameter not found", ConfigStreamAuthentication)
	}

	targetDir := configMap[ConfigTargetDir]

	if targetDir == "" {
		return fmt.Errorf("'%s' parameter not found", ConfigTargetDir)
	}

	return nil
}

//TODO improve with regexp
func parseConfig(configFile string) (map[string]string, error) {
	result := make(map[string]string)

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, ";") {
			lines := strings.Split(line, " ")
			if len(lines) >= 2 {
				result[lines[0]] = strings.Join(lines[1:], "")
			}
		}

	}

	return result, nil
}

func ConfigList() (map[string]interface{}, error) {
	ret, err := webControlGet("/config/list", func(body string) (interface{}, error) {
		ret := utils.RegexSubmatchTypedMap(listConfigParserRegex, body, ConfigTypeMapper)

		if len(ret) == 0 {
			return nil, fmt.Errorf("empty configuration")
		}
		return ret, nil
	})

	if err != nil {
		return nil, err
	}

	return ret.(map[string]interface{}), nil
}

func ConfigGetRO(param string) string {
	return readOnlyConfig[param]
}

func ConfigGet(param string) (map[string]interface{}, error) {
	queryURL := fmt.Sprintf("/config/get?query=%s", param)
	ret, err := webControlGet(queryURL, func(body string) (interface{}, error) {
		ret := utils.RegexSubmatchTypedMap(getConfigParserRegex, body, ConfigTypeMapper)

		if len(ret) == 0 {
			return nil, fmt.Errorf("invalid query (%s)", body)
		}
		return ret, nil
	})

	if err != nil {
		return nil, err
	}

	return ret.(map[string]interface{}), nil

}

func ConfigCanSet(name string) bool {
	b, _ := utils.InSlice(name, configReadOnlyParams)
	return !b
}

func ConfigSet(name string, value string) error {
	queryURL := fmt.Sprintf("/config/set?%s=%s", name, value)
	_, err := webControlGet(queryURL, func(body string) (interface{}, error) {
		if !utils.RegexMustMatch(fmt.Sprintf(setConfigParserRegex, name, value), body) {
			return nil, fmt.Errorf("there was an error on setting '%s' parameter", name)
		}

		return nil, nil
	})

	return err
}

func ConfigWrite() error {
	_, err := webControlGet("/config/write", func(body string) (interface{}, error) {
		if !utils.RegexMustMatch(configWriteRegex, body) {
			return nil, fmt.Errorf("unable to write config (%s)", body)
		}
		return nil, nil
	})

	return err
}
