package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"

	"../config"
	"../motion"
	"../utils"
)

var handlersMap = map[string]func(*gin.Context){
	"/control/startup":  startHandler,
	"/control/shutdown": stopHandler,
	"/control/status":   statusHandler,
	"/control/restart":  restartHandler,
	"/detection/status": isMotionDetectionEnabled,
	"/detection/start":  startDetectionHandler,
	"/detection/pause":  pauseDetectionHandler,
	"/camera":           utils.ReverseProxy("http://localhost:8081"),
	"/config/list":      listConfigHandler,
	"/config/set":       setConfigHandler,
	"/config/get":       getConfigHandler,
}

func Init() {
	glg.Info("Initializing REST API ...")
	var group *gin.RouterGroup
	r := gin.Default()

	if config.Get().Username != "" && config.Get().Password != "" {
		glg.Info("Username and password defined, authentication enabled")
		group = r.Group("/", gin.BasicAuth(gin.Accounts{config.Get().Username: config.Get().Password}))
	} else {
		glg.Warn("Username and password not defined, authentication disabled")
		group = r.Group("/")
	}

	for k, v := range handlersMap {
		group.GET(k, v)
	}

	r.Run(fmt.Sprintf("%s:%d", config.Get().Address, config.Get().Port))
}

func startHandler(c *gin.Context) {
	motionDetection, err := strconv.ParseBool(c.Query("detection"))

	if err != nil {
		motionDetection = false
	}

	err = motion.Startup(motionDetection)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion started"})
	}
}

func restartHandler(c *gin.Context) {
	err := motion.Restart()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion restarted"})
	}
}

func stopHandler(c *gin.Context) {
	err := motion.Shutdown()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion stopped"})
	}
}

func statusHandler(c *gin.Context) {
	started := motion.IsStarted()

	c.JSON(http.StatusOK, gin.H{"motionStarted": started})
}

func isMotionDetectionEnabled(c *gin.Context) {
	enabled, err := motion.IsMotionDetectionEnabled()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"motionDetectionEnabled": enabled})
	}
}

func startDetectionHandler(c *gin.Context) {
	err := motion.EnableMotionDetection(true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion detection started"})
	}
}

func pauseDetectionHandler(c *gin.Context) {
	err := motion.EnableMotionDetection(false)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion detection paused"})
	}
}

func listConfigHandler(c *gin.Context) {
	configMap, err := motion.ConfigGet(c.Query(""))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, configMap)
	}
}

func setConfigHandler(c *gin.Context) {
}

func getConfigHandler(c *gin.Context) {

}
