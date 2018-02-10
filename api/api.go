package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"

	"../config"
	"../motion"
)

func Init() {
	glg.Info("Initializing REST API ...")
	var group *gin.RouterGroup
	r := gin.Default()

	//eventGroup := r.Group("/event/")

	if config.Get().Username != "" && config.Get().Password != "" {
		glg.Info("Username and password defined, authentication enabled")
		group = r.Group("/api", gin.BasicAuth(gin.Accounts{config.Get().Username: config.Get().Password}))
	} else {
		glg.Warn("Username and password not defined, authentication disabled")
		group = r.Group("/api")
	}

	group.GET("/api/control/startup", startHandler)
	group.GET("/api/control/shutdown", stopHandler)
	group.GET("/api/control/restart", restartHandler)
	group.GET("/api/detection/status", isMotionDetectionEnabled)
	group.GET("/api/detection/start", startDetectionHandler)
	group.GET("/api/detection/pause", pauseDetectionHandler)
	group.GET("/api/stream", streamHandler)

	/*eventGroup.GET("/event/")
	eventGroup.GET("/event/")
	eventGroup.GET("/event/")*/

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
	err := motion.Restart()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "motion stopped"})
	}
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

func streamHandler(c *gin.Context) {

}
