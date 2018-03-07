package main

import (
	"flag"
	"fmt"

	"github.com/kpango/glg"

	"github.com/andreacioni/motionctrl/api"
	"github.com/andreacioni/motionctrl/backup"
	"github.com/andreacioni/motionctrl/config"
	"github.com/andreacioni/motionctrl/motion"
	"github.com/andreacioni/motionctrl/notify"
	"github.com/andreacioni/motionctrl/version"
)

var (
	configFile string
	logLevel   string
	autostart  bool
	detection  bool
)

func main() {
	fmt.Printf("%s is starting (version: %s)\n", version.Name, version.Number)

	//Parse command line arguments
	parseArgs()

	//Setup logger
	setupLogger()

	//Load motionctrl configuration file
	if err := config.Load(configFile); err != nil {
		glg.Fatalf("Error loading configuration: %v", err)
	}

	//Initialize motion package
	if err := motion.Init(config.GetConfig().MotionConfigFile, autostart, detection); err != nil {
		glg.Fatalf("Error initializing motion package: %v", err)
	}

	//Initialize backup  (if enabled)
	if err := backup.Init(config.GetBackupConfig(), motion.ConfigGetRO(motion.ConfigTargetDir)); err != nil {
		glg.Errorf("Error initializing backup package: %v", err)
	}

	//Initialize notify  (if enabled)
	if err := notify.Init(config.GetNotifyConfig()); err != nil {
		glg.Errorf("Error initializing notify package: %v", err)
	}

	//Initialize REST api
	api.Init()
}

func setupLogger() {
	glg.Get().SetMode(glg.STD).AddStdLevel(logLevel, glg.STD, false)
}

func parseArgs() {
	flag.StringVar(&configFile, "c", "config.json", "configuration file path")
	flag.StringVar(&logLevel, "l", "WARN", "set log level")
	flag.BoolVar(&autostart, "a", false, fmt.Sprintf("start motion right after %s", version.Name))
	flag.BoolVar(&detection, "d", false, "when -a is set, starts with motion detection enabled")

	flag.Parse()
}
