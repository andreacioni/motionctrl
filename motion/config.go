package motion

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/kpango/glg"
)

const (
	pat = "[^;#][[:word:]]+[[:blank:]]+[[:word:]]+"
)

var (
	re            *regexp.Regexp
	motionConfMap map[string]string
)

func init() {
	re = regexp.MustCompile(pat)
}

func Load(filename string) {
	glg.Infof("Loading motion configuration from %s...", filename)

	motionConfMap = Parse(filename)
}

func Parse(configFile string) map[string]string {
	result := make(map[string]string)

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
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

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

func Get(key string) string {
	return motionConfMap[key]
}
