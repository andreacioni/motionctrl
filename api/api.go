package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"

	"../config"
	"../motion"
)

const (
	Ok                  = 200
	InternalServerError = 500
)

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

	group.GET("/startup", startHandler)
	group.GET("/shutdown", stopHandler)
	group.GET("/restart", restartHandler)
	group.GET("/detection/status", isMotionDetectionEnabled)
	group.GET("/detection/start", startMotionDetection)
	group.GET("/detection/pause", pauseMotionDetection)

	r.Run(fmt.Sprintf("%s:%d", config.Get().Address, config.Get().Port))
}

func startHandler(c *gin.Context) {
	motionDetection, err := strconv.ParseBool(c.Query("detection"))

	if err != nil {
		motionDetection = false
	}

	err = motion.Startup(motionDetection)

	if err != nil {
		c.JSON(InternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(Ok, gin.H{"message": "motion started"})
	}
}

func stopHandler(c *gin.Context) {
	err := motion.Shutdown()

	if err != nil {
		c.JSON(InternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(Ok, gin.H{"message": "motion stopped"})
	}
}

func stopHandler(c *gin.Context) {
	err := motion.Restart()

	if err != nil {
		c.JSON(InternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(Ok, gin.H{"message": "motion stopped"})
	}
}

func isMotionDetectionEnabled(c *gin.Context) {
	enabled, err := motion.IsMotionDetectionEnabled()

	if err != nil {
		c.JSON(InternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(Ok, gin.H{"motionDetectionEnabled": enabled})
	}
}

func startMotionDetection(c *gin.Context) {
	err := motion.EnableMotionDetection(true)

	if err != nil {
		c.JSON(InternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(Ok, gin.H{"message": "motion detection started"})
	}
}

func pauseMotionDetection(c *gin.Context) {
	err := motion.EnableMotionDetection(false)

	if err != nil {
		c.JSON(InternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(Ok, gin.H{"message": "motion detection paused"})
	}
}
