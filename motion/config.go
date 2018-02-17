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
	KeyValueRegex            = "[a-zA-Z0-9_%\\/-]"
	ConfigWriteRegex         = "Camera [0-9]+ write\nDone\n"
	ConfigDefaultParserRegex = "(?m)^([^;#]" + KeyValueRegex + "+) (" + KeyValueRegex + "+)$"
	ListConfigParserRegex    = "(?m)^(" + KeyValueRegex + "+) = (" + KeyValueRegex + "+)$"
	GetConfigParserRegex     = "(" + KeyValueRegex + "+) = (" + KeyValueRegex + "+)\nDone"
	SetConfigParserRegex     = "%s = %s\nDone"
)

const (
	Daemon                   = "daemon"
	SetupMode                = "setup_mode"
	WebControlPort           = "webcontrol_port"
	StreamPort               = "stream_port"
	StreamAuthMethod         = "stream_auth_method"
	StreamAuthentication     = "stream_authentication"
	WebControlHTML           = "webcontrol_html_output"
	WebControlParms          = "webcontrol_parms"
	WebControlAuthentication = "webcontrol_authentication"
	ProcessIdFile            = "process_id_file"
)

var (
	ConfigReadOnlyParams = []string{Daemon,
		SetupMode,
		WebControlPort,
		StreamPort,
		StreamAuthMethod,
		StreamAuthentication,
		WebControlHTML,
		WebControlParms,
		WebControlAuthentication,
		ProcessIdFile,
	}
	motionConfMap map[string]string
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

	integer, err := strconv.Atoi(s)

	if err == nil {
		return integer
	}

	switch s {
	case "on":
		return true
	case "off":
		return false
	case "(null)":
		return nil
	default:
		return s
	}
}

func Load(filename string) error {

	temp, err := Parse(filename)

	if err == nil {
		err = Check(temp)

		if err == nil {
			motionConfMap = temp
		}
	}

	return err
}

func Check(configMap map[string]string) error {
	webControlPort := configMap[WebControlPort]

	if webControlPort == "" {
		return fmt.Errorf("missing %s", WebControlPort)
	}

	streamPort := configMap[StreamPort]

	if streamPort == "" {
		return fmt.Errorf("'%s' parameter not found", StreamPort)
	}

	webControlHTML := configMap[WebControlHTML]

	if webControlHTML == "" || webControlHTML == "on" {
		return fmt.Errorf("'%s' parameter not found or set to 'on' (must be 'off')", WebControlHTML)
	}

	webControlParms := configMap[WebControlParms]

	if webControlParms == "" || webControlParms != "2" {
		return fmt.Errorf("web control authentication is enabled in motion config, please disable it (set to '2'). %s already has login features to protect your camera", version.Name)
	}

	webControlAuth := configMap[WebControlAuthentication]

	if webControlAuth != "" {
		return fmt.Errorf("'%s' parameter need to be commented in motion config", WebControlAuthentication)
	}

	streamAuthMethod := configMap[StreamAuthMethod]

	if streamAuthMethod != "0" {
		return fmt.Errorf("stream authentication is enabled in motion config, please disable it, %s already has login features to protect your camera", version.Name)
	}

	streamAuth := configMap[StreamAuthentication]

	if streamAuth != "" {
		return fmt.Errorf("'%s' parameter need to be commented in motion config", StreamAuthentication)
	}

	pidFile := configMap[ProcessIdFile]

	if pidFile == "" {
		return fmt.Errorf("'%s' parameter not found", StreamAuthentication)
	}

	return nil
}

//TODO improve with regexp
func Parse(configFile string) (map[string]string, error) {
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
		ret := utils.RegexSubmatchTypedMap(ListConfigParserRegex, body, ConfigTypeMapper)

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

func ConfigGet(param string) (map[string]interface{}, error) {
	queryURL := fmt.Sprintf("/config/get?query=%s", param)
	ret, err := webControlGet(queryURL, func(body string) (interface{}, error) {
		ret := utils.RegexSubmatchTypedMap(GetConfigParserRegex, body, ConfigTypeMapper)

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
	b, _ := utils.InSlice(name, ConfigReadOnlyParams)

	return !b
}

func ConfigSet(name string, value string) error {
	queryURL := fmt.Sprintf("/config/set?%s=%s", name, value)
	_, err := webControlGet(queryURL, func(body string) (interface{}, error) {
		if !utils.RegexMustMatch(fmt.Sprintf(SetConfigParserRegex, name, value), body) {
			return nil, fmt.Errorf("there was an error on setting '%s' parameter", name)
		}

		return nil, nil
	})

	return err
}

func ConfigWrite() error {
	_, err := webControlGet("/config/write", func(body string) (interface{}, error) {
		if !utils.RegexMustMatch(ConfigWriteRegex, body) {
			return nil, fmt.Errorf("unable to write config (%s)", body)
		}
		return nil, nil
	})

	return err
}
