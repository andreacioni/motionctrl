package motion

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	WebControlPort = "webcontrol_port"
	StreamPort     = "stream_port"
)

var (
	motionConfMap map[string]string
)

func Load(filename string) error {

	temp, err := Parse(filename)

	if err == nil {
		err = Check(temp)
	}

	return err
}

func Check(configMap map[string]string) error {
	var err error
	webControlPort := configMap[WebControlPort]

	if webControlPort == "" {
		err = fmt.Errorf("missing %s", WebControlPort)
	}

	streamPort := configMap[StreamPort]

	if streamPort == "" {
		err = fmt.Errorf("missing %s", StreamPort)
	}

	return err
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
			if len(lines) == 2 {
				result[lines[0]] = lines[1]
			}
		}

	}

	return result, nil
}

func Get(key string) string {
	return motionConfMap[key]
}
