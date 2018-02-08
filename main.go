package main

import (
	"flag"
	"fmt"
	"os/exec"

	"github.com/kpango/glg"
	"github.com/parnurzeal/gorequest"

	"./api"
	"./config"
	"./version"
)

var (
	configFile string
	logLevel   string
)

func main() {
	fmt.Printf("%s is starting (version: %s build: %d)\n", version.Name, version.Number, version.Build)

	parseArgs()

	setupLogger()

	loadConfig()

	checkMotionInstalled()

	api.Init()
}

func checkMotionInstalled() {
	err := exec.Command("motion", "-h").Run()

	//TODO unfortunatelly motion doesn't return 0 when invoked with the "-h" parameter
	if err != nil && err.Error() != "exit status 1" {
		glg.Fatalf("Motion not found (%s)", err)
	}

	glg.Debug("Motion was found")
}

func setupLogger() {
	glg.Get().SetMode(glg.STD).AddStdLevel(logLevel, glg.STD, false)
}

func loadConfig() {
	err := config.Load(configFile)
	if err != nil {
		glg.Fatal(err)
	}
}

func parseArgs() {
	flag.StringVar(&configFile, "c", "./config.json", "configuration file path")
	flag.StringVar(&logLevel, "l", "WARN", "set log level")

	flag.Parse()
}

func getStream() {
	res, _, errs := gorequest.New().Get(fmt.Sprintf("http://%s:%d/", config.Conf.RemoteAddress, config.Conf.ControlPort)).End()
	if errs != nil {
		glg.Error(errs)
	} else {
		glg.Debugf("%+v", *res)
	}

}

func getControl() {
	res, _, errs := gorequest.New().Get(fmt.Sprintf("http://%s:%d/", config.Conf.RemoteAddress, config.Conf.StreamPort)).End()
	if errs != nil {
		glg.Fatal(errs)
	} else {
		glg.Debugf("%+v", *res)
	}

}
