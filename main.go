package main

import (
	"flag"
	"fmt"

	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/api"
	"github.com/andreacioni/motionctrl/backup"
	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/motion"
	"github.com/andreacioni/motionctrl/version"
)

var (
	configFile string
	logLevel   string
	autostart  bool
	detection  bool
)

func main() {
	fmt.Printf("%s is starting (version: %s build: %d)\n", version.Name, version.Number, version.Build)

	//Parse command line arguments
	parseArgs()

	//Setup logger
	setupLogger()

	//Load motionctrl configuration file
	config.Load(configFile)

	//Initialize motion package
	motion.Init(config.GetConfig().MotionConfigFile, autostart, detection)

	//Initialize REST api
	api.Init()

	//Initialize backup  (if enabled)
	backup.Init(config.GetBackupConfig(), motion.ConfigGetRO(motion.ConfigTargetDir))
}

func setupLogger() {
	glg.Get().SetMode(glg.STD).AddStdLevel(logLevel, glg.STD, false)
}

func parseArgs() {
	flag.StringVar(&configFile, "c", "github.com/andreacioni/motionctrl/config.json", "configuration file path")
	flag.StringVar(&logLevel, "l", "WARN", "set log level")
	flag.BoolVar(&autostart, "a", false, fmt.Sprintf("start motion right after %s", version.Name))
	flag.BoolVar(&detection, "d", false, "when -a is set, starts with motion detection enabled")

	flag.Parse()
}
