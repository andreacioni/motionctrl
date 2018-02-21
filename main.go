package main

import (
	"flag"
	"fmt"

	"github.com/kpango/glg"

	"./api"
	"./config"
	"./motion"
	"./version"
)

var (
	configFile string
	logLevel   string
	autostart  bool
	detection  bool
)

func main() {
	fmt.Printf("%s is starting (version: %s build: %d)\n", version.Name, version.Number, version.Build)

	parseArgs()

	setupLogger()

	config.Load(configFile)

	motion.Init(config.Get().MotionConfigFile, autostart, detection)

	api.Init()
}

func setupLogger() {
	glg.Get().SetMode(glg.STD).AddStdLevel(logLevel, glg.STD, false)
}

func parseArgs() {
	flag.StringVar(&configFile, "c", "./config.json", "configuration file path")
	flag.StringVar(&logLevel, "l", "WARN", "set log level")
	flag.BoolVar(&autostart, "a", false, fmt.Sprintf("start motion right after %s", version.Name))
	flag.BoolVar(&detection, "d", false, "when -a is set, starts with motion detection enabled")

	flag.Parse()
}
